package sound

import "time"

func ThemeHobbits() {
	notes := []note{
		{392.00, 300 * time.Millisecond}, // G4
		{440.00, 200 * time.Millisecond}, // A4
		{493.88, 200 * time.Millisecond}, // B4
		{523.25, 400 * time.Millisecond}, // C5
		{440.00, 300 * time.Millisecond}, // A4
		{392.00, 400 * time.Millisecond}, // G4
	}
	playSequence(notes)
}

func ThemeElves() {
	notes := []note{
		{659.25, 300 * time.Millisecond}, // E5
		{587.33, 300 * time.Millisecond}, // D5
		{523.25, 400 * time.Millisecond}, // C5
		{587.33, 300 * time.Millisecond}, // D5
		{659.25, 500 * time.Millisecond}, // E5
	}
	playSequence(notes)
}

func ThemeMinasTirith() {
	notes := []note{
		{440.00, 300 * time.Millisecond}, // A4
		{523.25, 300 * time.Millisecond}, // C5
		{587.33, 300 * time.Millisecond}, // D5
		{659.25, 400 * time.Millisecond}, // E5
		{523.25, 400 * time.Millisecond}, // C5
	}
	playSequence(notes)
}

func ThemeMountDoom() {
	notes := []note{
		{196.00, 400 * time.Millisecond}, // G3
		{174.61, 400 * time.Millisecond}, // F3
		{155.56, 500 * time.Millisecond}, // D#3
		{130.81, 600 * time.Millisecond}, // C3
	}
	playSequence(notes)
}

func ThemeAragorn() {
	notes := []note{
		{196.00, 400 * time.Millisecond},  // G3
		{220.00, 400 * time.Millisecond},  // A3
		{246.94, 400 * time.Millisecond},  // B3
		{329.63, 600 * time.Millisecond},  // E4
		{246.94, 300 * time.Millisecond},  // B3
		{196.00, 300 * time.Millisecond},  // G3
		{220.00, 600 * time.Millisecond},  // A3
	}

	playSequence(notes)
}
