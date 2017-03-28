package goleveldown

import (
	"testing"

	tests "github.com/fiatjaf/levelup/tests"
)

func TestAll(t *testing.T) {
	db := NewDatabase("/tmp/leveldowntest")
	defer db.Erase()

	tests.Test(db, t)
}
