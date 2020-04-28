package Parser

import (
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
	for _, tkn := range tokens {
		definition = append(definition, tkn)
		if tkn.ID == splitID {
			definitions = append(definitions, definition)
			definition = []token.Token{}
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

func Parse(tokens []token.Token) {
	ParseCharacters(GetCharacters(tokens))
	ParseKeywords(GetKeywords(tokens))
}
