package config

type PostgresDBParams struct {
	Host    string
	Port    int
	DB_name string
	User    string
	Pass    string
}

var PostgresConfig = PostgresDBParams{
	Host:    "localhost",
	Port:    5432,
	DB_name: "kvs",
	User:    "postgres",
	Pass:    "admin",
}
