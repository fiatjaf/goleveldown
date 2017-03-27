package goleveldown

import (
	"os"
	"testing"

	"github.com/fiatjaf/levelup"
)

func TestAll(t *testing.T) {
	db := NewDatabase("/tmp/leveldowntest")
	defer os.RemoveAll("/tmp/leveldowntest")

	levelup.BasicTests(db, t)
}
