package main

import (
	"bufio"
	"fmt"
	"bytes"
	"sort"
)

const SIZE = 5
const LENGTH = SIZE*SIZE
const RED = 'r'
const BLUE = 'b'
const EMPTY = ' '

type word []byte
type board [LENGTH]byte
type mask  [LENGTH]byte
type wordlist []word
type lettercount [26]int8
type wordsig int64

// this is by probability of finding *at least* this letter in the words
// var letters_order [26]byte =[26]byte{'e', 's', 'i', 'a', 'r', 'n', 't', 'o', 'l', 'c', 'u', 'd', 'p', 'm', 'g', 'h', 'b', 'y', 'f', 'v', 'k', 'w', 'z', 'x', 'q', 'j'}
// var primes [26]int = [26]int{ 2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101 }

// this is the prime associated with the alphabet in order to minimize the product values
var primes_letters_order [26]int64 = [26]int64{7, 59, 29, 37, 2, 67, 47, 53, 5, 101, 73, 23, 43, 13, 19, 41, 97, 11, 3, 17, 31, 71, 79, 89, 61, 83}

// define Interface for a natural sort by size
func (s wordlist) Len() int { return len(s) }
func (s wordlist) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s wordlist) Less(i, j int) bool { return len(s[i]) > len(s[j])}

func (s *mask) Zap() {
	for index := range s {
		s[index] = EMPTY
	}
}

func (mask *mask) vicinity_same_color(color byte, index int) bool {
	x, y := index%SIZE, index/SIZE
	if x == 0 {
		if y == 0 {
			if mask[index + 1] != color || mask[index + SIZE] != color {
				return false  // top left corner border is not dark
			}
		} else if y == SIZE - 1 {
			if mask[index - SIZE] != color || mask[index + 1] != color {
				return false // bottom left corner border is not dark
			}
		} else if mask[index - SIZE] != color || mask[index + SIZE] != color || mask[index + 1] != color {
			return false // left border is not dark
		}
	} else if x == SIZE - 1 {
		if y == 0 {
			if mask[index - 1] != color || mask[index + SIZE] != color {
				return false // top right corner border is not dark
			}
		} else if y == SIZE - 1 {
			if mask[index - 1] != color || mask[index - SIZE] != color {
				return false // bottom right corner border is not dark
			}
		} else if mask[index - SIZE] != color || mask[index + SIZE] != color || mask[index - 1] != color {
			return false // right border is not dark
		}
	}  else if y == 0 {
		if mask[index - 1] != color || mask[index + 1] != color || mask[index + SIZE] != color {
			return false // top border is not dark
		}
	} else if y == SIZE - 1 {
		if mask[index - 1] != color || mask[index + 1] != color || mask[index - SIZE] != color {
			return false // bottom border is dark
		}
	}  else if mask[index - SIZE] != color || mask[index + SIZE] != color || mask[index - 1] != color || mask[index + 1] != color {
		return false // a centered surrounder so dark too
	}
	return true
}

func make_board(s string) *board {
	var b board
	for index, letter := range s {
		b[index] = byte(letter)
	}

	return &b
}

func make_mask(s string) mask {
	var b mask
	for index, letter := range s {
		b[index] = byte(letter)
	}
	return b
}


func Readln(r *bufio.Reader) (word, error) {
	var (isPrefix bool = true
		err error      = nil
		line, ln word
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return ln, err
}


func count_letters(letters []byte, letter_count *lettercount) {
	for i := range letter_count {
		letter_count[i] = 0
	}
	for _, letter := range letters {
		letter_count[letter - 'a']++
	}
}

func count_to_string(counts lettercount) string {
	var response string = ""
	for index, count := range counts {
		response += fmt.Sprintf("%s = %d\n", string('a' + index), count)
	}
	return response
}

func all_possible_moves(board *board, words wordlist) wordlist {
	var result wordlist = make(wordlist, 0)
	var board_letter_count lettercount
	var word_letter_count lettercount
	count_letters(board[:], &board_letter_count)
next_word:
	for _, word := range words {
		count_letters(word, &word_letter_count)
		for index, count := range word_letter_count {
			if count > board_letter_count[index] {
				continue next_word
			}
		}
		result = append(result, word)
	}
	return result
}

func filter_out_subwords(words wordlist) wordlist {
	var result wordlist = make(wordlist,0, len(words))
	var index int
	result = append(result, words[0]) // The biggest one never match a smaller one by definition
nextword:
	for index = len(words) - 1; index > 0; index-- {
		var word = words[index]
		for _, otherword := range words[:index - 1] {
			if bytes.Index(otherword, word) != -1 {
				continue nextword
			}
		}
		result = append(result, word)
	}
	sort.Sort(result)
	return result
}

func clone_masks(origin []mask) []mask {
	var clone []mask = make([]mask,0, len(origin)*2)
	for _, m := range (origin) {
		clone = append(clone, m)
	}
	return clone
}

func is_mask_in_list(to_test *mask, list []mask) bool {
	for _, m := range list {
		if bytes.Equal(to_test[:], m[:]) {
			return true
		}
	}
	return false
}

func find_chr_indices(chr byte, in []byte) []int {
	var result = make([]int,0, len(in))
	for index, otherchr := range in {
		if chr == otherchr {
			result = append(result, index)
		}
	}
	return result
}

type moveiterator struct {
	letter_positions [][]int
	current_state    []int
	current_move     []int
}

func (mi *moveiterator) update() bool {
	current_move := mi.current_move
	for index := range (mi.current_state) {
		new_indice := mi.letter_positions[index][mi.current_state[index]]
		for index2 := 0; index2 < index; index2++ { // check if we don't have it yet
			if current_move[index2] == new_indice {
				return false
			}
		}
		current_move[index] = new_indice
	}
	return true
}

func (mi *moveiterator) Begin(board *board, word word) []int {
	l := len(word)
	mi.letter_positions = make([][]int,l, l)
	mi.current_state = make([]int,l, l)
	mi.current_move = make([]int,l, l)

	b := board[:]
	for index, letter := range word {
		mi.letter_positions[index] = find_chr_indices(letter, b)
	}

	if mi.update() {
		return mi.current_move
	}
	return mi.Next()
}

func (mi *moveiterator) Next() []int {
	right := len(mi.current_state) - 1

	for i := right; i >= 0; i-- {
		if mi.current_state[i] < len(mi.letter_positions[i]) - 1 {
			mi.current_state[i]++
			for a := i + 1; a <= right; a++ {
				mi.current_state[a] = 0
			}

			if mi.update() {
				return mi.current_move
			} else {
				i = right
				continue
			}
		}
	}
	return nil
}

func calculate_word_signature(word word) wordsig {
	var result int64 = 1
	for _, letter := range word {
		result*= primes_letters_order[letter - 'a']
	}
	return wordsig(result)
}

func number_maching_letters(signature wordsig, letters word) int {
	var result int
	sig := int64(signature)
	for _, letter := range letters {
		if (sig%primes_letters_order[letter - 'a']) == 0 {
			result++
		}
	}
	return result

}

func (s word) String() string {
	return string([]byte(s))
}

func (s word) Equal(other word) bool {
	if len(s) != len(other) {
		return false
	}

	for index, letter := range(s) {
		if other[index] != letter {
			return false
		}
	}
	return true

}

