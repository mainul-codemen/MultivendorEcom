package handler

import (
	"net/http"

	"github.com/MultivendorEcom/serviceutil/logger"
)

func (s *Server) adminindex(w http.ResponseWriter, r *http.Request) {
	tmpl := s.lookupTemplate("index.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}
