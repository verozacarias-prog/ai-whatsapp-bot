package main

import (
	"sync"
	"time"

	openai "github.com/openai/openai-go/v3"
)

type ConversationEntry struct {
	Messages []openai.ChatCompletionMessageParamUnion
	LastDate time.Time
}

var (
	conversationHistory = make(map[string]ConversationEntry)
	historyMu           sync.RWMutex
)

func GetHistory(phone string) ConversationEntry {
	historyMu.RLock()
	defer historyMu.RUnlock()
	
	entry := conversationHistory[phone]
	if !entry.LastDate.IsZero() && time.Now().Day() != entry.LastDate.Day() {
		delete(conversationHistory, phone)
		return ConversationEntry{}
	}
	return entry
}

func AddToHistory(phone, userMsg, assistantReply string) {
	historyMu.Lock()
	defer historyMu.Unlock()

	history := conversationHistory[phone]

	// Append new messages to the history slice
	history.Messages = append(history.Messages,
		openai.UserMessage(userMsg),
		openai.AssistantMessage(assistantReply),
	)
	history.LastDate = time.Now()
	conversationHistory[phone] = history

}
