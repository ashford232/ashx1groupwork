package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = pool.Ping(ctx)

	if err != nil {
		log.Println(err)
		pool.Close()
		return nil, err
	}
	log.Println("PostgreSQL connected")
	return pool, nil

}
