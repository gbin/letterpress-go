package main

import (
	"fmt"
	"strings"
	"testing"
)

var empty []string = make([]string, 0)

func TestGameCreation(t *testing.T) {
	game, _ := Make_empty_game("abcdeabcdeabcdeabcdeabcde", "", empty)
	if game.state.mask[12] != EMPTY {
		fmt.Println(game.state.mask)
		t.Errorf("mask should be empty")
	}

	if game.board[12] != 'c' {
		t.Errorf("game not initialized correctly")
	}
}

func TestBestSubset(t *testing.T) {
	game, _ := Make_empty_game("abcdeabcdeabcdeabcdeabcjy", "", empty)
	found := false
	for _, sword := range game.possible_words {
		if string(sword.word) == "deejayed" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Should have found deejayed in the list")
	}
	game.sort_possible_words_by_letter_subsets(word(""), word("cjy"))
	if string(game.possible_words[0].word) != "jaycee" || string(game.possible_words[1].word) != "deejayed" {
		t.Errorf("game not initialized correctly")
	}
}

func TestInterestingLetterset(t *testing.T) {
	game, _ := Make_empty_game("abcxxdexxxfxxxxxxxxxxxxxx", "rrrr rrr  rr   r         ", empty)
	r := game.interesting_letterset(BLUE)
	if r[0] != 'x' {
		t.Errorf("something is wrong")
	}
}

func TestEvaluation(t *testing.T) {
	var state GameState
	state.mask = make_mask("rrrr " +
		"rrr  " +
		"rr   " +
		"r    " +
		"     ")
	if state.evaluate() != -16 {
		t.Errorf("something is wrong in eval")
	}
}

func TestWordGen(t *testing.T) {
	game, _ := Make_empty_game("abcdeabcdeabcdeabcdeabcjy", "", empty)
	var state GameState
	state.played_moves = make([]signedword, 0)

	sw1, sw2, sw3 := game.possible_words[0], game.possible_words[1], game.possible_words[2]

	state.played_moves = append(state.played_moves, sw2)
	var wi worditerator
	first := wi.Begin(game, &state)
	if !first.Equal(&sw1) {
		t.Errorf("Should have been ", sw1.word)
	}
	state.played_moves = append(state.played_moves, sw1)
	first = wi.Begin(game, &state)
	if !first.Equal(&sw3) {
		t.Errorf("Should have been ", sw3.word)
	}

}

const g = //
`asdqw
zxcas
xcvsd
house
dogfg

rrbb.
.....
r.. r
    .
rrrrr

avocados
fougasses
house
`

func inPlayed(w string, sw []signedword) bool {
	for _, s := range sw {
		if string(s.word) == w {
			return true
		}
	}
	return false
}

func TestReadFromFile(t *testing.T) {
	sr := strings.NewReader(g)
	game, err := readGame(sr)
	if err != nil {
		t.Fatal(err)
	}
	if game.board[0] != 'a' || game.board[24] != 'g' {
		t.Errorf("board not read correctly %v", game.board)
	}
	if inPlayed("house", game.state.played_moves) {
		t.Errorf("house should NOT be in played %v as it is a subword anyway", game.state.played_moves)
	}
	if !inPlayed("avocados", game.state.played_moves) {
		t.Errorf("avocados should be in played %v", game.state.played_moves)
	}
	if !inPlayed("fougasses", game.state.played_moves) {
		t.Errorf("fougasses should be in played %v", game.state.played_moves)
	}

}
