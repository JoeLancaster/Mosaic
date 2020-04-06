package stats

import (
	"testing"
)

func TestOpenImage(t *testing.T) {
	_, err := OpenImage("white64x64.png")
	if err != nil {
		t.Errorf("%s", err)
	}
	_, err = OpenImage("white64x64.jpeg")
	if err != nil {
		t.Errorf("%s", err)
	}
	_, err = OpenImage("white64x64.gif")
	if err != nil {
		t.Errorf("%s", err)
	}
	_, err = OpenImage("util.go")
	if err == nil {
		t.Errorf("util.go which is not an image was erroneously loaded successfully.")
	}
}
