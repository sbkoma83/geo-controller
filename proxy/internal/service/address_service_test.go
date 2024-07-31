package service

import (
	"geo-controller/proxy/internal/models"
	"testing"
)

func TestAddressService_SearchAddress(t *testing.T) {
	// Используйте тестовые ключи API или замените их на моки
	addressService := NewAddressService("c4aab5f0a277fbaa6de6613c4c78930552172d28", "e1e61bbed8ab858bc7153ba44fc8344ba7681526")

	query := "Moscow"
	resp, err := addressService.SearchAddress(query)
	if err != nil {
		t.Errorf("SearchAddress failed: %v", err)
	}

	if resp == nil || len(resp.Addresses) == 0 {
		t.Error("Expected non-empty address response, got empty")
	}

	// Проверьте, что результаты содержат ожидаемые поля
	for _, addr := range resp.Addresses {
		if addr.Country == "" {
			t.Error("Expected non-empty country in address")
		}
		if addr.GeoLat == "" || addr.GeoLon == "" {
			t.Error("Expected non-empty geo coordinates in address")
		}
	}
}
func TestAddressService_Geocode_EmptyCoordinates(t *testing.T) {
	addressService := NewAddressService("c4aab5f0a277fbaa6de6613c4c78930552172d28", "e1e61bbed8ab858bc7153ba44fc8344ba7681526")

	// Создаем запрос с пустыми координатами
	geocodeReq := models.GeocodeRequest{Lat: "", Lng: ""}
	_, err := addressService.Geocode(geocodeReq)
	if err == nil {
		t.Error("Expected error when geocoding with empty coordinates, but got nil")
	}
}
