package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var database *sql.DB

func openDatabase() (*sql.DB, error) {
	database, err := sql.Open("sqlite", "thai-flashcards.db")
	if err != nil {
		return nil, err
	}

	if err := database.Ping(); err != nil {
		database.Close()
		return nil, err
	}

	_, err = database.Exec(`
        CREATE TABLE IF NOT EXISTS card_sets (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL
        )
    `)
	if err != nil {
		database.Close()
		return nil, err
	}

	_, err = database.Exec(`
        CREATE TABLE IF NOT EXISTS cards (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            set_id INTEGER NOT NULL,
            english TEXT NOT NULL,
            thai TEXT NOT NULL,
            transcription TEXT NOT NULL,
            FOREIGN KEY (set_id) REFERENCES card_sets(id)
        )
    `)
	if err != nil {
		database.Close()
		return nil, err
	}

	_, err = database.Exec(`
        CREATE TABLE IF NOT EXISTS card_progress (
			card_id INTEGER PRIMARY KEY,
			correct_count INTEGER NOT NULL DEFAULT 0,
			incorrect_count INTEGER NOT NULL DEFAULT 0,
			last_reviewed_at TEXT,
			FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE
		)
    `)
	if err != nil {
		database.Close()
		return nil, err
	}

	return database, nil
}

func getSets(database *sql.DB) ([]CardSet, error) {
	rows, err := database.Query(`
        SELECT id, title
        FROM card_sets
        ORDER BY id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []CardSet

	for rows.Next() {
		var set CardSet

		if err := rows.Scan(&set.ID, &set.Title); err != nil {
			return nil, err
		}

		result = append(result, set)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func getSetByID(database *sql.DB, id int) (CardSet, error) {
	var set CardSet

	err := database.QueryRow(`
        SELECT id, title
        FROM card_sets
        WHERE id = ?
    `, id).Scan(&set.ID, &set.Title)

	if err != nil {
		return CardSet{}, err
	}

	return set, nil
}

func updateSetTitle(database *sql.DB, id int, title string) (CardSet, error) {
	result, err := database.Exec(`
        UPDATE card_sets
        SET title = ?
        WHERE id = ?
    `, title, id)

	if err != nil {
		return CardSet{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return CardSet{}, err
	}

	if rowsAffected == 0 {
		return CardSet{}, sql.ErrNoRows
	}

	return getSetByID(database, id)
}

func deleteSet(database *sql.DB, id int) error {
	result, err := database.Exec(`
        DELETE FROM card_sets
        WHERE id = ?
    `, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func getCardsBySetID(database *sql.DB, setID int) ([]Card, error) {
	rows, err := database.Query(`
        SELECT id, set_id, english, thai, transcription
        FROM cards
        WHERE set_id = ?
        ORDER BY id
    `, setID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Card

	for rows.Next() {
		var card Card

		if err := rows.Scan(
			&card.ID,
			&card.SetID,
			&card.English,
			&card.Thai,
			&card.Transcription,
		); err != nil {
			return nil, err
		}

		result = append(result, card)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func createCard(database *sql.DB, card Card) (Card, error) {
	result, err := database.Exec(`
        INSERT INTO cards (set_id, english, thai, transcription)
        VALUES (?, ?, ?, ?)
    `, card.SetID, card.English, card.Thai, card.Transcription)
	if err != nil {
		return Card{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Card{}, err
	}

	card.ID = int(id)

	return card, nil
}

func getCardByID(database *sql.DB, id int) (Card, error) {
	var card Card

	err := database.QueryRow(`
        SELECT id, set_id, english, thai, transcription
        FROM cards
        WHERE id = ?
    `, id).Scan(
		&card.ID,
		&card.SetID,
		&card.English,
		&card.Thai,
		&card.Transcription,
	)

	if err != nil {
		return Card{}, err
	}

	return card, nil
}

func updateCard(database *sql.DB, id int, card Card) (Card, error) {
	result, err := database.Exec(`
        UPDATE cards
        SET english = ?, thai = ?, transcription = ?
        WHERE id = ?
    `, card.English, card.Thai, card.Transcription, id)
	if err != nil {
		return Card{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Card{}, err
	}

	if rowsAffected == 0 {
		return Card{}, sql.ErrNoRows
	}

	return getCardByID(database, id)
}

func deleteCard(database *sql.DB, id int) error {
	result, err := database.Exec(`
        DELETE FROM cards
        WHERE id = ?
    `, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func getCardProgress(database *sql.DB, cardID int) (CardProgress, error) {
	var progress CardProgress

	err := database.QueryRow(`
        SELECT card_id, correct_count, incorrect_count, last_reviewed_at
        FROM card_progress
        WHERE card_id = ?
    `, cardID).Scan(
		&progress.CardID,
		&progress.CorrectCount,
		&progress.IncorrectCount,
		&progress.LastReviewedAt,
	)

	if err != nil {
		return CardProgress{}, err
	}

	return progress, nil
}

func recordCardReview(
	database *sql.DB,
	cardID int,
	known bool,
) (CardProgress, error) {
	correctIncrement := 0
	incorrectIncrement := 1

	if known {
		correctIncrement = 1
		incorrectIncrement = 0
	}

	_, err := database.Exec(`
        INSERT INTO card_progress (
            card_id,
            correct_count,
            incorrect_count,
            last_reviewed_at
        )
        VALUES (?, ?, ?, CURRENT_TIMESTAMP)
        ON CONFLICT(card_id) DO UPDATE SET
            correct_count = correct_count + excluded.correct_count,
            incorrect_count = incorrect_count + excluded.incorrect_count,
            last_reviewed_at = CURRENT_TIMESTAMP
    `, cardID, correctIncrement, incorrectIncrement)

	if err != nil {
		return CardProgress{}, err
	}

	return getCardProgress(database, cardID)
}

func getSetProgress(database *sql.DB, setID int) (SetProgress, error) {
	var progress SetProgress

	err := database.QueryRow(`
        SELECT
            COUNT(cards.id),
            COUNT(card_progress.card_id),
            COALESCE(SUM(card_progress.correct_count), 0),
            COALESCE(SUM(card_progress.incorrect_count), 0)
        FROM cards
        LEFT JOIN card_progress
            ON card_progress.card_id = cards.id
        WHERE cards.set_id = ?
    `, setID).Scan(
		&progress.TotalCards,
		&progress.StudiedCards,
		&progress.CorrectCount,
		&progress.IncorrectCount,
	)

	if err != nil {
		return SetProgress{}, err
	}

	progress.SetID = setID

	return progress, nil
}
