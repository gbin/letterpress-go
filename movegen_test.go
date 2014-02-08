package main

import "testing"

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
		t.Error("Not the same length ", len(list1), "!=", len(list2))
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
	pos := board.first(word("supermans"))
	if pos[0] != 0 || pos[1] != 1 || pos[8] != 8 {
		t.Errorf("it should not be the position %v", pos)
	}

	if board.next(pos) {
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

func TestNumbermachingletters(t *testing.T) {
	wordsig := calculate_word_signature(word("zorglub"))
	matching := number_maching_letters(wordsig, word("zorh"))
	if matching != 3 {
		t.Errorf("3 letters should have matched")
	}
}

func same(t *testing.T, got, want move) {
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got %v but wanted %v", got, want)
			break
		}
	}
}

func TestNewMoveIterator(t *testing.T) {
	//                   0123456789012345678901234
	board := make_board("supermansmnsupeumnesruans")

	got := board.first(word("ssupsu"))
	same(t, got, []int{0, 8, 1, 2, 11, 12})

	// s positions [0, 8, 11, 19, 24]
	// u positions [1, 12, 15, 21]
	// p positions [2, 13]

	wants := [][]int{
		[]int{0, 8, 1, 2, 11, 15},
		[]int{0, 8, 1, 2, 11, 21},
		[]int{0, 8, 1, 2, 19, 12},
		[]int{0, 8, 1, 2, 19, 15},
		[]int{0, 8, 1, 2, 19, 21},
		[]int{0, 8, 1, 2, 24, 12},
		[]int{0, 8, 1, 2, 24, 15},
		[]int{0, 8, 1, 2, 24, 21},
		[]int{0, 8, 1, 13, 11, 12},
	}

	for _, want := range wants {
		if !board.next(got) {
			t.Errorf("Should have been true")
		}
		same(t, got, want)
	}
}

func BenchmarkNewMoveIterator(b *testing.B) {
	board := make_board("supermansmnsupeumnesruans")
	for i := 0; i < b.N; i++ {
		nb := 1
		move := board.first(word("supermans"))
		for board.next(move) {
			nb++
		}
		if nb != 11520 {
			b.Errorf("Something changed in the algorithm 11520 != %d", nb)
		}
	}
}
