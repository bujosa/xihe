package api

type CreateModelInput struct {
	Name   string `json:"name"`
	Brand  string `json:"brand"`
	Status string `json:"status"`
}

type CreateModelResponse struct {
	Id string `json:"id"`
}
