package Parser

import (
	token "github.com/ram16230/compis1/Token"
)

// isNotTokenID returns true if the source is inside of the tokensIds register and diferent to target
func isNotTokenID(target string, source string) bool {
	tokenIds := []string{"compiler", "characters", "keywords", "tokens", "productions", "end"}

	for _, tokenID := range tokenIds {
		if tokenID == source && target != tokenID {
			return true
		}
	}

	return false
}

// GetSections is a decorator that gets the list of tokens of a certain section of cocol file
func GetSections(tokenID string) func([]token.Token) []token.Token {
	return func(tokens []token.Token) []token.Token {
		newTokens := []token.Token{}
		write := false
		for _, tkn := range tokens {
			if tkn.ID == tokenID {
				write = true
				continue
			} else if isNotTokenID(tokenID, tkn.ID) {
				if write {
					return newTokens
				}
				write = false
			}
			if write {
				newTokens = append(newTokens, tkn)
			}
		}

		return newTokens
	}
}

// GetCompiler get list of tokens of list COMPILER
var GetCompiler = GetSections("compiler")

// GetCharacters get list of tokens of list CHARACTERS
var GetCharacters = GetSections("characters")

// GetKeywords get list of tokens of list KEYWORDS
var GetKeywords = GetSections("keywords")

// GetTokens get list of tokens of list TOKENS
var GetTokens = GetSections("tokens")

// GetProductions get list of tokens of list PRODUCTIONS
var GetProductions = GetSections("productions")
