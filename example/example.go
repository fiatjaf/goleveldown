package main

import (
	"github.com/fiatjaf/goleveldown"
	examples "github.com/fiatjaf/levelup/examples"
)

func main() {
	db := goleveldown.NewDatabase("/tmp/leveldownexample")
	defer db.Erase()

	examples.Example(db)
}
