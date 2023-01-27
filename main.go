package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/jepbarasgarov/slack-bot-AI/models"
	"github.com/jepbarasgarov/slack-bot-AI/utils"

	"github.com/shomali11/slacker"
	"github.com/sirupsen/logrus"
)

func startBot() error {
	bot := createBot()
	def := createAIChatDefinitionToBot()

	bot.Command("AI: <qstn>", def)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		return err
	}
	return nil
}

func createBot() (bot *slacker.Slacker) {
	bot = slacker.NewClient(utils.Config.SlackBotToken, utils.Config.SlackAppToken)
	return
}

func createAIChatDefinitionToBot() (def *slacker.CommandDefinition) {
	def = &slacker.CommandDefinition{
		Description: "AI chat",
		Examples:    []string{"chat"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			qstn := request.Param("qstn")

			//check question
			if len(qstn) == 0 || len(qstn) > 512 {
				response.Reply("Question size is not compatible")
				return
			}

			// Prepare the request body to openAI
			requestToOpenAI, err := generateRequestToOpenAI(qstn)
			if err != nil {
				response.Reply("Couldn't generate message")
				return
			}

			// Send the request to OpenAI
			client := &http.Client{}
			resp, err := client.Do(requestToOpenAI)
			if err != nil {
				response.Reply("couldn't get answer")
				return
			}
			defer resp.Body.Close()

			// Read the response
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				response.Reply("couldn't get answer")
				return
			}

			// Unmarshal the response
			var r map[string]interface{}
			json.Unmarshal([]byte(body), &r)

			// get the answer
			choices := r["choices"].([]interface{})
			var generatedText string
			for _, choice := range choices {
				choiceMap := choice.(map[string]interface{})
				generatedText = choiceMap["text"].(string)
				break
			}

			response.Reply(generatedText)

		},
	}

	return def
}

func generateRequestToOpenAI(qstn string) (*http.Request, error) {
	url := "https://api.openai.com/v1/engines/text-davinci-002/completions"

	// Prepare the request body to openAI
	req := models.RequestModelToOpenAI
	req.Prompt = qstn
	req.MaxToken = 1024
	req.Temperature = 0.7

	jsonString, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New("Couldn't generate message")
	}

	requestToOpenAI, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	requestToOpenAI.Header.Set("Content-Type", "application/json")
	requestToOpenAI.Header.Set("Authorization", "Bearer "+utils.Config.OpenApiKey)

	return requestToOpenAI, nil
}

func main() {

	err := utils.ReadConfig("config.json")
	if err != nil {
		eMsg := "Error in reading configuration"
		logrus.WithError(err).Error(eMsg)
		return
	}

	err = startBot()
	if err != nil {
		logrus.WithError(err).Error("Couldn't start slack bot")
		return
	}

}
