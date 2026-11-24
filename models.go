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
