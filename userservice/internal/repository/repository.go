package repository

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"log"
	"main/internal/entity"
	"os"
)

type Repositorer interface {
	GetUser(email, password string) (entity.User, error)
	ListUsers() ([]entity.User, error)
	CreateUser(user entity.User) error
	Migrate(ctx context.Context) error
}

type Repository struct{}

func NewRepository() Repositorer {
	var repo Repository
	return &repo
}
func (r *Repository) Migrate(ctx context.Context) error {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	db, err := sql.Open("postgres", string("host="+os.
		Getenv("DB_HOST")+" port="+os.
		Getenv("DB_PORT")+" entity="+os.
		Getenv("DB_USER")+" password="+os.
		Getenv("DB_PASSWORD")+" dbname="+os.
		Getenv("DB_NAME")+" sslmode=disable"))
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
	defer db.Close()
	if err := goose.Up(db, "migrations/user"); err != nil {
		log.Fatalf("Error applying migrations: %s", err)
	}
	return err
}
func Connect(ctx context.Context) (*pgx.Conn, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	conn, err := pgx.Connect(ctx, string("host="+os.
		Getenv("DB_HOST")+" port="+os.
		Getenv("DB_PORT")+" entity="+os.
		Getenv("DB_USER")+" password="+os.
		Getenv("DB_PASSWORD")+" dbname="+os.
		Getenv("DB_NAME")+" sslmode=disable"))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
func (r *Repository) CreateUser(user entity.User) error {
	conn, err := Connect(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	return err
}
func (r *Repository) GetUser(email, password string) (entity.User, error) {
	conn, err := Connect(context.Background())
	if err != nil {
		return entity.User{}, err
	}
	defer conn.Close(context.Background())
	var user entity.User
	err = conn.QueryRow(context.Background(), "SELECT username, password FROM users WHERE email = $1, password = $2", email, password).Scan(&user.Username, &user.Password)
	if err != nil {
		return entity.User{}, err
	}
	return entity.User{}, nil
}
func (r *Repository) ListUsers() ([]entity.User, error) {
	conn, err := Connect(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())
	rows, err := conn.Query(context.Background(), "SELECT id, username, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.Username, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
