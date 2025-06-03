package repository

import (
	"context"
	"javaneseivankov/url-shortener/internal/repository/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)


type ShortLinkRepositoryDB struct {
    store map[string]model.ShortLink
	 db *pgxpool.Pool
}

func NewShortLinkRepositoryDB(db *pgxpool.Pool) IShortLinkRepository {
    return &ShortLinkRepositoryDB{
        store: make(map[string]model.ShortLink),
		  db: db,
    }
}

func (s *ShortLinkRepositoryDB) CreateRedirectLink(ctx context.Context, shortLink model.ShortLink) (*model.ShortLink, error) {
	query := `INSERT INTO shortlinks()`
	s.db.Exec(ctx, query, shortLink.ShortName)
	panic("Not Implemented")
}

func (s *ShortLinkRepositoryDB) GetRedirectLink(ctx context.Context, shortName string) (*model.ShortLink, error) {
	panic("Not Implemented")
}

func (s *ShortLinkRepositoryDB) DeleteRedirectLink(ctx context.Context, shortName string, userId uuid.UUID) error {
	panic("Not Implemented")
}

func (s *ShortLinkRepositoryDB) EditShortLink(ctx context.Context, shortName string, url string, userId uuid.UUID) (*model.ShortLink, error) {
	panic("Not Implemented")
}
