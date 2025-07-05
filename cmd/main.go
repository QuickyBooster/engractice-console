package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os/exec"
	"sort"
	"strings"
	"time"

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
	// Read data from the sheet
	spreadsheetId := "1_xKMjnfCG3ADEH5nz5JOqvsFsdQ7UVPmc2ZDBtpvoc8"
	rangeData := "vocabulary!A2:E"

	// Initialize Sheets API client
	// service, err := sheets.NewService(nil, option.WithAuthCredentials())
	service, err := sheets.NewService(nil, option.WithCredentialsFile("credentials.json"))
	if err != nil {
		log.Fatalf("Unable to create Sheets client: %v", err)
	}

	resp, err := service.Spreadsheets.Values.Get(spreadsheetId, rangeData).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	//----------------------------------
	// service.Spreadsheets.
	//----------------------------------
	words := parseSheetData(resp.Values)
	fmt.Printf("Loaded %d words from the sheet.\n", len(words))

	// Initialize UI
	app := tview.NewApplication()
	mainMenu = tview.NewList()
	mainMenu.SetBorder(true).SetTitle("Main Menu")
	// mainMenu = tview.NewList().SetBorder(true).SetTitle("Options")
	mainMenu.AddItem("Create Test by Lowest Points", "", 'c', func() {
		createTestByLowestPoints(words, app)
	})
	mainMenu.AddItem("Refresh Data", "", 'r', func() {
		words = refreshData(app, service, spreadsheetId, rangeData)
	})
	//	Update data
	mainMenu.AddItem("Update Data", "", 'u', func() {
		updateData(app, words, service, spreadsheetId, rangeData)
	})
	//  Decrease points
	mainMenu.AddItem("Decrease Points Daily", "", 'd', func() {
		decreasePointsDaily(app, words)
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
func decreasePointsDaily(app *tview.Application, words []Word) {
	for i := range words {
		words[i].Point--
		if words[i].Point < 0 {
			words[i].Point = 0
		}
	}
	// Show screen
	modal := tview.NewModal()
	resultText := fmt.Sprintf("Decreased point for %d words!", len(words))

	modal.SetText(resultText).
		AddButtons([]string{"Back to Menu"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(mainMenu, true)
		})

	app.SetRoot(modal, true)
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
		startTest(words[:numWords], app, false)
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
	english := tview.NewTextView().SetText("").SetLabel("English: ")
	input := tview.NewInputField().SetLabel("Your Answer: ")

	if audioEnabled {
		go playAudio(testWords[index].MP3)
	}

	// Keep track of whether we're done with the test
	testFinished := false

	input.SetChangedFunc(func(text string) {
		// Skip processing if test is already finished
		if testFinished {
			return
		}

		if strings.EqualFold(text, testWords[index].English) {
			go playAudio(testWords[index].MP3)
			time.Sleep(1 * time.Second)
			testWords[index].Point += 2
			correctCount++
			index++

			if index < len(testWords) {
				// Move to next word
				vietnamese.SetText(testWords[index].Vietnamese)
				// wait for 1 second before clearing the input field
				input.SetText("")
				english.SetText("")
				if audioEnabled {
					go playAudio(testWords[index].MP3)
				}
			} else {
				// All words completed
				testFinished = true
				showResults(correctCount, testWords, app)
			}
		}
	})

	form.AddFormItem(vietnamese)
	form.AddFormItem(english)
	form.AddFormItem(input)
	form.AddButton("Play audio again", func() {
		go playAudio(testWords[index].MP3)
	})
	form.AddButton("Show Answer", func() {
		english.SetText(testWords[index].English)
	})
	form.AddButton("Turn on/off audio", func() {
		audioEnabled = !audioEnabled
	})
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

	cmd := exec.Command("mpg123", fileURL)
	if err := cmd.Run(); err != nil {
		log.Printf("Error playing audio: %v", err)
	}
}
func updateData(app *tview.Application, words []Word, service *sheets.Service, sheetId, rangeData string) {

	// Update the sheet with the new data
	var vr sheets.ValueRange
	for _, word := range words {
		vr.Values = append(vr.Values, []interface{}{
			word.English,
			word.Vietnamese,
			word.MP3,
			word.Tag,
			word.Point,
		})
	}

	_, err := service.Spreadsheets.Values.Update(sheetId, rangeData, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to update data in sheet: %v", err)
	}

	// Show screen
	modal := tview.NewModal()
	resultText := fmt.Sprintf("Updated %d words!", len(words))

	modal.SetText(resultText).
		AddButtons([]string{"Back to Menu"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(mainMenu, true)
		})

	app.SetRoot(modal, true)

}
func refreshData(app *tview.Application, service *sheets.Service, sheetId, rangeData string) []Word {
	// Read data from the sheet
	resp, err := service.Spreadsheets.Values.Get(sheetId, rangeData).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// Parse the new data
	words := []Word{}
	words = parseSheetData(resp.Values)

	// Show screen
	modal := tview.NewModal()
	resultText := fmt.Sprintf("Refreshed %d words!", len(words))

	modal.SetText(resultText).
		AddButtons([]string{"Back to Menu"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(mainMenu, true)
		})

	app.SetRoot(modal, true)

	return words
}
