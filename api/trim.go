package api

type CreateTrimInput struct {
	Alias        string   `json:"alias"`
	Name         string   `json:"name"`
	Status       string   `json:"status"`
	Model        string   `json:"carModel"`
	FuelType     string   `json:"fuelType"`
	DriveTrain   string   `json:"driveTrain"`
	Features     []string `json:"features"`
	BodyStyle    string   `json:"bodyStyle"`
	Transmission string   `json:"transmission"`
	Year         int      `json:"year"`
	Price        int      `json:"price"`
}

type CreateTrimResponse struct {
	Id string `json:"id"`
}
