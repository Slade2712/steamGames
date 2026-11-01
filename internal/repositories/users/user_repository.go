package users

import (
	"context"
	"database/sql"
	userModel "steam_games/internal/models/users"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, email string, hash string) (*userModel.User, error) {
	query := `
		INSERT INTO users (email, password_hash, role_id)
		VALUES ($1, $2, (SELECT id FROM roles WHERE name = 'USER'))
		RETURNING id, created_at
	`

	user := &userModel.User{Email: email, PasswordHash: hash, Role: "USER"}
	err := r.db.QueryRowContext(ctx, query, email, hash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*userModel.User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, r.name, u.created_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.email = $1
	`

	user := &userModel.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindByID(ctx context.Context, id int) (*userModel.User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, r.name, u.created_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
	`

	user := &userModel.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]userModel.User, error) {
	query := `
		SELECT u.id, u.email, r.name, u.created_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		ORDER BY u.id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []userModel.User{}
	for rows.Next() {
		var user userModel.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
