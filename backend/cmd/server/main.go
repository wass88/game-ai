package main

import (
	"os"

	"github.com/wass88/gameai/lib"
)

func main() {
	dbname := os.Getenv("MYSQL_DATABASE")
	lib.NewDB(dbname)
}
