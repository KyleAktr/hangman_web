package game

import (
	"fmt"
	"hang_web/savegame"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
)

// GetFile : Take the words.txt file convert it into []string
func GetFile(file string) []string {
	WordFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	WordStr := string(WordFile)                // Transform our file into string
	WordsSlice := strings.Split(WordStr, "\n") // strings.Split function slice a string when it found a parameter, here is "/n"
	return WordsSlice
}

// RandomNumber : Return a random int value
func RandomNumber(i int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomIndex := r1.Intn(i)
	return randomIndex
}

// RandomWord : Takes a Random word in a slice of words
func TakeRandomWord(txt []string) string {
	RandomIndex := RandomNumber(len(txt))
	ToFind := txt[RandomIndex]
	return ToFind
}

// Choosefile : this function chooses the appropriate words file and starts a new game or loads a saved game.
func Choosefile(arg string) {
	fmt.Println(arg)
	var Data savegame.GameData
	if arg == "displaytxt/words.txt" {
		file := "displaytxt/words.txt"
		WordsSlice := GetFile(file)
		if TestFile(WordsSlice) == false {
			return
		}
		Data = StartGame(WordsSlice)
	} else if arg == "displaytxt/words2.txt" {
		file := "displaytxt/words2.txt"
		WordsSlice := GetFile(file)
		if TestFile(WordsSlice) == false {
			return
		}
		Data = StartGame(WordsSlice)
	} else if arg == "displaytxt/words3.txt" {
		file := "displaytxt/words3.txt"
		WordsSlice := GetFile(file)
		if TestFile(WordsSlice) == false {
			return
		}
		Data = StartGame(WordsSlice)
	} else if arg == "displaytxt/save.txt" {
		file := "displaytxt/save.txt"
		Data = savegame.StartWithFlag(file) // Correction : Stocke la valeur retournée
	} else {
		fmt.Println("\033[31m", "Please enter displaytxt/words.txt to start a new game or --startWith save.txt to continue your last game.", "\033[0m")
		return
	}

	fmt.Println(Data.Word)
	Run(arg, Data)
}

// TestFile , tests if the file contain false words or empty lines
func TestFile(file []string) bool {
	for _, word := range file {
		if len(word) > 0 {
			for _, letter := range strings.ToUpper(word) {
				if letter < 65 || letter > 90 {
					fmt.Println("\033[31m", "This file must contain only letters. Please modify the words.txt file.")
					return false
				} else {
					continue
				}
			}
			continue
		} else {
			fmt.Println("\033[31m", "The words.txt file contain empty lines. Please modify it.")
			return false
		}
	}
	return true
}