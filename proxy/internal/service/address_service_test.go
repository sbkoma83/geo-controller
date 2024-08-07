// address_service_test.go
package service_test

import (
	"geo-controller/proxy/internal/models"
	"geo-controller/proxy/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressService_SearchAddress_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	addressService := service.NewAddressService("testApiKey", "testSecretKey")
	addressService.DaDataURL = server.URL

	result, err := addressService.SearchAddress("Test query")

	assert.Error(t, err)
	assert.Nil(t, result)
}
func TestAddressService_Geocode_InvalidCoordinates(t *testing.T) {
	addressService := service.NewAddressService("testApiKey", "testSecretKey")

	request := models.GeocodeRequest{
		Lat: "",
		Lng: "",
	}

	result, err := addressService.Geocode(request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "latitude and longitude cannot be empty", err.Error())
}
func TestAddressService_Geocode_EmptyResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"suggestions": []}`))
	}))
	defer ts.Close()

	svc := service.NewAddressService("testApiKey", "testSecretKey")
	svc.DaDataURL = ts.URL

	request := models.GeocodeRequest{
		Lat: "12.34",
		Lng: "56.78",
	}

	result, err := svc.Geocode(request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Suggestions)
}
