package repository

import (
	"context"
	"errors"
	"javaneseivankov/url-shortener/internal/errx"
	"javaneseivankov/url-shortener/internal/repository/model"
	"javaneseivankov/url-shortener/pkg/logger"
	"javaneseivankov/url-shortener/pkg/pgerror"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryDB struct {
    db *pgxpool.Pool
}

func NewUserRepositoryDB(db *pgxpool.Pool) IUserRepository {
    return &UserRepositoryDB{
        db: db,
    }
}

func (r *UserRepositoryDB) CreateUser(ctx context.Context, user model.User) error {
    logger.Info("UserRepositoryDB.CreateUser: creating user", "userID", user.ID, "email", user.Email)
    query := `
    INSERT INTO users(user_id, email, password)
    VALUES ($1, $2, $3)
    `
    _, err := r.db.Exec(ctx, query, user.ID, user.Email, user.Password)
    err = pgerror.NewPgErrHandler().
        AddPgErr(
            pgerror.UNIQUE_VIOLATION_CODE,
            "users_email_key",
            errx.ErrUserEmailAlreadyExists,
        ).
        AddPgErr(
            pgerror.UNIQUE_VIOLATION_CODE,
            "users_user_id_key",
            errx.ErrUserIdAlreadyExists,
        ).
        Handle(err)
    if err != nil {
        logger.Error("UserRepositoryDB.CreateUser: failed to create user", "userID", user.ID, "email", user.Email, "error", err)
        return err
    }

    logger.Info("UserRepositoryDB.CreateUser: user created successfully", "userID", user.ID, "email", user.Email)
    return nil
}

func (r *UserRepositoryDB) UpdateUser(ctx context.Context, user model.User) error {
    logger.Info("UserRepositoryDB.UpdateUser: updating user", "userID", user.ID, "email", user.Email)
    query := `
    UPDATE users
    SET email = $2,
    password = $3
    WHERE user_id = $1;
    `

    _, err := r.db.Exec(ctx, query, user.ID, user.Email, user.Password)
    err = pgerror.NewPgErrHandler().
        AddPgErr(
            pgerror.UNIQUE_VIOLATION_CODE,
            "users_email_key",
            errx.ErrUserEmailAlreadyExists,
        ).
        Handle(err)
    if err != nil {
        logger.Error("UserRepositoryDB.UpdateUser: failed to update user", "userID", user.ID, "email", user.Email, "error", err)
        return err
    }

    logger.Info("UserRepositoryDB.UpdateUser: user updated successfully", "userID", user.ID, "email", user.Email)
    return nil
}

func (r *UserRepositoryDB) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
    logger.Info("UserRepositoryDB.GetUserByEmail: retrieving user by email", "email", email)
    query := `
    SELECT user_id, email, password
    FROM users
    WHERE email = $1
    `

    var user model.User
    err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            logger.Error("UserRepositoryDB.GetUserByEmail: user with this email does not exist", "email", email)
            return nil, errx.ErrEmailDoesntExist
        }
        logger.Error("UserRepositoryDB.GetUserByEmail: failed to retrieve user", "email", email, "error", err)
        return nil, err
    }

    logger.Info("UserRepositoryDB.GetUserByEmail: user retrieved successfully", "email", email)
    return &user, nil
}

func (r *UserRepositoryDB) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
    logger.Info("UserRepositoryDB.GetUserByID: retrieving user by ID", "userID", id)
    query := `
    SELECT user_id, email, password
    FROM users
    WHERE user_id = $1
    `

    var user model.User
    err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            logger.Error("UserRepositoryDB.GetUserByID: user with this ID does not exist", "userID", id)
            return nil, errx.ErrEmailDoesntExist
        }
        logger.Error("UserRepositoryDB.GetUserByID: failed to retrieve user", "userID", id, "error", err)
        return nil, err
    }

    logger.Info("UserRepositoryDB.GetUserByID: user retrieved successfully", "userID", id)
    return &user, nil
}