package main

import (
	"net/http"
	"strings"

	"gofr.dev/pkg/gofr"
	vertex_ai "gofr.dev/pkg/gofr/ai/vertex-ai"
	gofrHTTPPkg "gofr.dev/pkg/gofr/http"
)

type ErrorSendMessage struct{}

func (e ErrorSendMessage) Error() string {
	return "failed to send message"
}

func (e ErrorSendMessage) StatusCode() int {
	return http.StatusBadRequest
}

func main() {
	app := gofr.New()

	creds := app.Config.Get("SVC_ACC_CREDS")

	vertexAIConfigs := &vertex_ai.Configs{
		ProjectID:         app.Config.Get("SVC_ACC_PROJECT_ID"),
		LocationID:        app.Config.Get("SVC_ACC_LOCATION_ID"),
		APIEndpoint:       app.Config.Get("SVC_ACC_ENDPOINT"),
		Datastore:         strings.Split(app.Config.Get("SVC_ACC_DATASTORE"), ","),
		SystemInstruction: []string{app.Config.Get("SVC_ACC_SYSTEM_INSTRUCTION")},
		ModelID:           app.Config.Get("SVC_ACC_MODEL_ID"),
		Credentials:       creds,
	}

	vertexAIClient, err := vertex_ai.NewVertexAIClientWithKey(vertexAIConfigs)
	if err != nil {
		app.Logger().Fatalf("failed to create vertex AI client: %v", err)
	}

	app.UseVertexAI(vertexAIClient)

	app.POST("/chat", func(c *gofr.Context) (interface{}, error) {
		var prompt []map[string]string

		err = c.Bind(&prompt)
		if err != nil {
			return nil, err
		}

		if len(prompt) == 0 {
			return nil, gofrHTTPPkg.ErrorMissingParam{Params: []string{"body"}}
		}

		resp, err := c.VertexAI.SendMessage(c, prompt)
		if err != nil {
			c.Logger.Errorf("failed to send message: %v", err)

			return nil, ErrorSendMessage{}
		}

		return resp, nil
	})

	app.Run()
}
