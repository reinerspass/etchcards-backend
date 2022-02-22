package waldego_test

import (
	"testing"

	"github.com/reinerspass/waldego/waldego"
)

func TestWaldego(t *testing.T) {
	if waldego.Mascot() != "Go Gopher" {
		t.Fatal("wrong mascot")
	}
}
