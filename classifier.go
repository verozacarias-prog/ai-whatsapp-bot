package main

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/openai/openai-go/v3"
)

type ClassifyRequest struct {
	Message string `json:"message"`
}

type ClassifyResponse struct {
	Intent     string `json:"intent"`
	Confidence string `json:"confidence"`
	TokensUsed int64  `json:"tokens_used"`
}

func Classify(request ClassifyRequest) (ClassifyResponse, error) {
	client := openai.NewClient()

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(`Sos un clasificador de mensajes para una peluquería.
Tu única tarea es clasificar el mensaje del cliente en UNA de estas categorías:
- consulta_precio
- reserva_turno
- cancelacion
- consulta_horario
- otro

Respondé SOLO con un JSON válido con estos campos:
{ "intent": string, "confidence": "high|medium|low" }

No agregues explicaciones. No salgas del JSON.`),
		openai.UserMessage(request.Message),
	}

	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model:       openai.ChatModelGPT4oMini,
		Messages:    messages,
		Temperature: openai.Float(0),
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &openai.ResponseFormatJSONObjectParam{},
		},
	})
	if err != nil {
		return ClassifyResponse{}, fmt.Errorf("error de API: %w", err)
	}

	var result ClassifyResponse
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return ClassifyResponse{}, fmt.Errorf("error parseando respuesta: %w", err)
	}
	result.TokensUsed = resp.Usage.TotalTokens
	return result, nil
}
