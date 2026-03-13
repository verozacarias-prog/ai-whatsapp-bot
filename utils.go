package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func WriteCSVLog(message string, intent, confidence string, tokensUsed int64) {
	csvFile, err := os.OpenFile("logs.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error al abrir el archivo CSV: %v", err)
		return
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	csvWriter.Write([]string{
		time.Now().Format(time.RFC3339),
		message,
		intent,
		confidence,
		fmt.Sprintf("%d", tokensUsed)})

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		log.Printf("error al escribir en el archivo CSV: %v", err)
	}
}
