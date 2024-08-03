package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Address содержит информацию об адресе.
type Address struct {
	Result     string `json:"result"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	Region     string `json:"region"`
	Street     string `json:"street"`
	GeoLat     string `json:"lat"`
	GeoLon     string `json:"lon"`
}

// GeocodeResponse представляет ответ на запрос геокодирования.
type GeocodeResponse struct {
	Suggestions []*Suggestion `json:"suggestions"`
}

// Suggestion содержит предложения по адресу.
type Suggestion struct {
	GeoLat string `json:"lat"`
	GeoLon string `json:"lon"`
	Value  string `json:"value"`
}

// SearchRequest представляет запрос на поиск адреса.
type SearchRequest struct {
	Query string `json:"query"`
}

// GeocodeRequest представляет запрос геокодирования.
type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

// SearchResponse представляет ответ на запрос поиска адреса.
type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}
