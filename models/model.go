package models

type Configuration struct {
	SlackBotToken string `json:"slack_bot_token"`
	SlackAppToken string `json:"slack_app_token"`
	OpenAIApiKey  string `json:"open_ai_api_key"`
}

var RequestModelToOpenAI struct {
	Prompt      string  `json:"prompt"`
	MaxToken    int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}
