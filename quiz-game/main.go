package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const (
	defaultLimit         = 10
	defaultFile          = "problems.csv"
	limitFlagDescription = "Limit in seconds to complete the entire quiz."
	fileFlagDescription  = "Input CSV file to read questions from."
)

type Question struct {
	Question string
	Answer   string
}

func (question Question) String() string {
	return fmt.Sprintf("Question: %s. Answer: %s\n", question.Question, question.Answer)
}

type Quiz struct {
	Questions []Question
}

func (quiz *Quiz) AddQuestion(question Question) {
	quiz.Questions = append(quiz.Questions, question)
}

func parseQuestions(file *os.File) (Quiz, error) {
	quiz := Quiz{}

	r := csv.NewReader(file)
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Quiz{}, fmt.Errorf("error reading CSV: %w", err)
		}
		if len(line) != 2 {
			log.Printf("Skipping invalid CSV row with %d columns", len(line))
			continue
		}
		quiz.AddQuestion(Question{
			Question: line[0],
			Answer:   strings.TrimSpace(strings.ToLower(line[1])),
		})
	}
	return quiz, nil
}

func playQuiz(quiz Quiz, timerChan <-chan time.Time) (bool, int) {
	score := 0

	for _, question := range quiz.Questions {
		fmt.Print(question.Question, "? ")
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timerChan:
			return false, score
		case answer := <-answerCh:
			userAnswer := strings.TrimSpace(strings.ToLower(answer))
			if question.Answer == userAnswer {
				score++
			}
		}
	}
	return true, score
}

func main() {

	fmt.Println("Quiz starting")

	questionsFile := flag.String("f", defaultFile, fileFlagDescription)
	limit := flag.Int("limit", defaultLimit, limitFlagDescription)
	flag.Parse()

	if *limit < 0 {
		fmt.Println("Provided limit is less than 0. Proceeding with default value.")
		*limit = defaultLimit
	}

	file, err := os.Open(*questionsFile)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	quiz, err := parseQuestions(file)
	if err != nil {
		log.Fatalf("Error parsing questions: %v", err)
	}
	fmt.Println("Get ready...")

	timer := time.NewTimer(time.Duration(*limit) * time.Second)

	completed, score := playQuiz(quiz, timer.C)

	switch completed {
	case true:
		timer.Stop()
		fmt.Printf("\nFinal score: %v / %v\n", score, len(quiz.Questions))
	case false:
		fmt.Printf("\nTime's up! Final score: %v / %v\n", score, len(quiz.Questions))
	}
}
