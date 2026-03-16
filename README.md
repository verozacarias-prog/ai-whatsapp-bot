# ai-whatsapp-bot

WhatsApp automation bot for small businesses. Classifies customer messages by intent, keeps short-term conversation memory per phone number, and generates context-aware responses using OpenAI's API.

## Features

- **Intent classification** — categorizes incoming messages into: `consulta_precio`, `reserva_turno`, `cancelacion`, `consulta_horario`, `otro`
- **Conversation memory (per phone, per day)** — keeps the dialog history for each phone number during the day to give more contextual answers
- **Automated responses** — generates natural language replies based on business configuration
- **Business hours aware** — uses a structured schedule from `business.yaml` to send a welcome message on first contact and a custom out-of-hours message when the business is closed
- **Prompt injection protection** — system prompt hardened against common attack vectors
- **Retry with exponential backoff** — handles OpenAI rate limits and transient errors gracefully
- **Per-request logging** — logs message, intent, confidence and token usage to CSV
- **Business config via YAML** — no code changes needed to onboard a new client

## Tech Stack

- Go 1.25
- OpenAI API (`gpt-4o-mini`)
- OpenAI Go SDK v3

## Project Structure

```
ai-whatsapp-bot/
├── main.go        # HTTP server, route registration
├── classifier.go  # Intent classification logic
├── responder.go   # Response generation + conversation memory
├── conversation.go# In-memory history per phone number
├── config.go      # Business config loader (YAML)
├── retry.go       # Exponential backoff retry utility
├── utils.go       # CSV logging
├── business.yaml  # Business configuration (per client)
├── .env           # API keys (not committed)
└── logs.csv       # Generated at runtime (not committed)
```

## Endpoints

### `POST /classify`

Classifies a customer message by intent.

**Request**
```json
{ "message": "quiero sacar turno para el jueves" }
```

**Response**
```json
{
  "intent": "reserva_turno",
  "confidence": "high",
  "tokens_used": 244
}
```

---

### `POST /respond`

Classifies the message, uses the phone number to retrieve the recent conversation history for that customer, and generates a natural language response based on business configuration and context.

**Request**
```json
{
  "message": "cuánto sale el corte con lavado?",
  "phone": "+5491112345678"
}
```

**Response**
```json
{
  "reply": "El corte con lavado tiene un precio de $7.000. Si tenés más preguntas, no dudes en consultar.",
  "intent": "consulta_precio",
  "confidence": "high",
  "tokens_used": 519
}
```

Behavior:

- **First message of the day for a phone**: the API responds with the configured `welcome_message` from `business.yaml` and does not call OpenAI yet.
- **Outside business hours**: the API always responds with the configured `out_of_hours_message` from `business.yaml` and skips the OpenAI call (no tokens consumed).
- **Within business hours with history**: the API uses the stored history for that phone (only today’s messages) plus the latest classified intent to generate the answer.

## Setup

**1. Clone the repo**
```bash
git clone https://github.com/vzacarias/ai-whatsapp-bot
cd ai-whatsapp-bot
```

**2. Create `.env`**
```env
OPENAI_API_KEY=your-api-key-here
```

**3. Configure the business**

Edit `business.yaml` with the client's data:
```yaml
name: "Peluquería Example"
phone: "11-1234-5678"
business_hours:
  - days: [1, 2, 3, 4, 5]  # 0=domingo, 6=sábado
    open: "09:00"
    close: "20:00"
  - days: [6]
    open: "09:00"
    close: "14:00"
welcome_message: "¡Hola! Bienvenido/a a Peluquería Example. ¿En qué te puedo ayudar?"
out_of_hours_message: "Estamos fuera de horario, te respondemos mañana."
services:
  - name: "Corte"
    price: "$5.000"
  - name: "Corte con lavado"
    price: "$7.000"
  - name: "Coloración"
    price: "$15.000"
```

**4. Run**
```bash
go run .
```

Server starts on `:8080`.

## Intent Classification

| Intent | Example message |
|--------|----------------|
| `consulta_precio` | "¿cuánto sale el corte?" |
| `reserva_turno` | "quiero turno para el jueves" |
| `cancelacion` | "no voy a poder ir mañana" |
| `consulta_horario` | "¿a qué hora abren?" |
| `otro` | anything outside scope |

Messages classified as `otro` receive a fallback response directing the customer to call the business. No API call is made for the response generation step, saving tokens.

## Security

The system prompt is hardened against prompt injection using:
- Explicit scope restriction
- Numbered strict rules
- Immutability instruction
- Triple-quote delimiters around user input

## Error Handling

Retries use exponential backoff (1s → 2s → 4s) for transient errors (429 rate limit, 500, 503). Quota exhaustion (429 `insufficient_quota`) is logged and not retried. After 3 failed attempts the API returns `503` with a user-friendly message.

## Cost Estimate

At ~550 tokens per request with `gpt-4o-mini`:

| Messages/month | Estimated API cost |
|---------------|-------------------|
| 200 | ~$0.07 USD |
| 600 | ~$0.20 USD |
| 1,000 | ~$0.33 USD |

## Roadmap

- [x] Conversation memory (multi-turn history per phone, resets daily)
- [ ] n8n integration via webhook
- [ ] WhatsApp connection via Meta API
- [ ] RAG over business catalog (dynamic pricing)
- [ ] Google Calendar integration for appointments
- [ ] Streamlit admin panel for business owners

## Author

Verónica Zacarías — Backend Engineer transitioning to AI Engineering.  
[LinkedIn](https://www.linkedin.com/in/veronicazacarias1983/) · [GitHub](https://github.com/verozacarias-prog)
