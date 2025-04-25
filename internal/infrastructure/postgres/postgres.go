package postgres

import "database/sql"

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{Db: db}
}

type PostgresBookRepository struct {
	db *sql.DB
}

func NewPostgresBookRepository(db *sql.DB) *PostgresBookRepository {
	return &PostgresBookRepository{db: db}
}

type PostgresUserRepository struct {
	Db *sql.DB
}

func NewPostgresAuthorRepository(db *sql.DB) *PostgresAuthorRepository {
	return &PostgresAuthorRepository{db: db}
}

type PostgresAuthorRepository struct {
	db *sql.DB
}

func NewPostgresAuthRepository(db *sql.DB) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}

type PostgresAuthRepository struct {
	db *sql.DB
}
