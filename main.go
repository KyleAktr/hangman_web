package main

import (
	"fmt"
	"hang_web/game"
	"hang_web/position"
	"net/http"
	"strings"
	"text/template"
)

type GameData struct {
	Word             string
	ToFind           string
	Index            []int
	Attempts         int
	Tries            []string
	LetterCheck      bool
	HangmanPositions int
	HangmanGraphic   string
	Message          string
	SuccessCount     int
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/index.html")
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		level := r.FormValue("level")
		var file string

		switch level {
		case "1":
			file = "displaytxt/words.txt"
		case "2":
			file = "displaytxt/words2.txt"
		case "3":
			file = "displaytxt/words3.txt"
		case "4":
			file = "displaytxt/save.txt"
		default:
			http.Error(w, "Invalid level choice", http.StatusBadRequest)
			return
		}
		words := game.GetFile(file)
		if len(words) == 0 {
			http.Error(w, "File is empty or contains invalid data", http.StatusInternalServerError)
			return
		}

		randomWord := game.TakeRandomWord(words)
		data := GameData{
			Word:             strings.Repeat("_", len(randomWord)),
			ToFind:           strings.ToUpper(randomWord),
			Attempts:         10,
			Tries:            []string{},
			HangmanPositions: 0,
			HangmanGraphic:   position.OpenJose()[0],
		}
		mode := r.FormValue("mode")

		switch mode {
		case "1":
			tmpl, err := template.ParseFiles("./html/hangman.html")
			if err != nil {
				http.Error(w, "Error loading hangman template", http.StatusInternalServerError)
				return
			}
			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, "Error rendering hangman template", http.StatusInternalServerError)
			}
		case "2":
			tmpl, err := template.ParseFiles("./html/hangmanInfinite.html")
			if err != nil {
				http.Error(w, "Error loading hangman template", http.StatusInternalServerError)
				return
			}
			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, "Error rendering hangman template", http.StatusInternalServerError)
			}
		default:
			http.Error(w, "Invalid mode choice", http.StatusBadRequest)
			return
		}
		// tmpl, err := template.ParseFiles("./html/hangman.html")
		// if err != nil {
		// 	http.Error(w, "Error loading hangman template", http.StatusInternalServerError)
		// 	return
		// }
		// if err := tmpl.Execute(w, data); err != nil {
		// 	http.Error(w, "Error rendering hangman template", http.StatusInternalServerError)
		// }
	} else if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./html/game.html")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// problème de code trop d'indentation
