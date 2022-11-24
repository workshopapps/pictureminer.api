package model

type MicroserviceResponse struct {
	Content string `json:"text_description"`
}

type MicroservicePromptResponse struct {
	CheckResult bool `json:"check_result"`
}
