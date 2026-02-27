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

    resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
        Model: openai.ChatModelGPT4oMini,
        Messages: []openai.ChatCompletionMessageParamUnion{
            openai.SystemMessage("Sos el asistente de Peluquería Sol. Respondé solo sobre turnos, precios y horarios."),
            openai.UserMessage("¿Cuánto sale el corte de pelo?"),
        },
        Temperature: openai.Float(0.2),
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp.Choices[0].Message.Content)
    fmt.Printf("Tokens usados: %d\n", resp.Usage.TotalTokens)
}