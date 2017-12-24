package db

import (
	"fmt"
	"os/exec"

	"github.com/mgierok/monujo/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func MustInitialize(c config.Db) {
	dbinfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Dbname,
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

func Dump(d config.Db, s config.Sys, dumptype string, file string) {
	if len(file) == 0 {
		fmt.Println("Output file is not set")
		return
	}

	var cmd *exec.Cmd
	if dumptype == "schema" {
		cmd = exec.Command(
			s.Pgdump,
			"--host",
			d.Host,
			"--port",
			d.Port,
			"--username",
			d.User,
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
			d.Dbname,
		)
	} else if dumptype == "data" {
		cmd = exec.Command(
			s.Pgdump,
			"--host",
			d.Host,
			"--port",
			d.Port,
			"--username",
			d.User,
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
			d.Dbname,
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
