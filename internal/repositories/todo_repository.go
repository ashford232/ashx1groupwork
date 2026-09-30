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

func CreateTodo(pool *pgxpool.Pool, input *types.CreateTodoRequest, user_id string) (*models.Todo, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO todos (
		id,
		title,
		description,
		is_completed,
		user_id,
		due_at,
		priority
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (id)
	DO UPDATE SET
		title = EXCLUDED.title,
		description = EXCLUDED.description,
		is_completed = EXCLUDED.is_completed,
		due_at = EXCLUDED.due_at,
		priority = EXCLUDED.priority,
		updated_at = NOW()
	RETURNING
		id,
		title,
		description,
		is_completed,
		user_id,
		is_deleted,
		created_at,
		updated_at,
		due_at,
		priority
`

	var todo models.Todo

	err := pool.QueryRow(
		ctx,
		query,
		input.ID,
		input.Title,
		input.Description,
		input.Completed,
		user_id,
		input.DueAt,
		input.Priority,
	).Scan(
		&todo.ID,          // id
		&todo.Title,       // title
		&todo.Description, // description
		&todo.IsCompleted, // is_completed
		&todo.UserId,      // user_id
		&todo.IsDeleted,   // is_deleted
		&todo.CreatedAt,   // created_at
		&todo.UpdatedAt,   // updated_at
		&todo.DueAt,       // due_at
		&todo.Priority,    // priority
	)
	if err != nil {
		return nil, err
	}
	if err := attachTodoLabels(pool, &todo); err != nil {
		return nil, err
	}

	return &todo, nil
}
func GetAllTodos(pool *pgxpool.Pool, user_id string) ([]models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, description, is_completed, user_id, is_deleted, created_at, updated_at, due_at, priority
		from todos 
		WHERE user_id = $1 AND is_deleted = FALSE
		ORDER BY updated_at DESC
	`

	rows, err := pool.Query(ctx, query, user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []models.Todo

	for rows.Next() {
		var todo models.Todo

		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.IsCompleted,
			&todo.UserId,
			&todo.IsDeleted,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.DueAt,
			&todo.Priority,
		)
		if err != nil {
			return nil, err
		}
		if err := attachTodoLabels(pool, &todo); err != nil {
			return nil, err
		}

		todos = append(todos, todo)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil

}
func GetAllTrashedTodos(pool *pgxpool.Pool, user_id string) ([]models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, description, is_completed, user_id, is_deleted, created_at, updated_at, due_at, priority
		from todos 
		WHERE user_id = $1 AND is_deleted = TRUE
		ORDER BY updated_at DESC
	`

	rows, err := pool.Query(ctx, query, user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []models.Todo

	for rows.Next() {
		var todo models.Todo

		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.IsCompleted,
			&todo.UserId,
			&todo.IsDeleted,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.DueAt,
			&todo.Priority,
		)
		if err != nil {
			return nil, err
		}
		if err := attachTodoLabels(pool, &todo); err != nil {
			return nil, err
		}

		todos = append(todos, todo)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil

}

func GetTodoByID(pool *pgxpool.Pool, id uuid.UUID, user_id string) (*models.Todo, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, description, is_completed, user_id, is_deleted, created_at, updated_at, due_at, priority
		FROM todos
		WHERE id = $1 AND user_id = $2
	`
	var todo models.Todo
	err := pool.QueryRow(ctx, query, id, user_id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.IsCompleted,
		&todo.UserId,
		&todo.IsDeleted,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.DueAt,
		&todo.Priority,
	)
	if err != nil {
		return nil, err
	}

	if err := attachTodoLabels(pool, &todo); err != nil {
		return nil, err
	}

	return &todo, nil
}

func UpdateTodo(pool *pgxpool.Pool, input *types.UpdateTodoRequest, id uuid.UUID, user_id string) (*models.Todo, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE todos 
		SET title = $1, description = $2, is_completed = $3, due_at = $4, priority = $5, is_deleted = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7 AND user_id = $8
		RETURNING id, title, description, is_completed, user_id, created_at, updated_at, due_at, priority

	`

	var todo models.Todo
	err := pool.QueryRow(
		ctx,
		query,
		input.Title,
		input.Description,
		input.IsCompleted,
		input.DueAt,
		input.Priority,
		input.IsDeleted,
		id,
		user_id,
	).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.IsCompleted,
		&todo.UserId,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.DueAt,
		&todo.Priority,
		&todo.IsDeleted,
	)
	if err != nil {
		return nil, err
	}
	if err := attachTodoLabels(pool, &todo); err != nil {
		return nil, err
	}

	return &todo, nil

}

func DeleteTodoById(pool *pgxpool.Pool, id string, user_id string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
	DELETE FROM todos
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

func DeleteAllTodos(pool *pgxpool.Pool, user_id string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
	DELETE FROM todos
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

func attachTodoLabels(pool *pgxpool.Pool, todo *models.Todo) error {
	labels, err := GetTodoLabels(pool, todo.ID)
	if err != nil {
		return err
	}

	todo.Labels = *labels

	return nil
}

func GetTodoLabels(pool *pgxpool.Pool, todo_id uuid.UUID) (*[]models.Label, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT l.id, l.name, l.user_id, l.created_at
		FROM labels l
		JOIN todo_labels lt ON lt.label_id = l.id
		WHERE lt.todo_id = $1
	`

	rows, err := pool.Query(ctx, query, todo_id)
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
func AddTodoLabels(pool *pgxpool.Pool, input *types.AddTodoLabelRequest, userId string) (*models.Label, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
		INSERT INTO todo_labels (todo_id, label_id)
		VALUES($1, $2)
		RETURNING todo_id, label_id
	`
	var labelRequest types.AddTodoLabelRequest

	err := pool.QueryRow(ctx, query, input.TodoID, input.LabelID).Scan(
		&labelRequest.TodoID,
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
