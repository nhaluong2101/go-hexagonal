package repositories

import (
	"context"
	"github.com/bagashiz/go_hexagonal/internal/app/adapters/storages/db/postgres"
	"github.com/bagashiz/go_hexagonal/internal/app/core/models"
	"github.com/bagashiz/go_hexagonal/internal/app/core/ports"
	"github.com/bagashiz/go_hexagonal/internal/app/core/utils"
	"go.uber.org/fx"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

/**
 * UserRepository implements ports.UserRepository interface
 * and provides an access to the postgres database
 */
type UserRepository struct {
	db *postgres.DB
}

// NewUserRepository creates a new user repositories instance
func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

// CreateUser creates a new user in the database
func (ur *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := ur.db.QueryBuilder.Insert("users").
		Columns("name", "email", "password").
		Values(user.Name, user.Email, user.Password).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errCode := ur.db.ErrorCode(err); errCode == "23505" {
			return nil, models.ErrConflictingData
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID gets a user by ID from the database
func (ur *UserRepository) GetUserByID(ctx context.Context, id uint64) (*models.User, error) {
	var user models.User

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmailAndPassword gets a user by email from the database
func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

// ListUsers lists all users from the database
func (ur *UserRepository) ListUsers(ctx context.Context, skip, limit uint64) ([]models.User, error) {
	var user models.User
	var users []models.User

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ur.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// UpdateUser updates a user by ID in the database
func (ur *UserRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	name := utils.NullString(user.Name)
	email := utils.NullString(user.Email)
	password := utils.NullString(user.Password)
	role := utils.NullString(string(user.Role))

	query := ur.db.QueryBuilder.Update("users").
		Set("name", sq.Expr("COALESCE(?, name)", name)).
		Set("email", sq.Expr("COALESCE(?, email)", email)).
		Set("password", sq.Expr("COALESCE(?, password)", password)).
		Set("role", sq.Expr("COALESCE(?, role)", role)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": user.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errCode := ur.db.ErrorCode(err); errCode == "23505" {
			return nil, models.ErrConflictingData
		}
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by ID from the database
func (ur *UserRepository) DeleteUser(ctx context.Context, id uint64) error {
	query := ur.db.QueryBuilder.Delete("users").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = ur.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

var UserRepositoryModule = fx.Module(
	"users-repositories-module",
	fx.Provide(
		fx.Annotate(NewUserRepository, fx.As(new(ports.UserRepository))),
	),
)
