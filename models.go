package main

type CardSet struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Card struct {
	ID            int    `json:"id"`
	SetID         int    `json:"setId"`
	English       string `json:"english"`
	Thai          string `json:"thai"`
	Transcription string `json:"transcription"`
}

type CardProgress struct {
	CardID         int    `json:"cardId"`
	CorrectCount   int    `json:"correctCount"`
	IncorrectCount int    `json:"incorrectCount"`
	LastReviewedAt string `json:"lastReviewedAt"`
}

type SetProgress struct {
	SetID          int `json:"setId"`
	TotalCards     int `json:"totalCards"`
	StudiedCards   int `json:"studiedCards"`
	CorrectCount   int `json:"correctCount"`
	IncorrectCount int `json:"incorrectCount"`
}
