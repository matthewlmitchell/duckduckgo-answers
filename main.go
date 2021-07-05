package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Options specifies all possible API arguments to be passed into the query URL
type Options struct {
	Format       string
	Pretty       int
	NoRedirect   int
	NoHTML       int
	SkipDisambig int
}

// Response specifies the exact json structure of a generic API query
type Response struct {
	Abstract       string         `json:"Abstract"`
	AbstractSource string         `json:"AbstractSource"`
	AbstractText   string         `json:"AbstractText"`
	AbstractURL    string         `json:"AbstractURL"`
	RelatedTopics  []RelatedTopic `json:"RelatedTopics"`
}

// TerminalColors is a short list of strings to pass to fmt.Println()
// to change the color of text in the terminal
var TerminalColors = map[string]string{
	"Reset":  "\033[0m",
	"Red":    "\033[31m",
	"Green":  "\033[33m",
	"Blue":   "\033[34m",
	"White":  "\033[37m",
	"Yellow": "\033[33m",
}

// RelatedTopic describes the structure of the underlying map[string]
// inside of the query response at "RelatedTopics": [{}]
type RelatedTopic struct {
	FirstURL string `json:"FirstURL"`
	Text     string `json:"Text"`
}

// searchPrompt() prompts the user for DuckDuckGo search query
func searchPrompt() string {
	fmt.Print("\nSearch: ")

	inputReader := bufio.NewReader(os.Stdin)
	query, _ := inputReader.ReadString('\n')
	return query
}

// getAPIURL() formats and returns a string for querying the DuckDuckGo API with http.Get()
func getAPIURL(queryString string, options Options) string {

	queryString = url.QueryEscape(queryString)

	return fmt.Sprintf("https://api.duckduckgo.com/?q=%s&format=%s&pretty=%d&no_redirect=%d&no_html=%d&skip_disambig=%d&t=duckduckgo-answers", queryString, options.Format, options.Pretty, options.NoRedirect, options.NoHTML, options.SkipDisambig)
}

func queryAPI(apiURL string) *http.Response {

	response, err := http.Get(apiURL)
	if err != nil {
		panic(err)
	}

	return response
}

func responseToString(response *http.Response) string {
	stringResponse := make([]string, 1)

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		stringResponse = append(stringResponse, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	response.Body.Close()

	return strings.Join(stringResponse[:], "")
}

func unmarshalResponse(jsonInput string) Response {
	jsonData := Response{}

	jsonBytes := []byte(jsonInput)

	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		panic(err)
	}

	return jsonData
}

func printResponse(input Response) {

	fmt.Printf("\n %s \n \n", input.AbstractText)
	fmt.Println(TerminalColors["Green"], "More info: ")

	for key := range input.RelatedTopics {
		fmt.Println(TerminalColors["Blue"], "\t"+input.RelatedTopics[key].FirstURL)
		fmt.Println(TerminalColors["White"], "\t"+input.RelatedTopics[key].Text+"\n")
	}
}

func main() {
	queryOptions := &Options{
		Format:       "json",
		Pretty:       1,
		NoRedirect:   1,
		NoHTML:       1,
		SkipDisambig: 1,
	}

	for {
		// Ask the user for a search query
		userInput := searchPrompt()

		// Encode the users input query into URL format, and return the formatted API url
		queryURL := getAPIURL(userInput, *queryOptions)

		// Use http.Get() to retrieve an HTTP response for our query
		apiResponse := queryAPI(queryURL)

		// Read the response into our buffer reader then combine it into a single string
		stringAnswer := responseToString(apiResponse)

		// Unmarshal the JSON-encoded string into our Response{} data structure
		parsedResponse := unmarshalResponse(stringAnswer)

		// Nicely print the response data
		printResponse(parsedResponse)
	}
}
