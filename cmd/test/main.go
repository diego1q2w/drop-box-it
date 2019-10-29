package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type ratedQuestion struct {
	question question
	rate     int
}

type question struct {
	ID              uint64   `json:"id"`
	Content         string   `json:"content"`
	CreateTimestamp uint64   `json:"createTimestamp"`
	Answers         []answer `json:"answers"`
}

type answer struct {
	ID      uint64 `json:"id"`
	Rating  int    `json:"rating"`
	Content string `json:"content"`
}

func getContent() []question {
	jsonContent := []byte(``)
	var questions []question
	if err := json.Unmarshal(jsonContent, &questions); err != nil {
		panic(err)
	}

	return questions
}

func getMaximumRateOfQuestion(question question) int {
	var maxRate = 0
	for _, a := range question.Answers {
		if a.Rating > maxRate {
			maxRate = a.Rating
		}
	}

	return maxRate
}

func findDuplicates(questions []question) map[string][]ratedQuestion {
	dedQuestions := make(map[string][]ratedQuestion)
	for _, q := range questions {
		rate := getMaximumRateOfQuestion(q)
		content := strings.ToLower(q.Content)
		_, ok := dedQuestions[content]
		if !ok {
			dedQuestions[content] = []ratedQuestion{{question: q, rate: rate}}
		} else {
			rq := ratedQuestion{question: q, rate: rate}
			dedQuestions[content] = append(dedQuestions[q.Content], rq)
		}
	}

	return dedQuestions
}

func deduplicate(questions []question) []question {
	dedQuestions := findDuplicates(questions)

	var finalQuestions []question
	for _, questions := range dedQuestions {
		sort.Slice(questions, func(i, j int) bool {
			if questions[i].rate == questions[j].rate {
				return questions[i].question.CreateTimestamp < questions[j].question.CreateTimestamp
			}
			return questions[i].rate > questions[j].rate
		})

		finalQuestions = append(finalQuestions, questions[0].question)
	}

	sort.Slice(finalQuestions[:], func(i, j int) bool {
		return finalQuestions[i].ID < finalQuestions[j].ID
	})

	return finalQuestions
}

func main() {
	content := getContent()

	deduplicated := deduplicate(content)
	for _, c := range deduplicated {
		fmt.Println(c)
	}
}
