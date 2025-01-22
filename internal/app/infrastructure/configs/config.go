package configs

import (
	"go.uber.org/fx"
	"os"

	"github.com/joho/godotenv"
)

// Container contains environment variables for the application, database, cache, token, and http server
type (
	Container struct {
		App   *App
		Token *Token
		Redis *Redis
		DB    *DB
		HTTP  *HTTP
	}
	// App contains all the environment variables for the application
	App struct {
		Name string
		Env  string
	}
	// Token contains all the environment variables for the token services
	Token struct {
		Duration string
	}
	// Redis contains all the environment variables for the cache services
	Redis struct {
		Addr     string
		Password string
	}
	// Database contains all the environment variables for the database
	DB struct {
		Connection string
		Host       string
		Port       string
		User       string
		Password   string
		Name       string
	}
	// HTTP contains all the environment variables for the http server
	HTTP struct {
		Env            string
		URL            string
		Port           string
		AllowedOrigins string
	}
)

// NewContainer creates a new container instance
func NewContainer() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	token := &Token{
		Duration: os.Getenv("TOKEN_DURATION"),
	}

	redis := &Redis{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	db := &DB{
		Connection: os.Getenv("DB_CONNECTION"),
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Name:       os.Getenv("DB_NAME"),
	}

	http := &HTTP{
		Env:            os.Getenv("APP_ENV"),
		URL:            os.Getenv("HTTP_URL"),
		Port:           os.Getenv("HTTP_PORT"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
	}

	return &Container{
		app,
		token,
		redis,
		db,
		http,
	}, nil
}

func ProvideToken(container *Container) *Token {
	return container.Token
}

func ProvideDB(container *Container) *DB {
	return container.DB
}

func ProvideRedis(container *Container) *Redis {
	return container.Redis
}

var Module = fx.Module(
	"configs",
	fx.Provide(
		NewContainer,
		ProvideToken,
		ProvideDB,
		ProvideRedis,
	),
)
