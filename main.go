package main

import (
	"fmt"
	"hang_web/game"
	"hang_web/position"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// GameData holds the state of the game
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
	Level            int
	Lang             string
}

var forbiddenChars = []rune{'@', '#', '$', '%', '&', '*', '(', ')', '-', '_', '+', '=', '[', ']', '{', '}', '|', '\\', '/', '<', '>', ',', '.', '?', '!', ';', ':', '"', '\'', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

// Checks if the input contains any forbidden characters
func containsForbiddenChar(input string) bool {
	for _, c := range input {
		for _, forbidden := range forbiddenChars {
			if c == forbidden {
				return true
			}
		}
	}
	return false
}

// HomeHandler serves the homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/index.html")
}

// GameHandler initializes the game based on user selection
func GameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		level := r.FormValue("level")
		mode := r.FormValue("mode")
		lang := r.FormValue("lang")
		var fileName string
		switch level {
		case "1":
			fileName = fmt.Sprintf("displaytxt/words%s.txt", getLangSuffix(lang))
		case "2":
			fileName = fmt.Sprintf("displaytxt/words2%s.txt", getLangSuffix(lang))
		case "3":
			fileName = fmt.Sprintf("displaytxt/words3%s.txt", getLangSuffix(lang))
		default:
			http.Error(w, "Invalid level choice", http.StatusBadRequest)
			return
		}

		words := game.GetFile(fileName)
		if len(words) == 0 {
			http.Error(w, "File is empty or contains invalid data", http.StatusInternalServerError)
			return
		}

		randomWord := game.TakeRandomWord(words)
		data := GameData{
			Word:             strings.Repeat("_", len([]rune(randomWord))),
			ToFind:           randomWord,
			Attempts:         10,
			Tries:            []string{},
			HangmanPositions: 0,
			HangmanGraphic:   position.OpenJose()[0],
			Level:            getLevelInt(level),
			Lang:             lang,
		}

		var templateFile string
		switch mode {
		case "1":
			templateFile = "./html/hangman.html"
		case "2":
			templateFile = "./html/hangmanInfinite.html"
		default:
			http.Error(w, "Invalid mode choice", http.StatusBadRequest)
			return
		}

		renderTemplate(w, templateFile, data)
	} else if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./html/game.html")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get language suffix for filename based on user choice
func getLangSuffix(lang string) string {
	switch lang {
	case "en":
		return "_en"
	case "es":
		return "_es"
	default:
		return ""
	}
}

// Convert level string to integer
func getLevelInt(level string) int {
	levelInt, _ := strconv.Atoi(level)
	return levelInt
}

