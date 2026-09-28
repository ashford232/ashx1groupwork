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

func CreateDoc(pool *pgxpool.Pool, input *types.CreateDocRequest, user_id string) (*models.Doc, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO docs(id, name, object_id, mime_type, size, user_id,)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, object_id, mime_type, size, user_id,is_deleted, created_at, updated_at
	`

	var doc models.Doc

	err := pool.QueryRow(
		ctx, 
		query, 
		input.ID, 
		input.Name, 
		input.ObjectKey, 
		input.MimeType, 
		input.Size, 
		user_id).Scan(
			&doc.ID, 
			&doc.Name, 
			&doc.ObjectKey, 
			&doc.MimeType, 
			&doc.Size, 
			&doc.UserId, 
			&doc.IsDeleted, 
			&doc.CreatedAt, 
			&doc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err := attachDocLabels(pool, &doc); err != nil {
		return nil, err
	}

	return &doc, nil
}
func GetAllDocs(pool *pgxpool.Pool, user_id string) ([]models.Doc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, object_id, mime_type, size, user_id, is_deleted, created_at, updated_at
		from docs 
		WHERE user_id = $1 AND is_deleted = FALSE
		ORDER BY updated_at DESC
	`

	rows, err := pool.Query(ctx, query, user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var docs []models.Doc

	for rows.Next() {
		var doc models.Doc

		err = rows.Scan(
			&doc.ID,
			&doc.Name,
			&doc.ObjectKey,
			&doc.MimeType,
			&doc.Size,
			&doc.UserId,
			&doc.IsDeleted,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if err := attachDocLabels(pool, &doc); err != nil {
			return nil, err
		}

		docs = append(docs, doc)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return docs, nil

}

func GetAllTrashedDocs(pool *pgxpool.Pool, user_id string) ([]models.Doc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, object_id, mime_type, size, is_deleted, user_id, created_at, updated_at
		from docs 
		WHERE user_id = $1 AND is_deleted = TRUE
		ORDER BY updated_at DESC
	`

	rows, err := pool.Query(ctx, query, user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var docs []models.Doc

	for rows.Next() {
		var doc models.Doc

		err = rows.Scan(
			&doc.ID,
			&doc.Name,
			&doc.ObjectKey,
			&doc.MimeType,
			&doc.Size,
			&doc.UserId,
			&doc.IsDeleted,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if err := attachDocLabels(pool, &doc); err != nil {
			return nil, err
		}

		docs = append(docs, doc)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return docs, nil

}

func GetDocByID(pool *pgxpool.Pool, id uuid.UUID, user_id string) (*models.Doc, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, object_id, mime_type, size, user_id, is_deleted, created_at, updated_at
		FROM docs
		WHERE id = $1 AND user_id = $2
	`
	var doc models.Doc
	err := pool.QueryRow(ctx, query, id, user_id).Scan(
		&doc.ID,
		&doc.Name,
		&doc.ObjectKey,
		&doc.MimeType,
		&doc.Size,
		&doc.UserId,
		&doc.IsDeleted,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := attachDocLabels(pool, &doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

func UpdateDoc(pool *pgxpool.Pool, input *types.UpdateDocRequest, id uuid.UUID, user_id string) (*models.Doc, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE notes
		SET name = $1, is_deleted = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $4
		RETURNING id, name, object_id, mime_type, size, user_id, is_deleted, created_at, updated_at

	`

	var doc models.Doc
	err := pool.QueryRow(ctx, query, input.Name,input.IsDeleted, id, user_id).Scan(
		&doc.ID,
		&doc.Name,
		&doc.ObjectKey,
		&doc.MimeType,
		&doc.Size,
		&doc.UserId,
		&doc.ObjectKey,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := attachDocLabels(pool, &doc); err != nil {
		return nil, err
	}

	return &doc, nil

}

func DeleteDocById(pool *pgxpool.Pool, id string, user_id string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
	DELETE FROM docs
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

func DeleteAllDocs(pool *pgxpool.Pool, user_id string) error {

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

func attachDocLabels(pool *pgxpool.Pool, doc *models.Doc) error {
	labels, err := GetNoteLabels(pool, doc.ID)
	if err != nil {
		return err
	}

	doc.Labels = *labels

	return nil
}

func GetDocLabels(pool *pgxpool.Pool, note_id uuid.UUID) (*[]models.Label, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT l.id, l.name, l.user_id, l.created_at
		FROM labels l
		JOIN doc_labels ld ON ld.label_id = l.id
		WHERE ld.todo_id = $1
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
func AddDocLabels(pool *pgxpool.Pool, input *types.AddDocLabelRequest, userId string) (*models.Label, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
		INSERT INTO doc_labels (doc_id, label_id)
		VALUES($1, $2)
		RETURNING doc_id, label_id
	`
	var labelRequest types.AddDocLabelRequest

	err := pool.QueryRow(ctx, query, input.DocID, input.LabelID).Scan(
		&labelRequest.DocID,
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
