package controllers


import (
	"geo-controller/proxy/internal/models"
	"geo-controller/proxy/internal/responder"
	"geo-controller/proxy/internal/service"
	"encoding/json"
	"net/http"
)

type AddressController struct {
	addressService *service.AddressService
	responder      *responder.Responder
}

func NewAddressController(addressService *service.AddressService) *AddressController {
	return &AddressController{
		addressService: addressService,
		responder:      responder.NewResponder(),
	}
}

func (c *AddressController) AddressSearchHandler(w http.ResponseWriter, r *http.Request) {
	var searchReq models.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&searchReq); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	searchResp, err := c.addressService.SearchAddress(searchReq.Query)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	c.responder.OutputJSON(w, searchResp)
}

func (c *AddressController) GeocodeHandler(w http.ResponseWriter, r *http.Request) {
	var geocodeReq models.GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&geocodeReq); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	geocodeResp, err := c.addressService.Geocode(geocodeReq)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	c.responder.OutputJSON(w, geocodeResp)
}

