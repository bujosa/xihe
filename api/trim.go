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

type Trim struct {
	Id string `json:"id"`
}

type CreateTrimResponse struct {
	Data struct {
		CreateTrim Trim `json:"createTrimLevel"`
	} `json:"data"`
}

func CreateTrim(createTrimInput CreateTrimInput) (Trim, error) {
	// TODO: Implement this
	return Trim{}, nil
}
