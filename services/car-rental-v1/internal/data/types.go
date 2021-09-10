package data

type Car struct {
	RentalCompany string  `json:"rental_company"`
	City          string  `json:"city"`
	BodyType      string  `json:"body_type"`
	Cost          float64 `json:"cost"`
	Name          string  `json:"name"`
	Country       string  `json:"country"`
	Image         string  `json:"image"`
	CarId         string  `json:"car_id"`
	Id            string  `json:"id"`
	Style         string  `json:"style"`
	DateFrom      string  `json:"dateFrom"`
	DateTo        string  `json:"dateTo"`
}

type CustomType struct {
	Implement string `json:"implement"`
}

type Error struct {
	Error string `json:"error"`
}
