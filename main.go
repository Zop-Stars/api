package main

import (
	"gofr.dev/pkg/gofr"
	vertex_ai "gofr.dev/pkg/gofr/ai/vertex-ai"
)

func main() {
	app := gofr.New()

	creds := app.Config.Get("SVC_ACC_CREDS")

	vertexAIConfigs := &vertex_ai.Configs{
		ProjectID:   "endless-fire-437206-j7",
		LocationID:  "us-central1",
		APIEndpoint: "us-central1-aiplatform.googleapis.com",
		Datastore:   "projects/endless-fire-437206-j7/locations/global/collections/default_collection/dataStores/gofr-datastore_1732298621027",
		ModelID:     "gemini-1.5-pro-002",
		Credentials: creds,
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

		return c.VertexAI.SendMessageUsingSystemInstruction(prompt, []string{"use only one documentation i.e. using-http-server"})
	})

	app.Run()
}
