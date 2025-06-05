package repository

import (
	"context"
	"errors"
	"javaneseivankov/url-shortener/internal/errx"
	"javaneseivankov/url-shortener/internal/repository/model"
	"javaneseivankov/url-shortener/pkg/pgerror"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type ShortLinkRepositoryDB struct {
	 db *pgxpool.Pool
}

func NewShortLinkRepositoryDB(db *pgxpool.Pool) IShortLinkRepository {
    return &ShortLinkRepositoryDB{
		  db: db,
    }
}

func (s *ShortLinkRepositoryDB) CreateRedirectLink(ctx context.Context, shortLink model.ShortLink) (*model.ShortLink, error) {
	query := `
	INSERT INTO shortlinks(id, owner_id, short_name, original_url)
	VALUES ($1, $2, $3, $4);
	`
	_, err := s.db.Exec(ctx, query, shortLink.Id, shortLink.UserId, shortLink.ShortName, shortLink.OriginalUrl);
	err = pgerror.NewPgErrHandler().
						AddPgErr(
							pgerror.UNIQUE_VIOLATION_CODE, 
							"shortlinks_short_name_key",
							errx.ErrShortLinkAlreadyExists,
						).
                  AddPgErr(
                     pgerror.FK_VIOLATION_CODE,
                     "shortlinks_owner_id_fkey",
                     errx.ErrUserIdDoesntExist,
                  ).
						Handle(err)
	if err != nil {
		return nil, err;
	}
	return &shortLink, err
}

func (s *ShortLinkRepositoryDB) GetRedirectLink(ctx context.Context, shortName string) (*model.ShortLink, error) {
	query := `
	SELECT id, short_name, original_url, created_at, owner_id
	FROM shortlinks
	WHERE short_name = $1;
	`
	var sLink model.ShortLink
	err := s.db.QueryRow(ctx, query, shortName).Scan(&sLink.Id, &sLink.ShortName, &sLink.OriginalUrl, &sLink.CreatedAt, &sLink.UserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errx.ErrShortLinkNotFound
		}
		return nil, err
	}
	return &sLink, nil
}

func (s *ShortLinkRepositoryDB) DeleteRedirectLink(ctx context.Context, shortName string, userId uuid.UUID) error {
	query := `
	UPDATE shortlinks 
	SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
	WHERE short_name = $1 AND owner_id = $2;
	`

	res, err := s.db.Exec(ctx, query, shortName, userId);
	if err != nil {
		return err
	}
	
	if res.RowsAffected() == 0 {
		return errx.ErrShortLinkNotFound
	}

	return nil
}

func (s *ShortLinkRepositoryDB) EditShortLink(ctx context.Context, shortName string, url string, userId uuid.UUID) (*model.ShortLink, error) {
	query := `
	UPDATE shortlinks
	SET original_url = $2,
		updated_at = CURRENT_TIMESTAMP
	WHERE short_name = $1 AND owner_id = $3
	RETURNING id, short_name, original_url, created_at, owner_id;
	`

	var sLink model.ShortLink
	err := s.db.QueryRow(ctx, query, shortName, url, userId).Scan(
		&sLink.Id,
		&sLink.ShortName,
		&sLink.OriginalUrl,
		&sLink.CreatedAt,
		&sLink.UserId,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errx.ErrShortLinkNotFound
		}
		return nil, err
	}

	return &sLink, nil
}
