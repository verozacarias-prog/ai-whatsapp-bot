package main

import (
	"sync"

	openai "github.com/openai/openai-go/v3"
)

var (
	conversationHistory = make(map[string][]openai.ChatCompletionMessageParamUnion)
	historyMu           sync.RWMutex
)

func GetHistory(phone string) []openai.ChatCompletionMessageParamUnion {
	historyMu.RLock()
	defer historyMu.RUnlock()
	return conversationHistory[phone]
}

func AddToHistory(phone, userMsg, assistantReply string) {
	historyMu.Lock()
	defer historyMu.Unlock()
	conversationHistory[phone] = append(conversationHistory[phone],
		openai.UserMessage(userMsg),
		openai.AssistantMessage(assistantReply),
	)

}
