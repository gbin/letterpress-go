package main

import (
	"fmt"
    "flag"
	"runtime/pprof"
	"log"
	"os"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	game := Make_empty_game(
		    "cszim" +
			"zrcur" +
			"aaaeb" +
			"pbsai" +
			"isuer",
			"     " +
			" r  r" +
			"    r" +
			"  rrr" +
			"   rr")


	best_evaluation, best_move, best_word, nb_moves := game.search(2) // at depth 1 and empty game, it should simply be the first proposal
	fmt.Println(game.possible_words[:50])
	fmt.Println(game.String())
	fmt.Println(best_move)
	game.state.play(best_move, best_word, BLUE)
	fmt.Println(game.String())
	fmt.Println("Eval: ", best_evaluation)
	fmt.Println("Total moves analyzed: ", nb_moves)

}
