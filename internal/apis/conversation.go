package apis

type AIResponse struct {
	Valid bool `json:"valid,omitempty"`
	//RawConfiguration map[string]interface{} `json:"raw_configuration,omitempty"`
	RawConfiguration string `json:"raw_configuration,omitempty"`
	ErrorMessages    string `json:"error_messages,omitempty"`
}
