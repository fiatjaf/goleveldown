package goleveldown

import (
	"os"
	"testing"

	tests "github.com/fiatjaf/levelup/tests"
)

func TestAll(t *testing.T) {
	db := NewDatabase("/tmp/leveldowntest")
	defer os.RemoveAll("/tmp/leveldowntest")

	tests.Test(db, t)
}
