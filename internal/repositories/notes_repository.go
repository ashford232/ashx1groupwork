package repositories

import (
	"context"
	"kolu/internal/models"
	"kolu/internal/types"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateNote(pool *pgxpool.Pool, input *types.CreateNoteRequest, user_id string) (*models.Note, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO notes(id, title, content, user_id,)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, content, , user_id, is_deleted, created_at, updated_at
	`

	var note models.Note

	err := pool.QueryRow(
		ctx,
		query,
		input.ID, 
		input.Title, 
		input.Content, 
		user_id).Scan(
			&note.ID, 
			&note.Title, 
			&note.Content, 
			&note.UserId, 
			&note.IsDeleted, 
			&note.CreatedAt, 
			&note.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err := attachNoteLabels(pool, &note); err != nil {
		return nil, err
	}

	return &note, nil
}
func GetAllNotes(pool *pgxpool.Pool, user_id string) ([]models.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, content, user_id, is_deleted, created_at, updated_at
		from notes 
		WHERE user_id = $1 AND is_deleted = FALSE
		ORDER BY updated_at DESC
	`

	rows, err := pool.Query(ctx, query, user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notes []models.Note

	for rows.Next() {
		var note models.Note

		err = rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.UserId,
			&note.IsDeleted,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if err := attachNoteLabels(pool, &note); err != nil {
			return nil, err
		}

		notes = append(notes, note)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil

}
func GetAllTrashedNotes(pool *pgxpool.Pool, user_id string) ([]models.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, content, user_id, is_deleted, created_at, updated_at
		from notes 
		WHERE user_id = $1 AND is_deleted = TRUE
		ORDER BY updated_at DESC
	`

	rows, err := pool.Query(ctx, query, user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notes []models.Note

	for rows.Next() {
		var note models.Note

		err = rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.UserId,
			&note.IsDeleted,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if err := attachNoteLabels(pool, &note); err != nil {
			return nil, err
		}

		notes = append(notes, note)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil

}

func GetNoteByID(pool *pgxpool.Pool, id uuid.UUID, user_id string) (*models.Note, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, content, user_id, is_deleted, created_at, updated_at
		FROM notes
		WHERE id = $1 AND user_id = $2
	`
	var note models.Note
	err := pool.QueryRow(ctx, query, id, user_id).Scan(
		&note.ID,
		&note.Title,
		&note.UserId,
		&note.IsDeleted,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := attachNoteLabels(pool, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

func UpdateNote(pool *pgxpool.Pool, input *types.UpdateNoteRequest, id uuid.UUID, user_id string) (*models.Note, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE notes
		SET title=$1, content = $2, is_deleted= $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND user_id = $5
		RETURNING id, title, content, user_id,is_deleted, created_at, updated_at, 

	`

	var note models.Note
	err := pool.QueryRow(ctx, query, input.Title, input.Content,input.IsDeleted, id, user_id).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.UserId,
		&note.IsDeleted,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := attachNoteLabels(pool, &note); err != nil {
		return nil, err
	}

	return &note, nil

}

func DeleteNoteById(pool *pgxpool.Pool, id string, user_id string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
	DELETE FROM notes
	WHERE id = $1 AND user_id = $2
	`

	tag, err := pool.Exec(ctx, query, id, user_id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func DeleteAllNotes(pool *pgxpool.Pool, user_id string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
	DELETE FROM notes
	WHERE user_id = $1
	`

	tag, err := pool.Exec(ctx, query, user_id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func attachNoteLabels(pool *pgxpool.Pool, note *models.Note) error {
	labels, err := GetNoteLabels(pool, note.ID)
	if err != nil {
		return err
	}

	note.Labels = *labels

	return nil
}

func GetNoteLabels(pool *pgxpool.Pool, note_id uuid.UUID) (*[]models.Label, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT l.id, l.name, l.user_id, l.created_at
		FROM labels l
		JOIN note_labels ln ON ln.label_id = l.id
		WHERE ln.todo_id = $1
	`

	rows, err := pool.Query(ctx, query, note_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []models.Label

	for rows.Next() {
		var label models.Label

		err := rows.Scan(
			&label.ID,
			&label.Name,
			&label.UserID,
			&label.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		labels = append(labels, label)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &labels, nil
}
func AddNoteLabels(pool *pgxpool.Pool, input *types.AddNoteLabelRequest, userId string) (*models.Label, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
		INSERT INTO note_labels (note_id, label_id)
		VALUES($1, $2)
		RETURNING note_id, label_id
	`
	var labelRequest types.AddNoteLabelRequest

	err := pool.QueryRow(ctx, query, input.NoteID, input.LabelID).Scan(
		&labelRequest.NoteID,
		&labelRequest.LabelID,
	)
	if err != nil {
		return nil, err
	}
	label, err := GetLabel(pool, labelRequest.LabelID, userId)

	if err != nil {
		return nil, err
	}

	return label, err

}