// hangmanHandler initializes the base for the game mode basic
func hangmanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		handleHTTPError(w, err)
		return
	}

	var data GameData
	data.ToFind = r.FormValue("toFind")
	data.Word = r.FormValue("word")
	data.Attempts, _ = strconv.Atoi(r.FormValue("attempts"))
	data.HangmanPositions, _ = strconv.Atoi(r.FormValue("hangmanPositions"))
	if tries := r.FormValue("tries"); tries != "" {
		data.Tries = strings.Split(tries, ",")
	}
	data.HangmanGraphic = r.FormValue("hangmanGraphic")
	letter := strings.TrimSpace(r.FormValue("letter"))
	if letter == "" {
		data.Message = "Veuillez entrer une lettre."
		renderTemplate(w, "html/hangman.html", data)
		return
	}
	if containsForbiddenChar(letter) {
		data.Message = "Les caractères spéciaux (@, #, $, %, etc.) et les chiffres ne sont pas autorisés."
		renderTemplate(w, "html/hangman.html", data)
		return
	}
	if len([]rune(letter)) > 1 {
		if strings.EqualFold(letter, data.ToFind) {
			data.Word = data.ToFind
			data.Message = "Bravo ! Vous avez trouvé le mot complet !"
		} else {
			data.Attempts -= 2
			data.HangmanPositions += 2
			if data.HangmanPositions >= len(position.OpenJose()) {
				data.HangmanPositions = len(position.OpenJose()) - 1
			}
			data.HangmanGraphic = position.OpenJose()[data.HangmanPositions]
			data.Message = "Mot incorrect ! Vous perdez 2 tentatives."
		}
	} else {
		if strings.Contains(strings.Join(data.Tries, ""), letter) {
			data.Message = "La lettre a déjà été essayée, essayez une autre lettre."
		} else {
			data.Tries = append(data.Tries, letter)
			if strings.Contains(data.ToFind, letter) {
				for i, ch := range data.ToFind {
					if string(ch) == letter {
						data.Word = data.Word[:i] + string(ch) + data.Word[i+1:]
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
	}

	if data.Word == data.ToFind {
		http.ServeFile(w, r, "./html/victory.html")
		return
	} else if data.Attempts <= 0 {
		http.ServeFile(w, r, "./html/defeat.html")
		return
	}
	renderTemplate(w, "html/hangman.html", data)
}

// hangmanInfiniteHandler initializes the base for the game mode infinite
func hangmanInfiniteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var data GameData
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	letter := normalizeInput(r.FormValue("letter"))
	if letter == "" {
		http.Error(w, "No letter provided", http.StatusBadRequest)
		return
	}
	wordToFind := r.FormValue("wordToFind")
	currentWord := r.FormValue("currentWord")
	attemptsStr := r.FormValue("attempts")
	triesStr := r.FormValue("tries")
	successCountStr := r.FormValue("successCount")
	levelStr := r.FormValue("level")
	langStr := r.FormValue("lang")
	attempts, _ := strconv.Atoi(attemptsStr)
	successCount, _ := strconv.Atoi(successCountStr)
	level, _ := strconv.Atoi(levelStr)
	tries := strings.Split(triesStr, ",")
	if len(tries) == 1 && tries[0] == "" {
		tries = []string{}
	}

	for _, t := range tries {
		if t == letter {
			data = GameData{
				Word:             currentWord,
				ToFind:           wordToFind,
				Attempts:         attempts,
				Tries:            tries,
				SuccessCount:     successCount,
				Level:            level,
				Lang:             langStr,
				HangmanPositions: 10 - attempts,
				HangmanGraphic:   position.OpenJose()[10-attempts],
			}
			data.Message = "La lettre a déjà été essayée, essayez une autre lettre."
			renderTemplate(w, "html/hangmanInfinite.html", data)
			return
		}
	}
	tries = append(tries, letter)
	if strings.Contains(wordToFind, letter) {
		newWord := []rune(currentWord)
		for i, char := range wordToFind {
			if string(char) == letter {
				newWord[i] = char
			}
		}
		currentWord = string(newWord)
	} else {
		attempts--
	}

	if string([]rune(currentWord)) == wordToFind {
		successCount++
		attempts += 3
		fileName := fmt.Sprintf("displaytxt/words%s%s.txt", getLevelSuffix(level), getLangSuffix(langStr))
		words := game.GetFile(fileName)
		wordToFind = game.TakeRandomWord(words)
		currentWord = strings.Repeat("_", len([]rune(wordToFind)))
		tries = []string{}
	}

	if attempts <= 0 {
		data = GameData{
			ToFind:       wordToFind,
			SuccessCount: successCount,
			Level:        level,
			Lang:         langStr,
		}
		renderTemplate(w, "./html/scoreboard.html", data)
		return
	}
	data = GameData{
		Word:             currentWord,
		ToFind:           wordToFind,
		Attempts:         attempts,
		Tries:            tries,
		SuccessCount:     successCount,
		Level:            level,
		Lang:             langStr,
		HangmanPositions: 10 - attempts,
		HangmanGraphic:   position.OpenJose()[10-attempts],
	}

	renderTemplate(w, "html/hangmanInfinite.html", data)
}

// generate a suffix for file names based on the game level
func getLevelSuffix(level int) string {
	if level == 1 {
		return ""
	}
	return strconv.Itoa(level)
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
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		data.ToFind = r.FormValue("wordToFind")
		successCountStr := r.FormValue("successCount")
		data.SuccessCount, _ = strconv.Atoi(successCountStr)
		data.Level, _ = strconv.Atoi(r.FormValue("level"))
		data.Lang = r.FormValue("lang")
		renderTemplate(w, "./html/scoreboard.html", data)
	} else {
		renderTemplate(w, "./html/scoreboard.html", data)
	}
}

// NormalizeInput nettoie et normalise l'entrée utilisateur
func normalizeInput(input string) string {
	input = strings.TrimSpace(input)
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)
	result, _, _ := transform.String(t, input)

	return result
}

// Renders the provided template with data
func renderTemplate(w http.ResponseWriter, templatePath string, data GameData) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

// Main function that sets up routes and starts the web server
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

	http.ListenAndServe(":8089", nil)
}
