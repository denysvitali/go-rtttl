package rtttl_test

import (
	"fmt"
	rtttl "github.com/denysvitali/go-rtttl/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T){
	hauntHouse := "HauntHouse: d=4,o=5,b=108: 2a4, 2e, 2d#, 2b4, 2a4, 2c, 2d, 2a#4, 2e., e, 1f4, 1a4, 1d#, 2e., d, 2c., b4, 1a4, 1p, 2a4, 2e, 2d#, 2b4, 2a4, 2c, 2d, 2a#4, 2e., e, 1f4, 1a4, 1d#, 2e., d, 2c., b4, 1a4"
	ringtone, err := rtttl.Parse(hauntHouse)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, ringtone)
	
	assert.Equal(t, "HauntHouse", ringtone.Name)
	assert.Equal(t, 4, ringtone.DefaultDuration)
	assert.Equal(t, 5, ringtone.DefaultOctave)
	assert.Equal(t, 108, ringtone.BPM)

	fmt.Printf("Notes: %v\n", ringtone.Notes)
}

func TestSimpsons(t *testing.T){
	simpsons := "Simpsons:d=4,o=5,b=160:4c6, 4e6, 4f#6, 8a6, 4.g6, 4e6, 4c6, 8a, 8f#, 8f#, 8f#, 2g, 8p, 8p, 8f#, 8f#, 8f#, 8g, 4a#, 8c6, 8c6, 8c6, 4c6"
	ringtone, err := rtttl.Parse(simpsons)
	if err != nil {
		t.Fatal(err)
	}

	ringtone.Play()
}

func TestKnightRider(t *testing.T){
	knightRider := "KnightRi:d=4,o=5,b=63:16e6, 32f6, 32e6, 8b6, 16e7, 32f7, 32e7, 8b6, 16e6, 32f6, 32e6, 16b6, 16e7, 4d7, 8p, 4p, 16e6, 32f6, 32e6, 8b6, 16e7, 32f7, 32e7, 8b6, 16e6, 32f6, 32e6, 16b6, 16e7, 4f7, 4p"
	ringtone, err := rtttl.Parse(knightRider)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, ringtone)

	assert.Equal(t, "KnightRi", ringtone.Name)
	assert.Equal(t, 4, ringtone.DefaultDuration)
	assert.Equal(t, 5, ringtone.DefaultOctave)
	assert.Equal(t, 63, ringtone.BPM)
	ringtone.Play()
}