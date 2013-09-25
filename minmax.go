package main

import (
	"fmt"
	"time"
)

const INFINITY = 1000000
const MAX_WORD_CUTOFF = INFINITY

var depth0_alpha int

type variation struct {
	eval    int
	move    []int
	word    *signedword
	nbMoves int
}

func (game *Game) max_value(state *GameState, alpha int, beta int, depth int, max_depth int, modulo int) variation {
	best := variation{eval: -INFINITY}
	if depth >= max_depth {
		best.eval = state.evaluate()
		return best
	}

	var last_nb_moves int
	timestampFloat := time.Now().UnixNano()
	nb_words := 0
	wi := NewWordIterator(game, state, modulo)

	for current_signedword := wi.Next(); current_signedword != nil; current_signedword = wi.Next() {
		nb_words++
		if depth > 1 && nb_words > 5 {
			return best
		}
		var current_move []int
		var mi moveiterator
		for current_move = mi.Begin(game.board, current_signedword.word); current_move != nil; current_move = mi.Next() {
			var eval, new_nb_moves int
			new_state := *state
			new_state.play(current_move, current_signedword, BLUE)

			if new_state.is_finished() {
				eval = new_state.evaluate()
				new_nb_moves = 1
			} else {
				min_variation := game.min_value(&new_state, alpha, beta, depth+1, max_depth)
				eval = min_variation.eval
				new_nb_moves = min_variation.nbMoves
			}
			best.nbMoves += new_nb_moves + 1
			if best.eval < eval {
				best.eval = eval
				best.move = make([]int, len(current_move))
				for index, chr := range current_move {
					best.move[index] = chr
				}
				best.word = current_signedword
			}
			if alpha < best.eval {
				alpha = best.eval
			}
			if depth == 0 {
				if alpha < depth0_alpha {
					alpha = depth0_alpha
				}
				if depth0_alpha < alpha {
					depth0_alpha = alpha
				}
			}
			if alpha >= beta { // alpha/beta cutoff
				return best
			}
		}
		if depth == 0 {
			new_ts := time.Now().UnixNano()
			word_per_seconds := (int64(best.nbMoves-last_nb_moves) * 1000000000) / (new_ts - timestampFloat)
			fmt.Println(modulo, "Best Eval", best.eval, "Best Word", best.word.word, "move", best.move, "     |  Current word", string(current_signedword.word), best.nbMoves, "Speed ", word_per_seconds)
			last_nb_moves = best.nbMoves
		}

	}
	return best
}

func (game *Game) min_value(state *GameState, alpha int, beta int, depth int, max_depth int) variation {
	best := variation{eval: INFINITY}
	if depth >= max_depth {
		best.eval = state.evaluate()
		return best
	}

	var nb_words int
	wi := NewWordIterator(game, state, -1)

	for current_signedword := wi.Next(); current_signedword != nil; current_signedword = wi.Next() {
		nb_words++
		if depth > 1 && nb_words > 1 /* len(game.possible_words)/depth*/ {
			return best
		}

		var current_move []int
		var mi moveiterator
		for current_move = mi.Begin(game.board, current_signedword.word); current_move != nil; current_move = mi.Next() {
			var eval, new_nb_moves int
			new_state := *state
			new_state.play(current_move, current_signedword, RED)
			if new_state.is_finished() {
				eval = new_state.evaluate()
				new_nb_moves = 1
			} else {
				max_variation := game.max_value(&new_state, alpha, beta, depth+1, max_depth, -1)
				eval = max_variation.eval
				new_nb_moves = max_variation.nbMoves
			}
			best.nbMoves += new_nb_moves + 1
			if best.eval > eval {
				best.eval = eval
				best.move = make([]int, len(current_move))
				for index, chr := range current_move {
					best.move[index] = chr
				}
				best.word = current_signedword
			}
			if beta > best.eval {
				beta = best.eval
			}
			if beta <= alpha { // alpha/beta cutoff
				return best
			}

		}

	}
	return best
}

func (game *Game) partition(result chan variation, max_depth int, moduloOffset int) {
	localGS := game.state
	result <- game.max_value(&localGS, -INFINITY, INFINITY, 0, max_depth, moduloOffset)
}

func (game *Game) search(max_depth int) variation {
	var best variation
	depth0_alpha = -INFINITY
	best.eval = -INFINITY
	fmt.Printf("Start \n%v\n", game)
	game.sort_possible_words_by_letter_subsets(game.unused_letterset(), game.interesting_letterset(BLUE))

	result := make(chan variation, MODULO)
	for i := 0; i < MODULO; i++ {
		go game.partition(result, max_depth, i)
	}

	for i := 0; i < MODULO; i++ {
		v := <-result
		fmt.Println("-- PARTIAL RESULT --\nBest word:", string(v.word.word), "Eval:", v.eval)
		if v.eval > best.eval {
			best = v
		}
	}

	fmt.Println("-- RESULT --\nBest word:", string(best.word.word), "Eval:", best.eval)
	return best
}
