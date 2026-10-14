package responses

type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
