package repositories

import (
	"context"
	"kolu/internal/models"
	"kolu/internal/types"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


func CreateLabel(pool *pgxpool.Pool, input *types.CreateLabelRequest, userId string) (*models.Label, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
		INSERT INTO labels (name, user_id)
		VALUES($1, $2)
		RETURNING id, name, user_id, created_at
		
	`
	var label models.Label

	err := pool.QueryRow(ctx, query, input.Name, userId).Scan(
		&label.ID,
		&label.Name,
		&label.UserID,
		&label.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &label, err

}
func GetLabel(pool *pgxpool.Pool, labelID string, userId string) (*models.Label, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
		SELECT id, name, user_id, created_at
		FROM labels
		WHERE id = $1 AND user_id = $2
		
	`
	var label models.Label

	err := pool.QueryRow(ctx, query, labelID, userId).Scan(
		&label.ID,
		&label.Name,
		&label.UserID,
		&label.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &label, err

}


func UpdateLabel(
	pool *pgxpool.Pool,
	labelID string,
	input *types.UpdateLabelRequest,
	userId string,
) (*models.Label, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE labels
		SET name = $1
		WHERE id = $2 AND user_id = $3
		RETURNING id, name, user_id, created_at
	`

	var label models.Label

	err := pool.QueryRow(
		ctx,
		query,
		input.Name,
		labelID,
		userId,
	).Scan(
		&label.ID,
		&label.Name,
		&label.UserID,
		&label.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &label, nil
}

func DeleteLabel(
	pool *pgxpool.Pool,
	labelID string,
	userId string,
) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		DELETE FROM labels
		WHERE id = $1 AND user_id = $2
	`

	tag, err := pool.Exec(ctx, query, labelID, userId)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
func GetUsersLabels(pool *pgxpool.Pool, user_id string) (*[]models.Label, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, user_id, created_at
		FROM labels
		WHERE user_id = $1
	`

	rows, err := pool.Query(ctx, query, user_id)

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
