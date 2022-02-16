package handler

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/MultivendorEcom/storage/postgres"
	"github.com/benbjohnson/hashfs"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Server struct {
	templates *template.Template
	env       string
	logger    *zap.Logger
	decoder   *schema.Decoder
	config    *viper.Viper
	assets    fs.FS
	assetFS   *hashfs.FS
	st        *postgres.Storage
}
type TemplateData struct {
	Env       string
	CSRFField template.HTML
}

const (
	// user
	userListPath         = "/user"
	createUserPath       = "/user/create"
	updateUserPath       = "/user/update/{id}"
	updateUserStatusPath = "/user/update/status/{id}"
	viewUserPath         = "/user/view/{id}"
	deleteUserPath       = "/user/delete/{id}"
	// designation
	designationListPath         = "/designation"
	designationCreate           = "/designation/create"
	deleteDesignationPath       = "/designation/delete/{id}"
	updateDesignationPath       = "/designation/update/{id}"
	updateDesignationStatusPath = "/designation/update/status/{id}"
	viewDesignationPath         = "/designation/view/{id}"
	// country
	countryListPath         = "/country"
	createCountryPath       = "/country/create"
	updateCountryPath       = "/country/update/{id}"
	viewCountryPath         = "/country/view/{id}"
	updateCountryStatusPath = "/country/update/status/{id}"
	deleteCountryPath       = "/country/delete/{id}"
	// district
	districtListPath         = "/district"
	createDistrictPath       = "/district/create"
	updateDistrictPath       = "/district/update/{id}"
	updateDistrictStatusPath = "/district/update/status/{id}"
	viewDistrictPath         = "/district/view/{id}"
	deleteDistrictPath       = "/district/delete/{id}"
	// station
	stationListPath         = "/station"
	createStationPath       = "/station/create"
	updateStationPath       = "/station/update/{id}"
	updateStationStatusPath = "/station/update/status/{id}"
	viewStationPath         = "/station/view/{id}"
	deleteStationPath       = "/station/delete/{id}"
	// hub
	hubListPath         = "/hub"
	createHubPath       = "/hub/create"
	updateHubPath       = "/hub/update/{id}"
	updateHubStatusPath = "/hub/update/status/{id}"
	viewHubPath         = "/hub/view/{id}"
	deleteHubPath       = "/hub/delete/{id}"
	// branch
	branchListPath         = "/branch"
	createBranchPath       = "/branch/create"
	updateBranchPath       = "/branch/update/{id}"
	updateBranchStatusPath = "/branch/update/status/{id}"
	viewBranchPath         = "/branch/view/{id}"
	deleteBranchPath       = "/branch/delete/{id}"
	// delivery-company
	deliveryCompanyListPath         = "/delivery-company"
	createDeliveryCompanyPath       = "/delivery-company/create"
	updateDeliveryCompanyPath       = "/delivery-company/update/{id}"
	updateDeliveryCompanyStatusPath = "/delivery-company/update/status/{id}"
	viewDeliveryCompanyPath         = "/delivery-company/view/{id}"
	deleteDeliveryCompanyPath       = "/delivery-company/delete/{id}"
	// delivery-charge
	deliveryChargeListPath         = "/delivery-charge"
	createDeliveryChargePath       = "/delivery-charge/create"
	updateDeliveryChargePath       = "/delivery-charge/update/{id}"
	updateDeliveryChargeStatusPath = "/delivery-charge/update/status/{id}"
	viewDeliveryChargePath         = "/delivery-charge/view/{id}"
	deleteDeliveryChargePath       = "/delivery-charge/delete/{id}"

	ErrorPath = "/error"
)

