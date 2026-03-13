package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error cargando .env")
	}

	if err := LoadBusinessConfig("business.yaml"); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/classify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
			return
		}

		var req ClassifyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "body inválido", http.StatusBadRequest)
			return
		}
		result, err := Classify(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/respond", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
			return
		}

		var req RespondRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "body inválido", http.StatusBadRequest)
			return
		}

		result, err := Respond(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	log.Println("Servidor corriendo en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
