package main

import (
	"testing"
)

func TestLetterCounter(t *testing.T) {
	var letter_count lettercount
	count_letters(word("aaabbczz"), &letter_count)
	if letter_count[0] != 3 {
		t.Errorf("There should be 3 a")
	}
	if letter_count[1] != 2 {
		t.Errorf("There should be 2 b")
	}
	if letter_count[2] != 1 {
		t.Errorf("There should be 1 c")
	}
	if letter_count[25] != 2 {
		t.Errorf("There should be 2 z")
	}

}

func AssertSameLists(t *testing.T, list1 wordlist, list2 wordlist) {
	if len(list1) != len(list2) {
		t.Error("Not the same length ", len(list1), "!=" , len(list2))
		return
	}
	for index, word := range list1 {
		if string(list2[index]) != string(word) {
			t.Errorf("%s not found, %s found instead", list2[index], word)
			return
		}
	}
}

func TestFindWordsInBoard(t *testing.T) {
	board := make_board("supermansqqqqqqqqqqqqqqqq")
	words_to_test := wordlist{word("super"), word("duper"), word("ssuper"), word("impossible"), word("supermans")}
	possible_words := all_possible_moves(board, words_to_test)
	espected_answer := wordlist{word("super"), word("ssuper"), word("supermans")}
	AssertSameLists(t, possible_words, espected_answer)
}

func TestFindSubWords(t *testing.T) {
	words_to_test := wordlist{word("supermans"), word("ssuper"), word("super"), word("duper"), word("mans")}
	words_expected := wordlist{word("supermans"), word("ssuper"), word("duper")}
	AssertSameLists(t, words_expected, filter_out_subwords(words_to_test))
}

func TestFindMovesForWord(t *testing.T) {
	board := make_board("supermansxxxxxxxxxxxxxxxx")
	var movesiter moveiterator
	pos1 := movesiter.Begin(board, word("supermans"))
	if pos1[0] != 0 || pos1[1] != 1 || pos1[8] != 8 {
		t.Errorf("it should not be the position")
	}
	pos2 := movesiter.Next()
	if pos2[0] != 8 || pos1[1] != 1 || pos1[8] != 0 {
		t.Errorf("it should not be the position")
	}
	pos3 := movesiter.Next()
	if pos3 != nil {
		t.Errorf("it should stop")

	}
}

func TestSignatureCalculation(t *testing.T) {
	if calculate_word_signature(word("e")) != 2 {
		t.Errorf("e signature should be 2")
	}

	if calculate_word_signature(word("es")) != 6 {
		t.Errorf("es signature should be 6")
	}

	if calculate_word_signature(word("eeees")) != 48 {
		t.Errorf("eeees signature should be 48")
	}


}

func TestFind_chr_indices(t *testing.T) {
	result := find_chr_indices('s', []byte{ 's', 'u', 'p', 'e', 'r', 'm', 'a', 'n', 's'})
	if result[0] != 0 {
		t.Errorf("should find one at the beginning")
	}
	if result[1] != 8 {
		t.Errorf("should find one at 8")
	}
	if len(result) != 2 {
		t.Errorf("should be only 2 results")
	}
}

func TestNumbermachingletters(t *testing.T) {
	wordsig := calculate_word_signature(word("zorglub"))
	matching := number_maching_letters(wordsig, word("zorh"))
	if matching != 3 {
		t.Errorf("3 letters should have matched")
	}
}
