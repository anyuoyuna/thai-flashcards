package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func healthHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status": "ok",
	}

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("write response: %v", err)
	}
}

func setsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	switch request.Method {
	case http.MethodGet:
		databaseSets, err := getSets(database)
		if err != nil {
			http.Error(writer, "failed to get sets", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(writer).Encode(databaseSets)
	case http.MethodPost:
		var set CardSet
		if err := json.NewDecoder(request.Body).Decode(&set); err != nil {
			http.Error(writer, "invalid request body", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(set.Title) == "" {
			http.Error(writer, "title is required", http.StatusBadRequest)
			return
		}

		result, err := database.Exec(
			"INSERT INTO card_sets (title) VALUES (?)",
			set.Title,
		)
		if err != nil {
			http.Error(writer, "failed to create set", http.StatusInternalServerError)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			http.Error(writer, "failed to get created set ID", http.StatusInternalServerError)
			return
		}
		set.ID = int(id)

		writer.WriteHeader(http.StatusCreated)
		json.NewEncoder(writer).Encode(set)
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func setByIDHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	if request.Method != http.MethodGet &&
		request.Method != http.MethodPut &&
		request.Method != http.MethodDelete {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(request.URL.Path, "/api/sets/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(writer, "invalid set ID", http.StatusBadRequest)
		return
	}

	switch request.Method {
	case http.MethodGet:
		set, err := getSetByID(database, id)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(writer, "set not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(writer, "failed to get set", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(writer).Encode(set)
	case http.MethodPut:
		var set CardSet

		if err := json.NewDecoder(request.Body).Decode(&set); err != nil {
			http.Error(writer, "invalid request body", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(set.Title) == "" {
			http.Error(writer, "title is required", http.StatusBadRequest)
			return
		}

		updatedSet, err := updateSetTitle(database, id, set.Title)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(writer, "set not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(writer, "failed to update set", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(writer).Encode(updatedSet)
	case http.MethodDelete:
		err := deleteSet(database, id)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(writer, "set not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(writer, "failed to delete set", http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusNoContent)
	}
}

func cardsBySetHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(request.URL.Path, "/api/sets/")
	path = strings.TrimSuffix(path, "/cards")

	setID, err := strconv.Atoi(path)
	if err != nil || setID <= 0 {
		http.Error(writer, "invalid set ID", http.StatusBadRequest)
		return
	}

	switch request.Method {
	case http.MethodGet:
		cards, err := getCardsBySetID(database, setID)
		if err != nil {
			http.Error(writer, "failed to get cards", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(writer).Encode(cards)

	case http.MethodPost:
		var card Card

		if err := json.NewDecoder(request.Body).Decode(&card); err != nil {
			http.Error(writer, "invalid request body", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(card.English) == "" ||
			strings.TrimSpace(card.Thai) == "" ||
			strings.TrimSpace(card.Transcription) == "" {
			http.Error(writer, "all card fields are required", http.StatusBadRequest)
			return
		}

		if _, err := getSetByID(database, setID); errors.Is(err, sql.ErrNoRows) {
			http.Error(writer, "set not found", http.StatusNotFound)
			return
		}

		card.SetID = setID

		createdCard, err := createCard(database, card)
		if err != nil {
			http.Error(writer, "failed to create card", http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusCreated)
		json.NewEncoder(writer).Encode(createdCard)

	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func cardByIDHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	if request.Method != http.MethodGet &&
		request.Method != http.MethodPut &&
		request.Method != http.MethodDelete {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(request.URL.Path, "/api/cards/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(writer, "invalid card ID", http.StatusBadRequest)
		return
	}

	switch request.Method {
	case http.MethodGet:
		card, err := getCardByID(database, id)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(writer, "card not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(writer, "failed to get card", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(writer).Encode(card)

	case http.MethodPut:
		var card Card

		if err := json.NewDecoder(request.Body).Decode(&card); err != nil {
			http.Error(writer, "invalid request body", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(card.English) == "" ||
			strings.TrimSpace(card.Thai) == "" ||
			strings.TrimSpace(card.Transcription) == "" {
			http.Error(writer, "all card fields are required", http.StatusBadRequest)
			return
		}

		updatedCard, err := updateCard(database, id, card)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(writer, "card not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(writer, "failed to update card", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(writer).Encode(updatedCard)

	case http.MethodDelete:
		err := deleteCard(database, id)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(writer, "card not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(writer, "failed to delete card", http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusNoContent)
	}
}

func setsRouter(writer http.ResponseWriter, request *http.Request) {
	if strings.HasSuffix(request.URL.Path, "/cards") {
		cardsBySetHandler(writer, request)
		return
	}

	setByIDHandler(writer, request)
}

func main() {
	var err error

	database, err = openDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	http.HandleFunc("/api/health", healthHandler)
	http.HandleFunc("/api/sets", setsHandler)
	http.HandleFunc("/api/sets/", setsRouter)
	http.HandleFunc("/api/cards/", cardByIDHandler)

	log.Println("server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
