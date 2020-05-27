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

func GetDefinitionName(tokens []token.Token) string {
	for _, tkn := range tokens {
		if tkn.ID == "ident" {
			return tkn.Lexema
		}
	}
	return "unk"
}

func GetDefinitionParams(tokens []token.Token) [][]token.Token {
	definitions := [][]token.Token{}
	definition := []token.Token{}
	find := false
	for _, tkn := range tokens {
		if tkn.ID == "equal" && !find {
			definitions = append(definitions, definition)
			definition = []token.Token{}
			definition = append(definition, tkn)
			definitions = append(definitions, definition)
			definition = []token.Token{}
			find = true
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
	"strconv"

	scanner "github.com/ram16230/compis1/Scanner"
	token "github.com/ram16230/compis1/Token"
)

var list []token.Token = []token.Token{}
var index int = 0

func Expect(item interface{}) bool {
	i := index
	next := index
	if next < len(list) {
		switch v := item.(type) {
		case string:
			{
				index = i
				return list[next].Lexema == v
			}
		case token.Token:
			{
				index = i
				return v.ID == list[next].ID
			}
		case bool:
			{
				return v
			}
		}
	}
	index = i
	return false
}

func Read(item interface{}) bool {
	status := true

	if index < len(list) {
		switch v := item.(type) {
		case string:
			{
				status = list[index].Lexema == v
			}
		case token.Token:
			{
				status = v.ID == list[index].ID
			}
		}
	} else {
		status = false
	}

	if status {
		index++
	}

	return status
}

`
}

func finishFile() string {
	return `
scan := scanner.MakeAFNS(tkns)
data, _ := ioutil.ReadFile("./test/test.txt")

tokens := scan.Simulate(string(data))
list = tokens
fmt.Printf("%v\n", tokens)
}
`
}

func Parse(tokens []token.Token) {
	text := initFile()
	char, _ := ParseCharacters(GetCharacters(tokens))
	keywords := ParseKeywords(GetKeywords(tokens))
	tkns, idsTkns := ParseTokens(GetTokens(tokens))
	production := ParseProductions(GetProductions(tokens), idsTkns) // TODO: change it to keywords, and tkns
	text = text + production + "func main() {\n" + char + "\n" + keywords + "\n" + tkns + "\n" + finishFile()
	ioutil.WriteFile("result.go", []byte(text), 0644)
}
