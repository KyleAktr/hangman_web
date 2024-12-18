package game

import (
	"fmt"
	"hang_web/savegame"
	"strings"
)

// GuessingWord : check if the word guessed is the word to find
func GuessingWord(guess string, data savegame.GameData) bool {
	if guess == strings.ToUpper(data.ToFind) {
		return true
	} else {
		fmt.Print("\033[31m")
		fmt.Printf("Sorry, the word is not %v . You lose 2 attemps ...\n", guess)
		fmt.Print("\033[0m")
		return false
	}
}

// WordGuessed : Check if we found the word or not
func WordGuessed(data savegame.GameData) bool {
	count := 0
	for _, v := range data.Word {
		if v == 95 {
			count++
		}
	}
	if count == 0 {
		return true
	}
	return false
}

// InitialLetters : chose 0<n<len(word) / 2 - 1 letters in the word to print and return the index of these letters in th word
func InitialLetters(data savegame.GameData) savegame.GameData {
	var initialLetters []string
	var n int
	if len(data.ToFind) < 4 { // error when len word < 3 ( 3/2-1 == 0 ) We cannot take a random number between 0 & 0
		n = 1

	} else {
		n = RandomNumber((len(data.ToFind) / 2) - 1)
	}
	if n == 0 {
		n = (len(data.ToFind) / 2) - 1
	}
	for i := 0; i < n; i++ {
		index := RandomNumber(len(data.ToFind))
		initialLetters = append(initialLetters, string(data.ToFind[index]))
		data.Tries = initialLetters
	}
	return data
}

// RevealInitialLetters : Take InitialLetters indexes and return the initial word to show
func RevealInitialLetters(data savegame.GameData) savegame.GameData {
	LetterInWord := strings.Split(data.ToFind, "")
	for _, letter := range data.Tries {
		for i, v := range LetterInWord {
			if letter == v {
				data.Index = append(data.Index, i)
			}
		}
	}
	InitialWord := make([]string, len(data.ToFind))
	for i := range data.ToFind {
		InitialWord[i] = "_"
	}
	for _, index := range data.Index {
		InitialWord[index] = string(data.ToFind[index])
	}
	reveal := strings.Join(InitialWord, "")
	data.Word = strings.ToUpper(reveal)
	return data
}

// RevealLetters : return the word with a letter revealed
func RevealLetters(data savegame.GameData) string {
	revealRune := []rune(data.Word)
	wordRune := []rune(data.ToFind)
	for _, v := range data.Index {
		revealRune[v] = wordRune[v]
	}
	data.Word = strings.ToUpper(string(revealRune))
	return data.Word
}

// IntputTesting : Test if the input is valid
func IntputTesting(guess string, data savegame.GameData) bool {
	if guess == "" {
		fmt.Println("\033[31m", "Please do not enter a blank entry.", "\033[0m")
		return false
	}
	if len(guess) != 1 {
		for _, v := range guess {
			if string(v) < "A" || string(v) > "Z" {
				fmt.Println("\033[31m", "Error. Please enter only a letter or a word composed of letters only.", "\033[0m")
				return false
			} else {
				continue
			}
		}
		return true
	}
	if guess >= "A" && guess <= "Z" {
		for _, tries := range data.Tries {
			if strings.ToUpper(tries) == guess {
				if strings.Contains(data.Word, guess) == true {
					fmt.Println("\033[33m", "This letter was already found, try with a different letter.", "\033[0m")
					return false
				} else {
					fmt.Println("\033[33m", "This letter was already tried, try with a different letter.", "\033[0m")
					return false
				}
			} else {
				continue
			}
		}
		return true
	} else {
		fmt.Println("\033[31m", "Error. Please enter only a letter or a word composed of letters only.", "\033[0m")
		return false
	}
}

// PrintWord : Simply prints spaces between letters
func PrintWord(word string) {
	for _, v := range word[:len(word)-1] {
		fmt.Print(string(v))
		fmt.Print(" ")
	}
	fmt.Print(string(word[len(word)-1]))
	fmt.Println()
}

// FindLetter : Test if the letter read in input is on the word or not. If yes, we store the index on the word to print it.
func FindLetter(data savegame.GameData) savegame.GameData {
	data.LetterCheck = false
	letter := strings.ToLower(data.Tries[len(data.Tries)-1])
	letters := strings.Split(data.ToFind, "")
	for i, v := range letters {
		if v == letter {
			data.LetterCheck = true
			data.Index = append(data.Index, i)
		}
	}
	return data
}
