package rtttl

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/padster/go-sound/output"
	"github.com/padster/go-sound/sounds"
)

type Ringtone struct {
	Name            string
	DefaultDuration int
	DefaultOctave   int
	BPM             int
	Notes           []Note
}

func Parse(rtttl string) (*Ringtone, error) {
	splitString := strings.Split(rtttl, ":")
	if len(splitString) != 3 {
		return nil, fmt.Errorf("invalid rtttl string: expected 3 parts, found %d", len(splitString))
	}

	ringtone := Ringtone{}
	// Part 1: Ringtone Name
	ringtone.Name = splitString[0]

	// Part 2: Settings
	settings := strings.Split(splitString[1], ",")

	for _, v := range settings {
		kv := strings.Split(v, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid settings format: %v", kv)
		}
		key := kv[0]
		value := kv[1]

		var err error
		switch strings.Trim(key, " ") {
		case "d":
			ringtone.DefaultDuration, err = strconv.Atoi(value)
		case "o":
			ringtone.DefaultOctave, err = strconv.Atoi(value)
		case "b":
			ringtone.BPM, err = strconv.Atoi(value)
		}

		if err != nil {
			return nil, err
		}
	}

	// Part 3: Notes
	notes := strings.Split(splitString[2], ",")
	for _, n := range notes {
		cleanNote := strings.Trim(n, " ")
		note, err := ringtone.parseNote(cleanNote)
		if err != nil {
			return nil, fmt.Errorf("invalid note: %v", err)
		}

		ringtone.Notes = append(ringtone.Notes, *note)
	}

	return &ringtone, nil
}

type Note struct {
	Duration int
	N        string
	Octave   int
	Dotted   bool
}

func (ringtone *Ringtone) parseNote(noteStr string) (*Note, error) {
	noteStr = strings.ToUpper(noteStr)
	noteReader := bytes.NewReader([]byte(noteStr))
	r, _, err := noteReader.ReadRune()
	if err != nil {
		return nil, err
	}

	var durationRunes []rune
	if unicode.IsDigit(r) {
		// Duration
		durationRunes = append(durationRunes, r)
		r, _, err = noteReader.ReadRune()
		if err != nil {
			return nil, err
		}
		if unicode.IsDigit(r) {
			// e.g: 16, 32
			durationRunes = append(durationRunes, r)
		} else {
			_ = noteReader.UnreadRune()
		}
	} else {
		_ = noteReader.UnreadRune()
	}


	note := Note{
		Dotted: false,
	}

	parseDot(noteReader, &note)

	if len(durationRunes) > 0 {
		note.Duration, err = strconv.Atoi(string(durationRunes))
		if err != nil {
			return nil, err
		}
	}
	
	r, _, err = noteReader.ReadRune()
	if err != nil {
		return nil, err
	}

	// Note
	var noteValue []rune
	if !unicode.IsLetter(r) {
		return nil, fmt.Errorf("invalid format: expected pitch but found %v", string(r))
	}

	noteValue = append(noteValue, r)
	r, _, err = noteReader.ReadRune()
	if err == io.EOF {
		// Octave specifier missing, note ended
		note.N = string(noteValue)
		note.Octave = ringtone.DefaultOctave
		return &note, nil
	}

	if r == '#' {
		// e.g: A#
		noteValue = append(noteValue, r)
		note.N = string(noteValue)
	} else {
		note.N = string(noteValue)
		_ = noteReader.UnreadRune()
	}

	parseDot(noteReader, &note)

	r, _, err = noteReader.ReadRune()
	if err == io.EOF {
		// Octave specifier missing
		note.Octave = ringtone.DefaultOctave
		return &note, nil
	}

	if !unicode.IsDigit(r) {
		return nil, fmt.Errorf("invalid octave specifier: %v", string(r))
	}

	note.Octave, err = strconv.Atoi(string([]rune{r}))
	if err != nil {
		return nil, err
	}

	parseDot(noteReader, &note)

	return &note, nil
}

func parseDot(reader *bytes.Reader, n *Note) {
	r, _, err := reader.ReadRune()
	if err == io.EOF {
		_ = reader.UnreadRune()
		return
	}
	
	if r == '.' {
		n.Dotted = true
	} else {
		_ = reader.UnreadRune()
	}
}

var frequencyMap = map[string]float64{
	"A":  440.00,
	"A#": 466.16,
	"B":  493.88,
	"C":  261.63,
	"C#": 277.18,
	"D":  293.66,
	"D#": 311.13,
	"E":  329.63,
	"F":  349.23,
	"F#": 369.99,
	"G":  392.00,
	"G#": 415.30,
	"P": 0,
}

func (r *Ringtone) Play() {
	var soundsElements []sounds.Sound
	for _, n := range r.Notes {
		theDuration := n.Duration
		if n.Dotted {
			// Dotted:
			// "the (...) dot increases the duration of the basic note by half"
			theDuration -= theDuration/2
		}
		noteLengthMs := (60.0 / float64(r.BPM)) / float64(theDuration) * 4
		var soundElement sounds.Sound
		if n.N == "P" {
			soundElement = sounds.NewTimedSilence(noteLengthMs * 1000)
		} else {
			freq := frequencyMap[n.N] + frequencyMap[n.N] * math.Pow(2, float64(n.Octave - 4))
			soundElement = sounds.NewTimedSound(sounds.NewSineWave(freq), noteLengthMs*1000)
		}
		soundsElements = append(soundsElements, soundElement)
	}
	toPlay := sounds.ConcatSounds(soundsElements...)
	output.Play(toPlay)
}
