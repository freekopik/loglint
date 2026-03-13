package a

import (
	"log"
	"log/slog"
)

func main() {
	log.Print("Starting server") // want "log message must start with a lowercase letter"

	slog.Info("ошибка входа") // want "log message must be in English only"

	log.Fatal("done!") // want "log message must not contain special characters"

	slog.Debug("user password: 123") // want "log message contains potentially sensitive data"

	log.Print("server started")
}
