package service

type GetRequest struct {
	IP string `json:"ip_address"`
}

type GeoLocationResponse struct {
	IP           string `json:"ip_address"`
	Country      string `json:"country"`
	CountryCode  string `json:"country_code"`
	City         string `json:"city"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	MysteryValue string `json:"mystery_value"`
}
