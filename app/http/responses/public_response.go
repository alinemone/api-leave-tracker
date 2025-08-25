package responses

type PublicResponse struct {
	Result interface{} `json:"result"`
}

func NewPublicResponse(result interface{}) PublicResponse {
	return PublicResponse{Result: result}
}
