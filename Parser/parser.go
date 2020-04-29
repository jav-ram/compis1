package Parser

import (
	"io/ioutil"

	token "github.com/ram16230/compis1/Token"
)

func FindTokenIndex(tokens []token.Token, tokenID string) int {
	for i, tkn := range tokens {
		if tkn.ID == tokenID {
			return i
		}
	}
	return -1
}

// SplitTokens splits a slice of tokens for splitID
func SplitTokens(tokens []token.Token, splitID string) [][]token.Token {
	definitions := [][]token.Token{}
	definition := []token.Token{}
	for i, tkn := range tokens {
		definition = append(definition, tkn)
		if tkn.ID == splitID {
			if i+1 < len(tokens) && i-1 >= 0 {
				if tokens[i-1].ID != "quote" || tokens[i+1].ID != "quote" {
					definitions = append(definitions, definition)
					definition = []token.Token{}
				}
			} else {
				definitions = append(definitions, definition)
				definition = []token.Token{}
			}
		}
	}
	return definitions
}

func GetDefinitions(tokens []token.Token) [][]token.Token {
	return SplitTokens(tokens, "finish")
}

func GetDefinitionParams(tokens []token.Token) [][]token.Token {
	definitions := [][]token.Token{}
	definition := []token.Token{}
	for _, tkn := range tokens {
		if tkn.ID == "equal" {
			definitions = append(definitions, definition)
			definition = []token.Token{}
			definition = append(definition, tkn)
			definitions = append(definitions, definition)
			definition = []token.Token{}
		} else {
			definition = append(definition, tkn)
		}
	}
	definitions = append(definitions, definition)
	return definitions
}

func initFile() string {
	return `package main

import (
	"fmt"
	"io/ioutil"
	scanner "github.com/ram16230/compis1/Scanner"
	token "github.com/ram16230/compis1/Token"
)

func main() {
`
}

func finishFile() string {
	return `
scan := scanner.MakeAFNS(tkns)
data, _ := ioutil.ReadFile("./test/DoubleAritmetica.ATG")

tokens := scan.Simulate(string(data))
fmt.Printf("%v\n", tokens)
}
`
}

func Parse(tokens []token.Token) {
	text := initFile()
	char, _ := ParseCharacters(GetCharacters(tokens))
	keywords := ParseKeywords(GetKeywords(tokens))
	tkns := ParseTokens(GetTokens(tokens))
	text = text + "\n" + char + "\n" + keywords + "\n" + tkns + "\n" + finishFile()
	ioutil.WriteFile("result.go", []byte(text), 0644)
}
