package handler

import (
	"net/http"

	"github.com/MultivendorEcom/serviceutil/logger"
)

const AlrEx = " already exists"
const PosEx = "position %d already exists"
const PhnEx = "phone number %s already exists"

func (s *Server) countryList(r *http.Request, w http.ResponseWriter, sts bool) []CountryForm {
	cntryList, err := s.st.GetAllCountry(r.Context(), sts)
	if err != nil {
		logger.Error("error while get country : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	cntryListForm := make([]CountryForm, 0)
	for _, item := range cntryList {
		cntryData := CountryForm{
			ID:       item.ID,
			Name:     item.Name,
			Status:   item.Status,
			Position: item.Position,
		}
		cntryListForm = append(cntryListForm, cntryData)
	}
	return cntryListForm
}

func (s *Server) strToDisFormByID(w http.ResponseWriter, r *http.Request, id string) DistrictForm {
	res, err := s.st.GetDistrictBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get district " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	Data := DistrictForm{
		ID:          res.ID,
		Name:        res.Name,
		CountryID:   res.CountryID,
		CountryName: res.CountryName.String,
		Status:      res.Status,
		Position:    res.Position,
		CreatedBy:   res.CreatedBy,
	}
	return Data
}

func (s *Server) stationList(r *http.Request, w http.ResponseWriter, sts bool) []StationForm {
	stnList, err := s.st.GetStationList(r.Context(), sts)
	if err != nil {
		logger.Error("error while get station : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	stnListForm := make([]StationForm, 0)
	for _, item := range stnList {
		cntryData := StationForm{
			ID:           item.ID,
			Name:         item.Name,
			Status:       item.Status,
			Position:     item.Position,
			DistrictName: item.DistrictName.String,
			CreatedBy:    item.CreatedBy,
		}
		stnListForm = append(stnListForm, cntryData)
	}
	return stnListForm
}

func (s *Server) districtList(r *http.Request, w http.ResponseWriter, sts bool) []DistrictForm {
	dist, err := s.st.GetDistrictList(r.Context(), sts)
	if err != nil {
		logger.Error("error while get district : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	disdata := make([]DistrictForm, 0)
	for _, item := range dist {
		disApnd := DistrictForm{
			ID:       item.ID,
			Name:     item.Name,
			Status:   item.Status,
			Position: item.Position,
		}
		disdata = append(disdata, disApnd)
	}
	return disdata
}

func (s *Server) hubList(r *http.Request, w http.ResponseWriter, sts bool) []HubForm {
	hubList, err := s.st.GetHubList(r.Context(),sts)
	if err != nil {
		logger.Error("error while get hub : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	hbdata := make([]HubForm, 0)
	for _, item := range hubList {
		disApnd := HubForm{
			ID:           item.ID,
			Name:         item.HubName,
			CountryID:    item.CountryID,
			CountryName:  item.CountryName.String,
			DistrictID:   item.DistrictID,
			DistrictName: item.DistrictName.String,
			StationID:    item.StationID,
			StationName:  item.StationName.String,
			HubPhone1:    item.HubPhone1,
			HubPhone2:    item.HubPhone2,
			HubEmail:     item.HubEmail,
			HubAddress:   item.HubAddress,
			Status:       item.Status,
			Position:     item.Position,
			CreatedAt:    item.CreatedAt,
			CreatedBy:    item.CreatedBy,
			UpdatedAt:    item.UpdatedAt,
			UpdatedBy:    item.UpdatedBy,
		}
		hbdata = append(hbdata, disApnd)
	}
	return hbdata
}
