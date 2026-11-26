package main

import "testing"

func TestRecordCardReview(t *testing.T) {
	card, err := createCard(database, Card{
		SetID:         1,
		English:       "Hello",
		Thai:          "สวัสดี",
		Transcription: "sawasdee",
	})
	if err != nil {
		t.Fatalf("failed to create test card: %v", err)
	}

	progress, err := recordCardReview(database, card.ID, true)
	if err != nil {
		t.Fatalf("failed to record correct answer: %v", err)
	}

	if progress.CorrectCount != 1 {
		t.Errorf(
			"expected correct count 1, got %d",
			progress.CorrectCount,
		)
	}

	if progress.IncorrectCount != 0 {
		t.Errorf(
			"expected incorrect count 0, got %d",
			progress.IncorrectCount,
		)
	}

	progress, err = recordCardReview(database, card.ID, false)
	if err != nil {
		t.Fatalf("failed to record incorrect answer: %v", err)
	}

	if progress.CorrectCount != 1 {
		t.Errorf(
			"expected correct count 1, got %d",
			progress.CorrectCount,
		)
	}

	if progress.IncorrectCount != 1 {
		t.Errorf(
			"expected incorrect count 1, got %d",
			progress.IncorrectCount,
		)
	}
}