func New(
	env string,
	config *viper.Viper,
	logger *zap.Logger,
	decoder *schema.Decoder,
	assets fs.FS,
	st *postgres.Storage,
) (*mux.Router, error) {
	s := &Server{
		env:     env,
		logger:  logger,
		decoder: decoder,
		config:  config,
		assets:  assets,
		st:      st,
	}
	r := mux.NewRouter()
	if err := s.parseTemplates(); err != nil {
		logger.Error("Error in parse templates:")
		return nil, err
	}
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/"))))
	ar := r.PathPrefix("/admin").Subrouter()
	ar.HandleFunc("/index", s.adminindex).Methods("GET")
	// User
	ar.HandleFunc(userListPath, s.userListHandler).Methods("GET")
	ar.HandleFunc(createUserPath, s.usrFormHandler).Methods("GET")
	ar.HandleFunc(createUserPath, s.submitUserHandler).Methods("POST")
	// ar.HandleFunc(updateUserPath, s.updateUserFormHandler).Methods("GET")
	ar.HandleFunc(updateUserPath, s.updateUserHandler).Methods("POST")
	ar.HandleFunc(viewUserPath, s.viewUserHandler).Methods("GET")
	ar.HandleFunc(updateUserStatusPath, s.updateUserStatusHandler).Methods("GET")
	ar.HandleFunc(deleteUserPath, s.deleteUserHandler).Methods("GET")
	// designation
	ar.HandleFunc(designationCreate, s.submitDesignation).Methods("POST")
	ar.HandleFunc(updateDesignationPath, s.updateDesignation).Methods("POST")
	ar.HandleFunc(deleteDesignationPath, s.deleteDesignation).Methods("GET")
	ar.HandleFunc(designationListPath, s.designationList).Methods("GET")
	ar.HandleFunc(viewDesignationPath, s.viewDesignation).Methods("GET")
	ar.HandleFunc(updateDesignationStatusPath, s.updateDesignationStatus).Methods("GET")
	// country
	ar.HandleFunc(countryListPath, s.countryListHandler).Methods("GET")
	ar.HandleFunc(createCountryPath, s.submitCountry).Methods("POST")
	ar.HandleFunc(updateCountryPath, s.updateCountryHadler).Methods("POST")
	ar.HandleFunc(viewCountryPath, s.viewCountryHandler).Methods("GET")
	ar.HandleFunc(updateCountryStatusPath, s.updateCountryStatusHandler).Methods("GET")
	ar.HandleFunc(deleteCountryPath, s.deleteCountryHandler).Methods("GET")
	// district
	ar.HandleFunc(districtListPath, s.districtListHandler).Methods("GET")
	ar.HandleFunc(createDistrictPath, s.districtFormHandler).Methods("GET")
	ar.HandleFunc(createDistrictPath, s.submitDistrictHandler).Methods("POST")
	ar.HandleFunc(updateDistrictPath, s.updateDistrictFormHandler).Methods("GET")
	ar.HandleFunc(updateDistrictPath, s.updateDistrictHandler).Methods("POST")
	ar.HandleFunc(viewDistrictPath, s.viewDistrictHandler).Methods("GET")
	ar.HandleFunc(updateDistrictStatusPath, s.updateDistrictStatusHandler).Methods("GET")
	ar.HandleFunc(deleteDistrictPath, s.deleteDistrictHandler).Methods("GET")
	// station
	ar.HandleFunc(stationListPath, s.stationListHandler).Methods("GET")
	ar.HandleFunc(createStationPath, s.stationFormHandler).Methods("GET")
	ar.HandleFunc(createStationPath, s.submitStationHandler).Methods("POST")
	ar.HandleFunc(updateStationPath, s.updateStationFormHandler).Methods("GET")
	ar.HandleFunc(updateStationPath, s.updateStationHandler).Methods("POST")
	ar.HandleFunc(viewStationPath, s.viewStationHandler).Methods("GET")
	ar.HandleFunc(updateStationStatusPath, s.updateStationStatusHandler).Methods("GET")
	ar.HandleFunc(deleteStationPath, s.deleteStationHandler).Methods("GET")
	// hub
	ar.HandleFunc(hubListPath, s.hubListHandler).Methods("GET")
	ar.HandleFunc(createHubPath, s.hubFormHandler).Methods("GET")
	ar.HandleFunc(createHubPath, s.submitHubHandler).Methods("POST")
	ar.HandleFunc(updateHubPath, s.updateHubFormHandler).Methods("GET")
	ar.HandleFunc(updateHubPath, s.updateHubHandler).Methods("POST")
	ar.HandleFunc(viewHubPath, s.viewHubHandler).Methods("GET")
	ar.HandleFunc(updateHubStatusPath, s.updateHubStatusHandler).Methods("GET")
	ar.HandleFunc(deleteHubPath, s.deleteHubHandler).Methods("GET")
	// branch
	ar.HandleFunc(branchListPath, s.branchListHandler).Methods("GET")
	ar.HandleFunc(createBranchPath, s.branchFormHandler).Methods("GET")
	ar.HandleFunc(createBranchPath, s.submitBranchHandler).Methods("POST")
	ar.HandleFunc(updateBranchPath, s.updateBranchFormHandler).Methods("GET")
	ar.HandleFunc(updateBranchPath, s.updateBranchHandler).Methods("POST")
	ar.HandleFunc(viewBranchPath, s.viewBranchHandler).Methods("GET")
	ar.HandleFunc(updateBranchStatusPath, s.updateBranchStatusHandler).Methods("GET")
	ar.HandleFunc(deleteBranchPath, s.deleteBranchHandler).Methods("GET")
	// deliveryCompany
	ar.HandleFunc(deliveryCompanyListPath, s.deliveryCompanyListHandler).Methods("GET")
	ar.HandleFunc(createDeliveryCompanyPath, s.deliveryCompanyFormHandler).Methods("GET")
	ar.HandleFunc(createDeliveryCompanyPath, s.submitDeliveryCompanyHandler).Methods("POST")
	ar.HandleFunc(updateDeliveryCompanyPath, s.updateDeliveryCompanyFormHandler).Methods("GET")
	ar.HandleFunc(updateDeliveryCompanyPath, s.updateDeliveryCompanyHandler).Methods("POST")
	ar.HandleFunc(viewDeliveryCompanyPath, s.viewDeliveryCompanyHandler).Methods("GET")
	ar.HandleFunc(updateDeliveryCompanyStatusPath, s.updateDeliveryCompanyStatusHandler).Methods("GET")
	ar.HandleFunc(deleteDeliveryCompanyPath, s.deleteDeliveryCompanyHandler).Methods("GET")

	// deliveryCharge
	ar.HandleFunc(deliveryChargeListPath, s.deliveryChargeListHandler).Methods("GET")
	ar.HandleFunc(createDeliveryChargePath, s.deliveryChargeFormHandler).Methods("GET")
	ar.HandleFunc(createDeliveryChargePath, s.submitDeliveryChargeHandler).Methods("POST")
	ar.HandleFunc(updateDeliveryChargePath, s.updateDeliveryChargeFormHandler).Methods("GET")
	ar.HandleFunc(updateDeliveryChargePath, s.updateDeliveryChargeHandler).Methods("POST")
	ar.HandleFunc(viewDeliveryChargePath, s.viewDeliveryChargeHandler).Methods("GET")
	ar.HandleFunc(updateDeliveryChargeStatusPath, s.updateDeliveryChargeStatusHandler).Methods("GET")
	ar.HandleFunc(deleteDeliveryChargePath, s.deleteDeliveryChargeHandler).Methods("GET")
	r.NotFoundHandler = s.getErrorHandler()
	return r, nil

}
