package main

import (
	"testing"
	"fmt"
)

func TestTrivialMinmax(t *testing.T) {
	var game *Game = Make_empty_game("abcdeabcdeabcdeabcdeabcjy", "")
	sw1 := game.possible_words[0]
	best_evaluation, best_move, best_word := game.search(1) // at depth 1 and empty game, it should simply be the first proposal
	if !best_word.Equal(&sw1) {
		t.Errorf("Should have been ", sw1.word)
	}

	if best_evaluation != 11 {
		t.Errorf("Should have been 11 but ", best_evaluation)
	}

	game.state.play(best_move, best_word, BLUE)
	fmt.Println(game.String())

}
