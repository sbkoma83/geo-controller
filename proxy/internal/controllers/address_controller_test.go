// file: controllers/address_controller_test.go

package controllers

import (
	"bytes"
	"encoding/json"
	"geo-controller/proxy/internal/models"
	"geo-controller/proxy/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddressController_AddressSearchHandler(t *testing.T) {
	addressService := service.NewAddressService("c4aab5f0a277fbaa6de6613c4c78930552172d28", "e1e61bbed8ab858bc7153ba44fc8344ba7681526")
	addressController := NewAddressController(addressService)

	// Создаем запрос с телом JSON
	searchReq := models.SearchRequest{Query: "Moscow"}
	reqBody, _ := json.Marshal(searchReq)
	req, err := http.NewRequest("POST", "/api/address/search", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addressController.AddressSearchHandler)

	// Вызываем обработчик
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем содержимое ответа
	var response models.SearchResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if len(response.Addresses) == 0 {
		t.Error("Expected non-empty address response, got empty")
	}
}

func TestAddressController_GeocodeHandler(t *testing.T) {
	addressService := service.NewAddressService("c4aab5f0a277fbaa6de6613c4c78930552172d28", "e1e61bbed8ab858bc7153ba44fc8344ba7681526")
	addressController := NewAddressController(addressService)

	// Создаем запрос с телом JSON
	geocodeReq := models.GeocodeRequest{Lat: "55.7558", Lng: "37.6176"}
	reqBody, _ := json.Marshal(geocodeReq)
	req, err := http.NewRequest("POST", "/api/address/geocode", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addressController.GeocodeHandler)

	// Вызываем обработчик
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем содержимое ответа
	var response models.GeocodeResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if len(response.Suggestions) == 0 {
		t.Error("Expected non-empty geocode response, got empty")
	}
}
func TestAddressController_AddressSearchHandler_BadRequest(t *testing.T) {
	addressService := service.NewAddressService("c4aab5f0a277fbaa6de6613c4c78930552172d28", "e1e61bbed8ab858bc7153ba44fc8344ba7681526")
	addressController := NewAddressController(addressService)

	// Создаем запрос с некорректным телом JSON
	req, err := http.NewRequest("POST", "/api/address/search", bytes.NewBuffer([]byte("{invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addressController.AddressSearchHandler)

	// Вызываем обработчик
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
func TestAddressController_GeocodeHandler_BadRequest(t *testing.T) {
	addressService := service.NewAddressService("c4aab5f0a277fbaa6de6613c4c78930552172d28", "e1e61bbed8ab858bc7153ba44fc8344ba7681526")
	addressController := NewAddressController(addressService)

	// Создаем запрос с некорректным телом JSON
	req, err := http.NewRequest("POST", "/api/address/geocode", bytes.NewBuffer([]byte("{invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addressController.GeocodeHandler)

	// Вызываем обработчик
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
