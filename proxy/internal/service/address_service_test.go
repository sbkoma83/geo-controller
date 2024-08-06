// address_service_test.go
package service_test

import (
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
