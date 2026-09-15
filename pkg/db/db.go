package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(255) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);
`

var database *sql.DB

func Init(dbFile string) error {
	var err error

	database, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	_, err = database.Exec(schema)
	if err != nil {
		return err
	}

	return nil
}
