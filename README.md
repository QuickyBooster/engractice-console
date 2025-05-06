# English Learning Console Application

A terminal-based English learning application that helps users practice vocabulary using data from Google Sheets. The application features dynamic word selection, audio pronunciation, and progress tracking.

## Features

- **Test by Lowest Points**: Practice words that need the most attention based on a point system
- **Audio Pronunciation**: Option to play audio pronunciations of words during practice
- **Progress Tracking**: Points system that adjusts based on correct/incorrect answers
- **Daily Point Decrease**: Automatic decrease of word points to encourage regular practice
- **Google Sheets Integration**: Read and update vocabulary data from Google Sheets

## Requirements

- Go 1.21 or later
- mpg123 (for audio playback)
- Google Sheets API credentials

## Setup

1. Ensure you have Go installed on your system
2. Install mpg123 for audio playback:
   - Windows: Download from the official mpg123 website
   - Linux: `sudo apt-get install mpg123`
   - macOS: `brew install mpg123`

3. Place your Google Sheets API credentials in `credentials.json` at the root of the project

4. Build the application:
   ```
   make build
   ```

This will create an executable file `app.exe`.

## Usage

Run the application:
```
./app.exe
```

### Main Menu Options:

1. **Create Test by Lowest Points**: 
   - Select number of words to practice
   - Type English translations for Vietnamese words
   - Get immediate feedback on correct/incorrect answers

2. **Refresh Data**: 
   - Update the local word list from Google Sheets

3. **Update Data**: 
   - Save progress and point changes back to Google Sheets

4. **Decrease Points Daily**: 
   - Decrease points for all words to maintain practice schedule

5. **Quit**: 
   - Exit the application

### During Practice:

- Vietnamese words are displayed one at a time
- Type the English translation
- Correct answers are automatically detected
- Points are updated based on performance:
  - Correct answer: +2 points
  - Skip word: No point change
- Audio pronunciation is played when available

## Google Sheets Format

The application expects the following columns in your Google Sheet:
1. English word
2. Vietnamese translation
3. Audio URL (MP3)
4. Tag
5. Points

## License

This project is licensed under the MIT License - see the LICENSE file for details.