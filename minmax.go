package main

import "fmt"

const INFINITY = 1000000
const MAX_WORD_CUTOFF = INFINITY

var depth0_alpha int

type variation struct {
	eval    int
	move    []move
	word    []*signedword
	nbMoves int
}

func (game *Game) max_value(variation *variation, state *GameState, alpha int, beta int, depth int, max_depth int) int {
	if depth >= max_depth {
		return state.evaluate()
	}
	if state.is_finished() {
		for i := depth; i < max_depth; i++ {
			variation.move[i] = nil
			variation.word[i] = nil
		}
		return state.evaluate()
	}

	max_eval := -INFINITY
	wi := NewWordIterator(game, state)

	for current_signedword := wi.Next(); current_signedword != nil; current_signedword = wi.Next() {
		current_move := game.board.first(current_signedword.word)
		for {
			new_state := *state
			new_state.play(current_move, current_signedword, BLUE)

			eval := game.min_value(variation, &new_state, alpha, beta, depth+1, max_depth)
			if max_eval < eval {
				max_eval = eval
				if alpha < max_eval {
					alpha = max_eval
					if alpha >= beta { // alpha/beta cutoff
						return max_eval
					}
				}
				bestmove := make([]int, len(current_move))
				copy(bestmove, current_move)
				variation.move[depth] = bestmove
				variation.word[depth] = current_signedword
				if depth == 0 {
					fmt.Printf("Eval=%d Variation=%v moves=%v\n", max_eval, variation.word, variation.move)
				}
			}
			if !game.board.next(current_move) {
				break
			}
		}
	}
	return max_eval
}

func (game *Game) min_value(variation *variation, state *GameState, alpha int, beta int, depth int, max_depth int) int {
	if depth >= max_depth {
		return state.evaluate()
	}
	if state.is_finished() {
		for i := depth; i < max_depth; i++ {
			variation.move[i] = nil
			variation.word[i] = nil
		}
		return state.evaluate()
	}
	min_eval := INFINITY
	wi := NewWordIterator(game, state)

	for current_signedword := wi.Next(); current_signedword != nil; current_signedword = wi.Next() {
		current_move := game.board.first(current_signedword.word)
		for {
			new_state := *state
			new_state.play(current_move, current_signedword, RED)
			eval := game.max_value(variation, &new_state, alpha, beta, depth+1, max_depth)
			if min_eval > eval {
				min_eval = eval
				if beta > min_eval {
					beta = min_eval
					if beta <= alpha { // alpha/beta cutoff
						return min_eval
					}
				}
				bestmove := make([]int, len(current_move))
				copy(bestmove, current_move)
				variation.move[depth] = bestmove
				variation.word[depth] = current_signedword
			}
			if !game.board.next(current_move) {
				break
			}
		}

	}
	return min_eval
}

func (game *Game) search(max_depth int) variation {
	best := variation{move: make([]move, max_depth), word: make([]*signedword, max_depth)}
	fmt.Printf("Start \n%v\n", game)
	game.sort_possible_words_by_letter_subsets(game.unused_letterset(), game.interesting_letterset(BLUE))
	best.eval = game.max_value(&best, &game.state, -INFINITY, INFINITY, 0, max_depth)
	return best
}
