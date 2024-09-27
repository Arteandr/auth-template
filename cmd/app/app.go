package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"mzhn/auth/internal/app"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("cannot load env: %w", err))
	}
}

func main() {
	app, _, err := app.New()
	if err != nil {
		panic(err)
	}

	app.Run()
}
