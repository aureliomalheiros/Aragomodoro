package sound

import "testing"

func TestThemeHobbits(t *testing.T) {
	Mute = true 
	ThemeHobbits()
}

func TestThemeElves(t *testing.T) {
	Mute = true
	ThemeElves()
}

func TestThemeAragorn(t *testing.T) {
	Mute = true
	ThemeAragorn()
}

