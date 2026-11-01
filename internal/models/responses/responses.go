package responses

type GenericResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Invalid request"`
}
