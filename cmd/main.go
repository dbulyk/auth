package main

import (
	"context"
	"log"
	"time"

	"auth/internal/app"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("Не удалось инициализировать приложение: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("Не удалось запустить приложение: %s", err.Error())
	}
}
