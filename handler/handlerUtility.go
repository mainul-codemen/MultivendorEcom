package handler

import (
	crand "crypto/rand"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
	"github.com/gorilla/csrf"
)

const (
	ewrt    = "error while read template"
	ult     = "unable to load template"
	ewte    = "error with template execution"
	errMsg  = "parsing form. "
	deErr   = "error while decoding form. "
	nameReq = `The name is required`
	posReq  = `The Position is required`
)

var msg = storage.Message{
	Status:  true,
	Message: "Successfully Save Data",
}

func cacheStaticFiles(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if asset is hashed extend cache to 180 days
		e := `"4FROTHS24N"`
		w.Header().Set("Etag", e)
		w.Header().Set("Cache-Control", "max-age=15552000")
		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, e) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func (s *Server) parseTemplates() error {
	templates := template.New("templates").Funcs(template.FuncMap{
		"assetHash": func(n string) string {
			return path.Join("/", s.assetFS.HashName(strings.TrimPrefix(path.Clean(n), "/")))
		},
	}).Funcs(sprig.FuncMap())

	tmpl, err := templates.ParseFS(s.assets, "templates/*/*.html")
	if err != nil {
		return err
	}
	s.templates = tmpl
	return nil
}

func (s *Server) lookupTemplate(name string) *template.Template {
	if s.env == "development" {
		if err := s.parseTemplates(); err != nil {
			s.logger.Error("template reload")
			return nil
		}
	}
	return s.templates.Lookup(name)
}

func (s *Server) getErrorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl := s.lookupTemplate("error.html")
		if tmpl == nil {
			logger.Error(ult)
			http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
		}
		if err := tmpl.Execute(w, nil); err != nil {
			logger.Error(ewte + err.Error())
			http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
		}
	})
}

func isPartialTemplate(name string) bool {
	return strings.HasSuffix(name, ".part.html")
}

func (s *Server) templateData(r *http.Request) TemplateData {
	return TemplateData{
		Env:       s.env,
		CSRFField: csrf.TemplateField(r),
	}
}

func RandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, n)
	read, err := crand.Read(bytes)
	if err != nil {
		return ""
	}
	if read != n {
		return ""
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func (s *Server) validateSingleFileType(r *http.Request, name string, types []string) error {
	f, _, err := r.FormFile(name)
	if err != nil {
		return nil
	}
	defer f.Close()
	c, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	ft := http.DetectContentType(c)
	rtn := false
	for _, t := range types {
		if ft == t {
			rtn = true
		}
	}
	if !rtn {
		return errors.New("invalid file format")
	}
	return nil
}

func (s *Server) saveImage(file multipart.File, fileHeader *multipart.FileHeader, imagePath string) (string, error) {
	tt := time.Now().Local()
	image := tt.Format("20060102") + RandomString(5) + fileHeader.Filename
	dest, err := os.Create(fmt.Sprintf(imagePath+"%s", image))
	if err != nil {
		return "", err
	}
	defer dest.Close()
	if _, err := io.Copy(dest, file); err != nil {
		fmt.Println(err.Error())
	}
	return image, err
}

func (s *Server) stringToDate(date string) time.Time {
	layout := "2006-01-02"
	fdate, err := time.Parse(layout, date)
	if err != nil {
		fmt.Println(err)
	}
	return fdate
}

func trim(str string) string {
	return strings.Trim(str, " ")
}
