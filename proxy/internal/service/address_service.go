package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"geo-controller/proxy/internal/models"
	"net/http"

	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
)

type AddressService struct {
	daDataApiKey    string
	daDataSecretKey string
}

func NewAddressService(apiKey, secretKey string) *AddressService {
	return &AddressService{
		daDataApiKey:    apiKey,
		daDataSecretKey: secretKey,
	}
}

func (s *AddressService) SearchAddress(query string) (*models.SearchResponse, error) {
	creds := client.Credentials{
		ApiKeyValue:    s.daDataApiKey,
		SecretKeyValue: s.daDataSecretKey,
	}

	api := dadata.NewSuggestApi(client.WithCredentialProvider(&creds))

	params := suggest.RequestParams{
		Query: query,
	}

	suggestions, err := api.Address(context.Background(), &params)
	if err != nil {
		return nil, err
	}

	searchResp := &models.SearchResponse{}
	for _, s := range suggestions {
		addr := models.Address{
			Result:     s.Value,
			PostalCode: s.Data.PostalCode,
			Country:    s.Data.Country,
			Region:     s.Data.Region,
			Street:     s.Data.Street,
			GeoLat:     s.Data.GeoLat,
			GeoLon:     s.Data.GeoLon,
		}
		searchResp.Addresses = append(searchResp.Addresses, &addr)
	}

	return searchResp, nil
}

func (s *AddressService) Geocode(request models.GeocodeRequest) (*models.GeocodeResponse, error) {
	if request.Lat == "" || request.Lng == "" {
		return nil, errors.New("latitude and longitude cannot be empty")
	}

	requestData := map[string]string{
		"lat": request.Lat,
		"lon": request.Lng,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	url := "http://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token "+s.daDataApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var geocodeResp models.GeocodeResponse
	err = json.NewDecoder(resp.Body).Decode(&geocodeResp)
	if err != nil {
		return nil, err
	}

	return &geocodeResp, nil
}
