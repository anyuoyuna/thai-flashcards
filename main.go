package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func importCardsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(request.URL.Path, "/api/sets/")
	idStr = strings.TrimSuffix(idStr, "/cards/import")

	setID, err := strconv.Atoi(idStr)
	if err != nil || setID <= 0 {
		http.Error(writer, "invalid set ID", http.StatusBadRequest)
		return
	}

	if _, err := getSetByID(database, setID); errors.Is(err, sql.ErrNoRows) {
		http.Error(writer, "set not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(writer, "failed to find set", http.StatusInternalServerError)
		return
	}

	file, _, err := request.FormFile("file")
	if err != nil {
		http.Error(writer, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		http.Error(writer, "failed to read CSV header", http.StatusBadRequest)
		return
	}

	if len(header) != 3 ||
		header[0] != "english" ||
		header[1] != "thai" ||
		header[2] != "transcription" {
		http.Error(writer, "invalid CSV header", http.StatusBadRequest)
		return
	}

	var createdCards []Card
	rowNumber := 1
	for {
		record, err := reader.Read()
		rowNumber++

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			http.Error(
				writer,
				fmt.Sprintf("invalid CSV at row %d", rowNumber),
				http.StatusBadRequest,
			)
			return
		}

		if len(record) != 3 ||
			strings.TrimSpace(record[0]) == "" ||
			strings.TrimSpace(record[1]) == "" ||
			strings.TrimSpace(record[2]) == "" {
			http.Error(
				writer,
				fmt.Sprintf("invalid CSV row %d", rowNumber),
				http.StatusBadRequest,
			)
			return
		}

		card := Card{
			SetID:         setID,
			English:       record[0],
			Thai:          record[1],
			Transcription: record[2],
		}

		createdCard, err := createCard(database, card)
		if err != nil {
			http.Error(writer, "failed to create card", http.StatusInternalServerError)
			return
		}

		createdCards = append(createdCards, createdCard)
	}

	json.NewEncoder(writer).Encode(createdCards)
}

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

func cardReviewHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(request.URL.Path, "/api/cards/")
	idStr = strings.TrimSuffix(idStr, "/review")

	cardID, err := strconv.Atoi(idStr)
	if err != nil || cardID <= 0 {
		http.Error(writer, "invalid card ID", http.StatusBadRequest)
		return
	}

	if _, err := getCardByID(database, cardID); errors.Is(err, sql.ErrNoRows) {
		http.Error(writer, "card not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(writer, "failed to find card", http.StatusInternalServerError)
		return
	}

	var requestBody struct {
		Known bool `json:"known"`
	}

	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		http.Error(writer, "invalid request body", http.StatusBadRequest)
		return
	}

	progress, err := recordCardReview(database, cardID, requestBody.Known)
	if err != nil {
		http.Error(writer, "failed to save card review", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(progress)
}

func cardProgressHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(request.URL.Path, "/api/cards/")
	idStr = strings.TrimSuffix(idStr, "/progress")

	cardID, err := strconv.Atoi(idStr)
	if err != nil || cardID <= 0 {
		http.Error(writer, "invalid card ID", http.StatusBadRequest)
		return
	}

	if _, err := getCardByID(database, cardID); errors.Is(err, sql.ErrNoRows) {
		http.Error(writer, "card not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(writer, "failed to find card", http.StatusInternalServerError)
		return
	}

	progress, err := getCardProgress(database, cardID)

	if errors.Is(err, sql.ErrNoRows) {
		progress = CardProgress{
			CardID: cardID,
		}
	} else if err != nil {
		http.Error(writer, "failed to get card progress", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(progress)
}

func cardsRouter(writer http.ResponseWriter, request *http.Request) {
	if strings.HasSuffix(request.URL.Path, "/review") {
		cardReviewHandler(writer, request)
		return
	}
	if strings.HasSuffix(request.URL.Path, "/progress") {
		cardProgressHandler(writer, request)
		return
	}

	cardByIDHandler(writer, request)
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
	if strings.HasSuffix(request.URL.Path, "/cards/import") {
		importCardsHandler(writer, request)
		return
	}
	if strings.HasSuffix(request.URL.Path, "/cards") {
		cardsBySetHandler(writer, request)
		return
	}
	if strings.HasSuffix(request.URL.Path, "/progress") {
		setProgressHandler(writer, request)
		return
	}

	setByIDHandler(writer, request)
}

func setProgressHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(request.URL.Path, "/api/sets/")
	idStr = strings.TrimSuffix(idStr, "/progress")

	setID, err := strconv.Atoi(idStr)
	if err != nil || setID <= 0 {
		http.Error(writer, "invalid set ID", http.StatusBadRequest)
		return
	}

	if _, err := getSetByID(database, setID); errors.Is(err, sql.ErrNoRows) {
		http.Error(writer, "set not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(writer, "failed to find set", http.StatusInternalServerError)
		return
	}

	progress, err := getSetProgress(database, setID)
	if err != nil {
		http.Error(writer, "failed to get set progress", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(progress)
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
	http.HandleFunc("/api/cards/", cardsRouter)

	log.Println("server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
