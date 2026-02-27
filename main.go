package main

import (
    "context"
    "fmt"
    "log"

    "github.com/openai/openai-go"
    "github.com/joho/godotenv"
)

func main() {
    // Cargar .env
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error cargando .env")
    }

    client := openai.NewClient()

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(`Sos el asistente de Peluquería Sol. Clasificá el mensaje del usuario y respondé ÚNICAMENTE con un JSON con este formato exacto:
	{
	  "intencion": "consulta_precio" | "reserva_turno" | "cancelacion" | "otro",
	  "confianza": "alta" | "media" | "baja"
	}
	No agregues texto antes ni después del JSON.`),
		openai.UserMessage("¿cuánto sale cortarse el pelo?"),
	}

    resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
        Model: openai.ChatModelGPT4oMini,
        Messages: messages,
        Temperature: openai.Float(0.2),
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp.Choices[0].Message.Content)
    fmt.Printf("Tokens usados: %d\n", resp.Usage.TotalTokens)
}