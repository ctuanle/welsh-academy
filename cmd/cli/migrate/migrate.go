package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var dns, src string
	var up, down, force bool
	var version int

	flag.StringVar(&dns, "dns", "postgres://welsh:azerty@localhost/welsh", "Postgres DNS")
	flag.StringVar(&src, "src", "file://migrations", "Source to migrations files")
	flag.BoolVar(&up, "up", false, "Run migrate up")
	flag.BoolVar(&down, "down", false, "Run migrate down")
	flag.BoolVar(&force, "force", false, "Run migrate force version")
	flag.IntVar(&version, "version", 0, "Version to force")
	flag.Parse()

	action := "none"

	if up {
		if down || force {
			log.Fatal(errors.New("One action at a time (up | down | force). --up for example"))
		}
		action = "up"
	} else if down {
		if force {
			log.Fatal(errors.New("One action at a time (up | down | force). --up for example"))
		}
		action = "down"
	} else if force {
		if version < 1 {
			log.Fatal(errors.New("Please specify a positive version to force (-version=1 for example)"))
		}
		action = "force"
	}

	m, err := migrate.New(src, dns)
	if err != nil {
		log.Fatal(err)
	}

	switch action {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
	case "force":
		if err := m.Force(version); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Nothing to do. Run this command with --help to see possible actions")
	}
	fmt.Println("Done without error.")
}
