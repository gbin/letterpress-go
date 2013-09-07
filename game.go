package main

import (
	"os"
	"fmt"
	"bufio"
	"sort"
)

func ReadWords() wordlist {
	var result wordlist = make(wordlist,0, 271377)

	f, err := os.Open("word_cache.txt")
	if err != nil {
		fmt.Println("error opening file= ", err)
		os.Exit(1)
	}
	reader := bufio.NewReader(f)
	line, err := Readln(reader)
	for err == nil {
		line, err = Readln(reader)
		if len(line) != 0 {
			result = append(result, line)
		}
	}
	return result
}

type signedword struct {
	word      word
	signature wordsig
}

func (s signedword) String() string {
	return fmt.Sprintf("%v [%d]", s.word, s.signature)
}

type Game struct {
	board *board
	state          GameState
	possible_words []signedword
}

func (g *Game) String() string {
	var result string
	for i:=0;i<=20; i+=5 {
		result+=string(g.board[i:i+5])
		result+=" "
		result+=string(g.state.mask[i:i+5])
		result+="\n"
	}
	result += fmt.Sprintf("%v", g.state.played_moves)
	return result
}


type GameState struct {
	mask          mask
	played_moves  []signedword
}

type BestmatchSignedWords struct {
	possible_words  []signedword
	subset_criteria word // subset of letters to best match on
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
	ii := number_maching_letters(left.signature, b.subset_criteria)
	jj := number_maching_letters(right.signature, b.subset_criteria)
	if ii > jj {
		return true
	} else if ii < jj {
		return false
	}
	return len(left.word) > len(right.word) // if it is equal, sort by length

}

func Make_empty_game(board_str string, mask_str string) *Game {
    game := Game{board:make_board(board_str)}

	if mask_str == "" {
		game.state.mask.Zap()
	} else {
		game.state.mask = make_mask(mask_str)
	}

    possible_moves := ReadWords()
    possible_signed_words := make([]signedword,0, len(possible_moves))
	for _, word := range all_possible_moves(game.board, possible_moves) {
		var signedword signedword
		signedword.word = word
		signedword.signature = calculate_word_signature(word)
		possible_signed_words = append(possible_signed_words, signedword)
	}
	game.possible_words = possible_signed_words
	return &game
}

func (g *Game) sort_possible_words_by_letter_subset(subset word) {
	var bsw BestmatchSignedWords
	bsw.possible_words = g.possible_words
	bsw.subset_criteria = subset
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

func (game *Game) interesting_letterset() word {
	var word word
	for index := range (game.state.mask) {
		if !game.state.mask.vicinity_same_color(game.state.mask[index], index) {
			word = append(word, game.board[index])
		}
	}
	return word.uniqueLetters()
}

type worditerator struct {
	game *Game
	gamestate *GameState
	current_index int
}

func (wi *worditerator) Begin(game *Game, gamestate *GameState) *signedword {
	wi.game = game
	wi.gamestate = gamestate
	return wi.Next()
}

func (wi *worditerator) Next() *signedword {
already_played:
	for _, psw := range (wi.game.possible_words[wi.current_index:]) {
		wi.current_index += 1
		for _, signed_played_word := range (wi.gamestate.played_moves) {
			if psw.Equal(&signed_played_word) {
				continue already_played
			}
		}
		return &psw
	}
	return nil
}


func (state *GameState) is_finished() bool {
	for index := range (state.mask) {
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
	for index := range (state.mask) {
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
		total *= 100  // make that urgent and corrolate the total score
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
	for _, index := range (move) {
		if !old_mask.vicinity_same_color(othercolor, index) { // protect it if it is locked by the opponent
			state.mask[index] = color
		}
	}
	state.played_moves = append(state.played_moves, *signed_word)

}
