package api

type CreateModelInput struct {
	Name   string `json:"name"`
	Brand  string `json:"brand"`
	Status string `json:"status"`
}

type Model struct {
	Id string `json:"id"`
}

type CreateModelResponse struct {
	Data struct {
		CreateModel Model `json:"createModel"`
	} `json:"data"`
}

func CreateModel(createModelInput CreateModelInput) (Model, error) {
	return Model{}, nil
}
