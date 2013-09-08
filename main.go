package main

import (
	"os"

	"fmt"
	"log"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("letterpress-go <game.txt>")
		os.Exit(-1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Could not open %q", os.Args[1])
		os.Exit(-1)
	}
	defer file.Close()
	game, err := readGame(file)
	if err != nil {
		log.Fatalf("Could not initialize the game %v", err)
	}
	best_evaluation, best_move, best_word, nb_moves := game.search(2) // at depth 1 and empty game, it should simply be the first proposal
	fmt.Println(game.possible_words[:50])
	fmt.Println(game.String())
	fmt.Println(best_move)
	game.state.play(best_move, best_word, BLUE)
	fmt.Println(game.String())
	fmt.Println("Eval: ", best_evaluation)
	fmt.Println("Total moves analyzed: ", nb_moves)

}
