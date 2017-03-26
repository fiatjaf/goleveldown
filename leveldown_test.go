package leveldown

import (
	"os"
	"testing"

	"github.com/fiatjaf/go-levelup"
)

func TestAll(t *testing.T) {
	db := NewDatabase("/tmp/leveldowntest")
	defer os.Remove("/tmp/leveldowntest")

	levelup.BasicTests(db, t)
}
