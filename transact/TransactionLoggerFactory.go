package transact

import (
	"cloud_native_go/core"
	"fmt"
	"os"
	"strconv"
)

func NewTransactionLogger() (core.TransactionLogger, error) {
	loggerType := os.Getenv("TLOG_TYPE")
	switch loggerType {
	case "file":
		return NewFileTransactionLogger(os.Getenv("TLOG_FILENAME"))

	case "postgres":
		port, err := strconv.Atoi(os.Getenv("TLOG_DB_PORT"))
		if err != nil {
			return nil, fmt.Errorf("could not create PostgresTransactionLogger: %w", err)
		}
		return NewPostgresTransactionLogger(
			PostgresDBParams{
				Host:    os.Getenv("TLOG_DB_HOST"),
				Port:    port,
				DB_name: os.Getenv("TLOG_DB_NAME"),
				User:    os.Getenv("TLOG_DB_USER"),
				Pass:    os.Getenv("TLOG_DB_PASS"),
			},
		)

	case "":
		return nil, fmt.Errorf("transaction logger type not defined")

	default:
		return nil, fmt.Errorf("transaction logger type not correct: %s", loggerType)

	}
}
