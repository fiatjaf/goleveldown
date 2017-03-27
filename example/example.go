package main

import (
	"os"

	"github.com/fiatjaf/goleveldown"
	examples "github.com/fiatjaf/levelup/examples"
)

func main() {
	db := goleveldown.NewDatabase("/tmp/leveldownexample")
	defer os.RemoveAll("/tmp/leveldownexample")

	examples.Example(db)
}
