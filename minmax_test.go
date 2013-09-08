package main

import (
	"fmt"
	"testing"
)

func TestTrivialMinmax(t *testing.T) {
	game, _ := Make_empty_game("abcdeabcdeabcdeabcdeabcjy", "", empty)
	sw1 := game.possible_words[0]
	best_evaluation, best_move, best_word, _ := game.search(1) // at depth 1 and empty game, it should simply be the first proposal
	if !best_word.Equal(&sw1) {
		t.Errorf("Should have been %q but it is %v", sw1.word, best_word)
	}

	if best_evaluation != 11 {
		t.Errorf("Should have been 11 but it is %d", best_evaluation)
	}

	game.state.play(best_move, best_word, BLUE)
	fmt.Println(game.String())

}
