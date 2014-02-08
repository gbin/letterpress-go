package main

import (
	"fmt"
	"testing"
)

func TestTrivialMinmax(t *testing.T) {
	game, _ := Make_empty_game("abcdeabcdeabcdeabcdeabcjy", "", empty)
	sw1 := game.possible_words[0]
	best := game.search(1) // at depth 1 and empty game, it should simply be the first proposal
	if !best.word[0].Equal(sw1) {
		t.Errorf("Should have been %q but it is %v", sw1.word, best.word[0])
	}

	if best.eval != 11 {
		t.Errorf("Should have been 11 but it is %d", best.eval)
	}

	game.state.play(best.move[0], best.word[0], BLUE)
	fmt.Println(game.String())

}
