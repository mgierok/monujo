package db

import (
	"fmt"
	"os/exec"

	"github.com/mgierok/monujo/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func MustInitialize() {
	dbinfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Db().Host,
		config.Db().Port,
		config.Db().User,
		config.Db().Password,
		config.Db().Dbname,
	)

	var err error
	db, err = sqlx.Connect("postgres", dbinfo)

	if err != nil {
		panic(err)
	}
}

func Connection() *sqlx.DB {
	return db
}

func Dump(dumptype string, file string) {
	if len(file) == 0 {
		fmt.Println("Output file is not set")
		return
	}

	var cmd *exec.Cmd
	if dumptype == "schema" {
		cmd = exec.Command(
			"/usr/bin/pg_dump",
			"--host",
			config.Db().Host,
			"--port",
			config.Db().Port,
			"--username",
			config.Db().User,
			"--no-password",
			"--format",
			"plain",
			"--schema-only",
			"--no-owner",
			"--no-privileges",
			"--no-tablespaces",
			"--no-unlogged-table-data",
			"--file",
			file,
			config.Db().Dbname,
		)
	} else if dumptype == "data" {
		cmd = exec.Command(
			"/usr/bin/pg_dump",
			"--host",
			config.Db().Host,
			"--port",
			config.Db().Port,
			"--username",
			config.Db().User,
			"--no-password",
			"--format",
			"plain",
			"--data-only",
			"--inserts",
			"--disable-triggers",
			"--no-owner",
			"--no-privileges",
			"--no-tablespaces",
			"--no-unlogged-table-data",
			"--file",
			file,
			config.Db().Dbname,
		)
	} else {
		fmt.Println("Invalid dump type, please specify 'schema' or 'data'")
		return
	}

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))
}
