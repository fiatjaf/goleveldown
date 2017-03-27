package main

import (
	"os"

	"github.com/fiatjaf/goleveldown"
	"github.com/fiatjaf/levelup"
)

func main() {
	db := goleveldown.NewDatabase("/tmp/leveldownexample")
	defer os.RemoveAll("/tmp/leveldownexample")

	levelup.Example(db)
}