func hangmanHandler(w http.ResponseWriter, r *http.Request) {
	var data GameData

	if r.Method == http.MethodPost {
		r.ParseForm()
		data.Word = r.FormValue("word")
		data.ToFind = r.FormValue("toFind")
		fmt.Sscanf(r.FormValue("attempts"), "%d", &data.Attempts)
		fmt.Sscanf(r.FormValue("hangmanPositions"), "%d", &data.HangmanPositions)
		tries := r.FormValue("tries")
		if tries != "" {
			data.Tries = strings.Split(tries, ",")
		}
		data.HangmanGraphic = r.FormValue("hangmanGraphic")
		letter := strings.ToUpper(r.FormValue("letter"))
		if len(letter) == 1 {
			if strings.Contains(strings.Join(data.Tries, ""), letter) {
				data.Message = "La lettre a déjà été essayée, essayez une autre lettre."
			} else {
				data.Tries = append(data.Tries, letter)
				if strings.Contains(data.ToFind, letter) {
					for i, ch := range data.ToFind {
						if string(ch) == letter {
							data.Word = data.Word[:i] + letter + data.Word[i+1:]
						}
					}
				} else {
					data.Attempts--
					data.HangmanPositions++
					if data.HangmanPositions < len(position.OpenJose()) {
						data.HangmanGraphic = position.OpenJose()[data.HangmanPositions]
					}
				}
			}
		} else {
			data.Message = "Veuillez entrer une seule lettre."
		}
		if strings.ToUpper(data.Word) == strings.ToUpper(data.ToFind) {
			http.ServeFile(w, r, "./html/victory.html")
			return
		} else if data.Attempts <= 0 {
			http.ServeFile(w, r, "./html/defeat.html")
			return
		}
		//uniquement pour les testes ou guillian
	} else if r.Method == http.MethodGet {
		words := game.GetFile("displaytxt/words.txt")
		if len(words) == 0 {
			http.Error(w, "Default word file is empty or invalid", http.StatusInternalServerError)
			return
		}

		randomWord := game.TakeRandomWord(words)
		data = GameData{
			Word:             strings.Repeat("_", len(randomWord)),
			ToFind:           strings.ToUpper(randomWord),
			Attempts:         10,
			Tries:            []string{},
			HangmanPositions: 0,
			HangmanGraphic:   position.OpenJose()[0],
			SuccessCount:     0,
		}
	}
	tmpl, err := template.ParseFiles("html/hangman.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// problème de code trop d'indentation
func hangmanInfiniteHandler(w http.ResponseWriter, r *http.Request) {
	var data GameData

	if r.Method == http.MethodPost {
		r.ParseForm()
		data.Word = r.FormValue("word")
		data.ToFind = r.FormValue("toFind")
		fmt.Sscanf(r.FormValue("attempts"), "%d", &data.Attempts)
		fmt.Sscanf(r.FormValue("hangmanPositions"), "%d", &data.HangmanPositions)
		fmt.Sscanf(r.FormValue("SuccessCount"), "%d", &data.SuccessCount)
		tries := r.FormValue("tries")
		if tries != "" {
			data.Tries = strings.Split(tries, ",")
		}
		data.HangmanGraphic = r.FormValue("hangmanGraphic")
		letter := strings.ToUpper(r.FormValue("letter"))
		if len(letter) == 1 {
			if strings.Contains(strings.Join(data.Tries, ""), letter) {
				data.Message = "La lettre a déjà été essayée, essayez une autre lettre."
			} else {
				data.Tries = append(data.Tries, letter)
				if strings.Contains(data.ToFind, letter) {
					for i, ch := range data.ToFind {
						if string(ch) == letter {
							data.Word = data.Word[:i] + letter + data.Word[i+1:]
						}
					}
				} else {
					data.Attempts--
					data.HangmanPositions++
					if data.HangmanPositions < len(position.OpenJose()) {
						data.HangmanGraphic = position.OpenJose()[data.HangmanPositions]
					}
				}
			}
		} else {
			data.Message = "Veuillez entrer une seule lettre."
		}
		if strings.ToUpper(data.Word) == strings.ToUpper(data.ToFind) {
			data.SuccessCount++
			data.Attempts += 3
			newWord := game.TakeRandomWord(game.GetFile("displaytxt/words2.txt"))
			data.ToFind = strings.ToUpper(newWord)
			data.Word = strings.Repeat("_", len(newWord))
			data.Tries = []string{}
			data.HangmanPositions = 0
			data.HangmanGraphic = position.OpenJose()[0]
		} else if data.Attempts <= 0 {
			tmpl, err := template.ParseFiles("./html/scoreboard.html")
			if err != nil {
				http.Error(w, "Error loading scoreboard template", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, data)
			return
		}
	} else if r.Method == http.MethodGet {
		words := game.GetFile("displaytxt/words.txt")
		if len(words) == 0 {
			http.Error(w, "Default word file is empty or invalid", http.StatusInternalServerError)
			return
		}

		randomWord := game.TakeRandomWord(words)
		data = GameData{
			Word:             strings.Repeat("_", len(randomWord)),
			ToFind:           strings.ToUpper(randomWord),
			Attempts:         10,
			Tries:            []string{},
			HangmanPositions: 0,
			HangmanGraphic:   position.OpenJose()[0],
			SuccessCount:     0,
		}
	}
	tmpl, err := template.ParseFiles("html/hangmanInfinite.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func ContactHandle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/contact.html")
}

func handleHTTPError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func erreurHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/404.html")
}

func victoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Victory page accessed")
	http.ServeFile(w, r, "./html/victory.html")
}

func DefeatHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/defeat.html")
}

func ScoreboardHandler(w http.ResponseWriter, r *http.Request) {
	var data GameData
	data.SuccessCount = 0
	http.ServeFile(w, r, "./html/scoreboard.html")
}

func setupRoutes() {
	http.HandleFunc("/index", HomeHandler)
	http.HandleFunc("/game", GameHandler)
	http.HandleFunc("/hangman", hangmanHandler)
	http.HandleFunc("/hangmanInfinite", hangmanInfiniteHandler)
	http.HandleFunc("/contact", ContactHandle)
	http.HandleFunc("/victory", victoryHandler)
	http.HandleFunc("/defeat", DefeatHandler)
	http.HandleFunc("/scoreboard", ScoreboardHandler)

}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	setupRoutes()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		erreurHandler(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
