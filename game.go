package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

func ReadWords() (wordlist, error) {

	f, err := os.Open("word_cache.txt")
	if err != nil {
		return nil, err
	}
	result := make(wordlist, 0, 271377)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		w := scanner.Bytes()
		if len(w) != 0 {
			c := make(word, len(w))
			copy(c, w)
			result = append(result, c)
		}
	}
	return result, scanner.Err()
}

type signedword struct {
	word      word
	signature wordsig
}

func (s signedword) String() string {
	return fmt.Sprintf("%q(%d)", s.word, s.signature)
}

type Game struct {
	board          *board
	state          GameState
	possible_words []signedword
}

func (g *Game) String() string {
	result := "\x1b[0m\x1b[30m"
	for i, color := range g.state.mask {
		if i%5 == 0 {
			result += "\n"
		}
		intense := color != EMPTY && g.state.mask.vicinity_same_color(color, i)
		switch color {
		case RED:
			if intense {
				result += "\x1b[48;5;196m"
			} else {
				result += "\x1b[48;5;203m"
			}
		case BLUE:
			if intense {
				result += "\x1b[48;5;21m"
			} else {
				result += "\x1b[48;5;62m"
			}
		default:
			result += "\x1b[48;5;248m"
		}
		result += string(g.board[i])
	}
	result += "\x1b[0m\n"
	result += fmt.Sprintf("%v", g.state.played_moves)
	return result
}

type GameState struct {
	mask         mask
	played_moves []signedword
}

type BestmatchSignedWords struct {
	possible_words  []signedword
	all_in_criteria word // subset which is best if they all match (end of game)
	criteria        word // the next best subset of lettes to best match on
}

func (b BestmatchSignedWords) Len() int {
	return len(b.possible_words)
}

func (b BestmatchSignedWords) Swap(i, j int) {
	b.possible_words[i], b.possible_words[j] = b.possible_words[j], b.possible_words[i]
}

func (b BestmatchSignedWords) Less(i, j int) bool {
	left := b.possible_words[i]
	right := b.possible_words[j]

	left_all_in := number_maching_letters(left.signature, b.all_in_criteria) == len(b.all_in_criteria)
	right_all_in := number_maching_letters(right.signature, b.all_in_criteria) == len(b.all_in_criteria)

	if left_all_in && !right_all_in {
		return true
	}
	if !left_all_in && right_all_in {
		return false
	}
	// they are equal so discriminate on the rest
	ii := number_maching_letters(left.signature, b.criteria)
	jj := number_maching_letters(right.signature, b.criteria)
	if ii > jj {
		return true
	} else if ii < jj {
		return false
	}
	return len(left.word) > len(right.word) // if it is equal, sort by length

}

func readGame(r io.Reader) (*Game, error) {
	s := bufio.NewScanner(r)
	boardStr := ""
	maskStr := ""
	played := make([]string, 0)
	for i := 0; i < 5; i++ {
		s.Scan()
		boardStr += s.Text()
	}
	if len(boardStr) != 25 {
		return nil, fmt.Errorf("Could not parse the game %q", boardStr)
	}
	s.Scan() // empty line
	for i := 0; i < 5; i++ {
		s.Scan()
		maskStr += s.Text()
	}
	if len(maskStr) != 25 {
		return nil, fmt.Errorf("Could not parse the mask: %q", maskStr)
	}
	s.Scan() // empty line
	for s.Scan() {
		played = append(played, s.Text())
	}
	return Make_empty_game(boardStr, maskStr, played)
}

func Make_empty_game(board_str string, mask_str string, played_moves_str []string) (*Game, error) {
	game := &Game{board: make_board(board_str)}
	played_moves := make([]word, len(played_moves_str))
	for i := range played_moves_str {
		played_moves[i] = []byte(played_moves_str[i])
	}

	if mask_str == "" {
		game.state.mask.Zap()
	} else {
		game.state.mask = make_mask(mask_str)
	}

	possible_moves, err := ReadWords()
	if err != nil {
		return nil, err
	}
	possible_signed_words := make([]signedword, 0, len(possible_moves))
	for _, word := range filter_out_subwords(all_possible_moves(game.board, possible_moves)) {
		possible_signed_words = append(possible_signed_words, signedword{word, calculate_word_signature(word)})
	}
	game.possible_words = possible_signed_words
	return game, err
}

func (g *Game) sort_possible_words_by_letter_subsets(allIn, criteria word) {
	var bsw BestmatchSignedWords
	bsw.possible_words = g.possible_words
	bsw.all_in_criteria = allIn.uniqueLetters() // TODO improve the algo to do that correctly
	bsw.criteria = criteria.uniqueLetters()
	sort.Sort(bsw)
}

func (sw *signedword) Equal(osw *signedword) bool {
	if sw.signature != osw.signature {
		return false
	}
	return sw.word.Equal(osw.word)
}

func (w word) uniqueLetters() word {
	var result word
	seen := map[byte]byte{}
	for _, val := range w {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}

func (game *Game) unused_letterset() word {
	var word word
	for index, color := range game.state.mask {
		if color == EMPTY {
			word = append(word, game.board[index])
		}
	}
	return word
}

func (game *Game) interesting_letterset(forColor byte) word {
	var word word
	for index, color := range game.state.mask {
		if color == EMPTY || (color != forColor && !game.state.mask.vicinity_same_color(game.state.mask[index], index)) {
			word = append(word, game.board[index])
		}
	}
	fmt.Println("Interesting letterset:", word)
	return word
}

type worditerator struct {
	game          *Game
	gamestate     *GameState
	current_index int
}

func (wi *worditerator) Begin(game *Game, gamestate *GameState) *signedword {
	wi.game = game
	wi.gamestate = gamestate
	return wi.Next()
}

func (wi *worditerator) Next() *signedword {
already_played:
	for _, psw := range wi.game.possible_words[wi.current_index:] {
		wi.current_index += 1
		for _, signed_played_word := range wi.gamestate.played_moves {
			if psw.Equal(&signed_played_word) {
				continue already_played
			}
		}
		return &psw
	}
	return nil
}

func (state *GameState) is_finished() bool {
	for index := range state.mask {
		if state.mask[index] == EMPTY {
			return false
		}
	}
	return true
}

func (state *GameState) evaluate() int {
	var points int
	var total int
	one_empty := false
	for index := range state.mask {
		color := state.mask[index]
		if color == EMPTY {
			one_empty = true
			continue
		}
		if color == BLUE {
			points = 1
		} else {
			points = -1
		}
		if state.mask.vicinity_same_color(color, index) {
			points *= 2
		}
		total += points
	}
	if !one_empty {
		total *= 100 // make that urgent and corrolate the total score
	}
	return total
}

func (state *GameState) play(move []int, signed_word *signedword, color byte) {
	var othercolor byte
	if color == BLUE {
		othercolor = RED
	} else {
		othercolor = BLUE
	}
	old_mask := state.mask
	for _, index := range move {
		if !old_mask.vicinity_same_color(othercolor, index) { // protect it if it is locked by the opponent
			state.mask[index] = color
		}
	}
	state.played_moves = append(state.played_moves, *signed_word)

}
