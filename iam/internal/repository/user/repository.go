package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	errs "github.com/Anton119/rocket-service-/iam/internal/errors"
	"github.com/Anton119/rocket-service-/iam/internal/model"
	repoconv "github.com/Anton119/rocket-service-/iam/internal/repository/converter"
	"github.com/Anton119/rocket-service-/iam/internal/repository/record"
)

const userSelectColumns = `
	uuid, login, password_hash, created_at, updated_at`

// Repository — PostgreSQL-хранилище пользователей.
type Repository struct {
	pool *pgxpool.Pool
}

// New создаёт репозиторий пользователей.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create сохраняет нового пользователя.
func (r *Repository) Create(ctx context.Context, user model.User) error {
	rec := repoconv.UserToRecord(user)

	const query = `
		INSERT INTO users (uuid, login, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx, query,
		rec.UUID,
		rec.Login,
		rec.PasswordHash,
		rec.CreatedAt,
		rec.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errs.ErrUserAlreadyExists
		}

		return fmt.Errorf("создать пользователя: %w", err)
	}

	return nil
}

// GetByUUID возвращает пользователя по UUID или ErrUserNotFound.
func (r *Repository) GetByUUID(ctx context.Context, id uuid.UUID) (model.User, error) {
	query := `SELECT` + userSelectColumns + `
		FROM users
		WHERE uuid = $1`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return model.User{}, fmt.Errorf("получить пользователя: %w", err)
	}
	defer rows.Close()

	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[record.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errs.ErrUserNotFound
		}

		return model.User{}, fmt.Errorf("получить пользователя: %w", err)
	}

	return repoconv.UserFromRecord(rec), nil
}

// GetByLogin возвращает пользователя по логину или ErrUserNotFound.
func (r *Repository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	query := `SELECT` + userSelectColumns + `
		FROM users
		WHERE login = $1`

	rows, err := r.pool.Query(ctx, query, login)
	if err != nil {
		return model.User{}, fmt.Errorf("получить пользователя по логину: %w", err)
	}
	defer rows.Close()

	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[record.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errs.ErrUserNotFound
		}

		return model.User{}, fmt.Errorf("получить пользователя по логину: %w", err)
	}

	return repoconv.UserFromRecord(rec), nil
}
