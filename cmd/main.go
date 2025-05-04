package main

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	// "time"

	"github.com/rivo/tview"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Word struct {
	English    string
	Vietnamese string
	MP3        string
	Tag        string
	Point      int
}

// Global menu variable so it can be accessed from other functions
var mainMenu *tview.List

func main() {
	fmt.Println("Welcome to the English Practice App!")

	// Initialize Sheets API client
	service, err := sheets.NewService(nil, option.WithCredentialsFile("credentials.json"))
	if err != nil {
		log.Fatalf("Unable to create Sheets client: %v", err)
	}

	// Read data from the sheet
	spreadsheetId := "1_xKMjnfCG3ADEH5nz5JOqvsFsdQ7UVPmc2ZDBtpvoc8"
	rangeData := "vocabulary!A2:E"
	resp, err := service.Spreadsheets.Values.Get(spreadsheetId, rangeData).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	//----------------------------------
	// service.Spreadsheets.
	//----------------------------------
	words := parseSheetData(resp.Values)
	fmt.Printf("Loaded %d words from the sheet.\n", len(words))

	// Example: Decrease points daily
	decreasePointsDaily(words)

	// Initialize UI
	app := tview.NewApplication()
	mainMenu = tview.NewList()
	mainMenu.SetBorder(true).SetTitle("Main Menu")
	// mainMenu = tview.NewList().SetBorder(true).SetTitle("Options")
	mainMenu.AddItem("Create Test by Lowest Points", "", 'a', func() {
		createTestByLowestPoints(words, app)
	})
	mainMenu.AddItem("Quit", "", 'q', func() {
		app.Stop()
	})

	app.SetRoot(mainMenu, true)

	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

func parseSheetData(data [][]interface{}) []Word {
	var words []Word
	for _, row := range data {
		if len(row) < 5 {
			continue
		}
		point := 0
		if p, ok := row[4].(string); ok {
			fmt.Sscanf(p, "%d", &point)
		}
		words = append(words, Word{
			English:    fmt.Sprintf("%v", row[0]),
			Vietnamese: fmt.Sprintf("%v", row[1]),
			MP3:        fmt.Sprintf("%v", row[2]),
			Tag:        fmt.Sprintf("%v", row[3]),
			Point:      point,
		})
	}
	return words
}

func decreasePointsDaily(words []Word) {
	for i := range words {
		words[i].Point--
		if words[i].Point < 0 {
			words[i].Point = 0
		}
	}
	fmt.Println("Points decreased for all words.")
}

func createTestByLowestPoints(words []Word, app *tview.Application) {
	form := tview.NewForm()
	form.AddInputField("Number of words", "", 10, nil, nil)
	form.AddButton("Start Test", func() {
		numWordsText := form.GetFormItemByLabel("Number of words").(*tview.InputField).GetText()
		numWords := 0
		fmt.Sscanf(numWordsText, "%d", &numWords)
		rand.Shuffle(len(words), func(i, j int) {
			words[i], words[j] = words[j], words[i]
		})
		sort.Slice(words, func(i, j int) bool {
			return words[i].Point < words[j].Point
		})

		if numWords > len(words) {
			numWords = len(words)
		}
		startTest(words[:numWords], app, true)
	})
	form.AddButton("Cancel", func() {
		app.SetRoot(mainMenu, true)
	})

	app.SetRoot(form, true)
}

func startTest(testWords []Word, app *tview.Application, audioEnabled bool) {
	if len(testWords) == 0 {
		log.Println("No words available for the test.")
		app.SetRoot(mainMenu, true)
		return
	}

	index := 0
	correctCount := 0

	form := tview.NewForm()
	vietnamese := tview.NewTextView().SetText(testWords[index].Vietnamese).SetLabel("Vietnamese: ")
	input := tview.NewInputField().SetLabel("Your Answer: ")

	if audioEnabled {
		playAudio(testWords[index].MP3)
	}

	// Keep track of whether we're done with the test
	testFinished := false

	input.SetChangedFunc(func(text string) {
		// Skip processing if test is already finished
		if testFinished {
			return
		}

		if strings.EqualFold(text, testWords[index].English) {
			correctCount++
			index++

			if index < len(testWords) {
				// Move to next word
				vietnamese.SetText(testWords[index].Vietnamese)
				input.SetText("")
				if audioEnabled {
					playAudio(testWords[index].MP3)
				}
			} else {
				// All words completed
				testFinished = true
				showResults(correctCount, testWords, app)
			}
		}
	})

	form.AddFormItem(vietnamese)
	form.AddFormItem(input)
	form.AddButton("Quit", func() {
		app.SetRoot(mainMenu, true)
	})

	app.SetRoot(form, true)
}

func showResults(correctCount int, testWords []Word, app *tview.Application) {
	modal := tview.NewModal()
	resultText := fmt.Sprintf("You got %d/%d correct!", correctCount, len(testWords))

	modal.SetText(resultText).
		AddButtons([]string{"Back to Menu"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(mainMenu, true)
		})

	app.SetRoot(modal, true)
}
func playAudio(fileURL string) {
	resp, err := http.Get(fileURL)
	if err != nil {
		log.Printf("Error fetching audio file: %v", err)
		return
	}
	defer resp.Body.Close()

	tempFile, err := os.CreateTemp("", "audio-*.mp3")
	if err != nil {
		log.Printf("Error creating temporary file: %v", err)
		return
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		log.Printf("Error saving audio file: %v", err)
		return
	}

	tempFile.Close()

	cmd := exec.Command("mpg123", tempFile.Name())
	if err := cmd.Run(); err != nil {
		log.Printf("Error playing audio: %v", err)
	}
}
