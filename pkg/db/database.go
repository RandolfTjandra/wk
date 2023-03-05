package db

import (
	"database/sql"
	"fmt"
	"os"
)

func InitDatabase() (database *sql.DB, err error) {
	location := fmt.Sprintf("%s", os.ExpandEnv("${HOME}/.local/state/wk"))
	os.Mkdir(location, os.ModePerm)
	database, err = sql.Open("sqlite3", location+"/db.sqlite")
	sts := `
CREATE TABLE IF NOT EXISTS subjects(id INTEGER PRIMARY KEY, subject TEXT);
CREATE TABLE IF NOT EXISTS voice_actors(id INTEGER PRIMARY KEY, voice_actor TEXT);
`
	_, err = database.Exec(sts)
	return
}
