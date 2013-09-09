package main

import (
	"fmt"
	"time"
)

const INFINITY = 1000000
const MAX_WORD_CUTOFF = INFINITY

func (game *Game) max_value(state *GameState, alpha int, beta int, depth int, max_depth int) (best_evaluation int, best_move []int, best_word *signedword, nb_moves int) {
	if depth >= max_depth {
		return state.evaluate(), best_move, best_word, nb_moves
	}
	best_evaluation = -INFINITY

	var current_signedword *signedword
	var wi worditerator
	var last_nb_moves int
	timestampFloat := time.Now().UnixNano()
	nb_words := 0
	for current_signedword = wi.Begin(game, state); current_signedword != nil; current_signedword = wi.Next() {
		nb_words++
		if nb_words > MAX_WORD_CUTOFF {
			return
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
				// fmt.Println("FINISHED BLUE", eval, current_move, best_move)
			} else {
				eval, _, _, new_nb_moves = game.min_value(&new_state, alpha, beta, depth+1, max_depth)
			}
			nb_moves += new_nb_moves + 1
			if best_evaluation < eval {
				best_evaluation = eval
				best_move = make([]int, len(current_move))
				for index, chr := range current_move {
					best_move[index] = chr
				}
				best_word = current_signedword
				//fmt.Println("BETTER MOVE", eval, current_move, best_move)
			}
			if best_evaluation >= beta { // alpha/beta cutoff
				return
			}
			if alpha < best_evaluation {
				alpha = best_evaluation
			}
		}
		if depth == 0 {
			new_ts := time.Now().UnixNano()
			word_per_seconds := (int64(nb_moves-last_nb_moves) * 1000000000) / (new_ts - timestampFloat)
			fmt.Println("Best Eval", best_evaluation, "Best Word", best_word.word, "move", best_move, "     |  Current word", string(current_signedword.word), nb_moves, "Speed ", word_per_seconds)
			last_nb_moves = nb_moves
		}

	}
	return
}

func (game *Game) min_value(state *GameState, alpha int, beta int, depth int, max_depth int) (best_evaluation int, best_move []int, best_word *signedword, nb_moves int) {
	if depth >= max_depth {
		return state.evaluate(), best_move, best_word, nb_moves
	}
	best_evaluation = INFINITY

	var current_signedword *signedword
	var wi worditerator
	var nb_words int

	for current_signedword = wi.Begin(game, state); current_signedword != nil; current_signedword = wi.Next() {
		nb_words++
		if nb_words > MAX_WORD_CUTOFF {
			return
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
				//fmt.Println("FINISHED RED", eval)
			} else {
				eval, _, _, new_nb_moves = game.max_value(&new_state, alpha, beta, depth+1, max_depth)
			}
			nb_moves += new_nb_moves + 1
			if best_evaluation > eval {
				best_evaluation = eval
				best_move = make([]int, len(current_move))
				for index, chr := range current_move {
					best_move[index] = chr
				}
				best_word = current_signedword
			}
			if best_evaluation <= alpha { // alpha/beta cutoff
				return best_evaluation, best_move, best_word, nb_moves
			}
			if beta > best_evaluation {
				beta = best_evaluation
			}
		}

	}
	return best_evaluation, best_move, best_word, nb_moves
}

func (game *Game) search(max_depth int) (best_evaluation int, best_move move, best_word *signedword, nb_moves int) {
	fmt.Printf("Start \n%v\n", game)
	game.sort_possible_words_by_letter_subsets(game.unused_letterset(), game.interesting_letterset(BLUE))
	fmt.Println(game.possible_words)
	best_evaluation, best_move, best_word, nb_moves = game.max_value(&game.state, -INFINITY, INFINITY, 0, max_depth)
	fmt.Println("-- RESULT --\nBest word:", string(best_word.word), "Eval:", best_evaluation)
	return best_evaluation, best_move, best_word, nb_moves
}
