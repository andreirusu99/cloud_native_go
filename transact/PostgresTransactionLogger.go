package transact

import (
	"cloud_native_go/core"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PostgresTransactionLogger struct {
	events chan<- core.Event
	errors <-chan error
	db     *sql.DB
}

type PostgresDBParams struct {
	Host    string
	Port    int
	DB_name string
	User    string
	Pass    string
}

func NewPostgresTransactionLogger(config PostgresDBParams) (core.TransactionLogger, error) {

	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		config.Host,
		config.Port,
		config.DB_name,
		config.User,
		config.Pass,
	))

	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	fmt.Printf("--> Successfully connected to DB \"%s\" on %s:%d\n", config.DB_name, config.Host, config.Port)

	logger := &PostgresTransactionLogger{db: db}

	exists, err := logger.verifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("could not verify Transactions table exists: %w", err)
	}
	if !exists {
		fmt.Println("--> Transactions table not found, please create it")
	} else {
		fmt.Println("--> Transactions table found")
	}

	return logger, nil
}

func (l *PostgresTransactionLogger) verifyTableExists() (bool, error) {
	query := `select count(*) from public."Transactions"`
	rows, err := l.db.Query(query)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return true, nil
}

func (l *PostgresTransactionLogger) LogPut(key, value string) {
	l.events <- core.Event{Type: core.EventPut, Key: key, Value: value}
}

func (l *PostgresTransactionLogger) LogDelete(key string) {
	l.events <- core.Event{Type: core.EventDelete, Key: key}
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *PostgresTransactionLogger) Run() {
	events := make(chan core.Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		defer l.db.Close()
		query := `insert into public."Transactions"
				(event_type, key, value)
				values ($1, $2, $3)`

		for event := range events {
			_, err := l.db.Exec(query, event.Type, event.Key, event.Value)
			if err != nil {
				errors <- fmt.Errorf("failed to insert row in log: %w", err)
			}
		}
	}()
}

func (l *PostgresTransactionLogger) ReplayEvents() (<-chan core.Event, <-chan error) {
	events := make(chan core.Event)
	errors := make(chan error, 1)

	go func() {
		defer func() {
			close(events)
			close(errors)
		}()

		query := `select tr_index, event_type, key, value
				from public."Transactions"
				order by tr_index`

		rows, err := l.db.Query(query)
		if err != nil {
			errors <- fmt.Errorf("sql query error: %w", err)
			return
		}

		defer rows.Close()

		event := core.Event{}

		for rows.Next() {
			if err = rows.Scan(&event.Index, &event.Type, &event.Key, &event.Value); err != nil {
				errors <- fmt.Errorf("error reading row: %w", err)
				return
			}
			events <- event
		}

		if err = rows.Err(); err != nil {
			errors <- fmt.Errorf("error reading Transactions table: %w", err)
		}
	}()

	return events, errors
}
