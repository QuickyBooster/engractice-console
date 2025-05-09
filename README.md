# English Learning Console Application

A terminal-based English learning application that helps users practice vocabulary using data from Google Sheets. The application features dynamic word selection, audio pronunciation, and progress tracking.

## Features

- Practice English words with Vietnamese translations
- Audio pronunciation for words (requires mpg123)
- Points-based learning system
- Integration with Google Sheets for data storage
- Daily point decrease system to encourage regular practice
- Real-time feedback during practice sessions

## Prerequisites

- Go 1.21 or later
- Google Sheets API credentials

- For linux:
   - package:
      - mpg123 (for audio playback)
      - gcc
      - pkgconf
   - enable:
      - `CGO_ENABLED=1`
## Setup

1. Clone the repository
2. Place your Google Sheets credentials file (`credentials.json`) in the project root
3. Run `make build` to build the application
4. Run `make run` to start the application (this will install mpg123 if not present)

## Usage

### Main Menu Options:

1. **Create Test by Lowest Points**
   - Select number of words to practice
   - Type English translations for Vietnamese words
   - Get immediate feedback
   - Points increase by 2 for correct answers

2. **Refresh Data**
   - Update word list from Google Sheets

3. **Update Data**
   - Save progress back to Google Sheets

4. **Decrease Points Daily**
   - Decrease points for all words
   - Helps maintain regular practice schedule

5. **Quit**
   - Exit the application

### During Practice:

- Type the English translation for the displayed Vietnamese word
- Audio pronunciation plays automatically (if enabled)
- Correct answers are detected automatically
- Press Enter to skip a word
- Use the Quit button to return to the main menu

## Google Sheets Format

The application expects a Google Sheet with the following columns:
1. English word
2. Vietnamese translation
3. Audio URL (MP3)
4. Tag
5. Points

## Development

To modify or extend the application:

1. Make your changes in the `cmd/main.go` file
2. Run `make build` to rebuild
3. Run `make run` to test your changes

## License

This project is licensed under the MIT License.