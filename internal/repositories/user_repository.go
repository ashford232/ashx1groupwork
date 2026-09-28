package repositories

import (
	"context"
	"errors"
	"kolu/internal/config"
	errork "kolu/internal/error-k"
	"kolu/internal/models"
	"kolu/internal/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(pool *pgxpool.Pool, input *types.SignUpUserRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT into users (name, username, email, password)
		VALUES($1,$2,$3,$4)
		RETURNING id, name, username, email, created_at, updated_at
	`
	var user models.User
	err = pool.QueryRow(ctx, query, input.Name, input.Username, input.Email, hashedPassword).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, username, email,photo_url, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user models.User
	err := pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.PhotoURL,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {

		return nil, err
	}

	return &user, nil
}

func GetUserByID(pool *pgxpool.Pool, id string) (*models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, username, email,photo_url, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user models.User
	err := pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.PhotoURL,
		
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {

		return nil, err
	}

	return &user, nil
}

func SignIn(pool *pgxpool.Pool, input *types.SingInUserRequest, cfg *config.Config) (*types.SignInUserResponse, error) {

	existingUser, err := GetUserByEmail(pool, input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errork.ErrInvalidCredentials
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(input.Password))
	if err != nil {
		return nil, errork.ErrInvalidCredentials
	}

	user, err := GetUserByID(pool, existingUser.ID)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {

		return nil, err
	}

	return &types.SignInUserResponse{User: user, Token: tokenStr}, nil
}

func SetupProfile(pool *pgxpool.Pool, input *types.SetupUserRequest, id uuid.UUID) (*models.User, error) {
	query := `
		UPDATE users 
		SET username = $1, photo_url = COALESCE(NULLIF($2, ''), photo_url), updated_at = NOW()
		WHERE id = $3
		RETURNING id, email, username, photo_url, created_at, updated_at
	`

	var user models.User
	err := pool.QueryRow(
		context.Background(),
		query,
		input.Username,
		input.PhotoURL,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PhotoURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
