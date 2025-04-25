package postgres

import (
	"context"

	"studentgit.kata.academy/Zhodaran/go-kata/internal/entities"
)

func (r *PostgresUserRepository) Create(ctx context.Context, user entities.User) error {
	query := "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)"
	_, err := r.Db.ExecContext(ctx, query, user.ID, user.Name, user.Email)
	return err
}

// GetByID получает пользователя по ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (entities.User, error) {
	var user entities.User
	query := "SELECT id, name, email, deleted_at FROM users WHERE id = $1"
	err := r.Db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.DeletedAt)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

// Update обновляет данные пользователя
func (r *PostgresUserRepository) Update(ctx context.Context, user entities.User) error {
	query := "UPDATE users SET name = $1, email = $2 WHERE id = $3"
	_, err := r.Db.ExecContext(ctx, query, user.Name, user.Email, user.ID)
	return err
}

// Delete помечает пользователя как удаленного
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := "UPDATE users SET deleted_at = NOW() WHERE id = $1"
	_, err := r.Db.ExecContext(ctx, query, id)
	return err
}

// List возвращает список пользователей с пагинацией
func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]entities.User, error) {
	query := "SELECT id, name, email, deleted_at FROM users WHERE deleted_at IS NULL LIMIT $1 OFFSET $2"
	rows, err := r.Db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.DeletedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
