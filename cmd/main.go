package main

import (
	"github.com/bagashiz/go_hexagonal/internal/app"
)

func main() {
	application := app.NewApp()
	application.Run()
}
