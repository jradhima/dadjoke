/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Get a random dad joke",
	Long:  `This command fetches a random dad joke from the icanhazdadjoke api.`,
	Run: func(cmd *cobra.Command, args []string) {
		jokeTerm, err := cmd.Flags().GetString("term")
		if err != nil {
			fmt.Printf("Error encountered - %v", err)
		}
		switch jokeTerm {
		case "":
			getRandomJoke()
		default:
			getRandomJokeWithTerm(jokeTerm)
		}

	},
}

func init() {
	rootCmd.AddCommand(randomCmd)

	randomCmd.PersistentFlags().String("term", "", "A search term for a dad joke.")
}

type Joke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

type SearchResult struct {
	Results    json.RawMessage `json:"results"`
	SearchTerm string          `json:"search_term"`
	Status     int             `json:"status"`
	TotalJokes int             `json:"total_jokes"`
}

var url string = "https://icanhazdadjoke.com"

func getJokeDataWithTerm(jokeTerm string) (int, []Joke) {
	url := fmt.Sprintf("%s/search?term=%s", url, jokeTerm)
	responseBytes, err := getJokeData(url)
	if err != nil {
		fmt.Println(err)
	}
	results := SearchResult{}
	err = json.Unmarshal(responseBytes, &results)
	if err != nil {
		fmt.Println(err)
	}
	jokes := []Joke{}
	err = json.Unmarshal(results.Results, &jokes)
	if err != nil {
		fmt.Println(err)
	}
	return results.TotalJokes, jokes
}

func getRandomJokeWithTerm(jokeTerm string) {
	_, jokes := getJokeDataWithTerm(jokeTerm)
	if len(jokes) <= 0 {
		fmt.Printf("Couldn't find any jokes with term: %s\n", jokeTerm)
		return
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	jokeNum := r1.Intn(len(jokes))
	fmt.Println(jokes[jokeNum].Joke)
}

func getRandomJoke() {
	responseBytes, err := getJokeData(url)
	if err != nil {
		fmt.Println(err)
	}
	joke := Joke{}
	err = json.Unmarshal(responseBytes, &joke)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(joke.Joke)
}

func getJokeData(baseAPI string) ([]byte, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		baseAPI,
		nil,
	)
	if err != nil {
		fmt.Errorf("Could not request a datjoke - %v", err)
		return nil, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "CLI tutorial (github.com/jradhima/dadjoke)")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Errorf("Could not complete dadjoke request - %v", err)
		return nil, err
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Errorf("Could not parse response body - %v", err)
		return nil, err
	}

	return responseBytes, nil
}
