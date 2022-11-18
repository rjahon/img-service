package postgres

import (
	"context"
	"fmt"

	"github.com/rjahon/img-service/config"
	"github.com/rjahon/img-service/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db  *pgxpool.Pool
	img storage.ImgRepoI
}

func NewPostgres(ctx context.Context, cfg config.Config) (storage.StorageI, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	))
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: pool,
	}, err
}

func (s *Store) CloseDB() {
	s.db.Close()
}

func (s *Store) Img() storage.ImgRepoI {
	if s.img == nil {
		s.img = NewImgRepo(s.db)
	}

	return s.img
}
