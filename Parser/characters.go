package Parser

import (
	"fmt"
	"strconv"

	token "github.com/ram16230/compis1/Token"
)

func removeFromSlice(slice []token.Token, i int) []token.Token {
	copy(slice[i:], slice[i+1:])
	slice[len(slice)-1] = token.Token{}
	return slice[:len(slice)-1]
}

func fromStringToSet(text string) string {
	if len(text) > 0 {
		if len(text) == 1 {
			return text
		}
		tmp := ""
		for i, t := range text {
			if i == len(text)-1 {
				tmp = tmp + string(t)
			} else {
				tmp = tmp + string(t) + "|"
			}
		}
		return tmp
	}
	return ""
}

func ParseCharactersDefinitions(tokens []token.Token) string {
	for {
		if FindTokenIndex(tokens, "chr") != -1 {
			i := FindTokenIndex(tokens, "chr")
			tmp, _ := strconv.Atoi(tokens[i+2].Lexema)
			char := string(rune(tmp))
			newToken := token.NewToken("set", char)
			// delete i, i+1, i+2, i+3
			tokens[i] = newToken
			tokens = removeFromSlice(tokens, i+3)
			tokens = removeFromSlice(tokens, i+2)
			tokens = removeFromSlice(tokens, i+1)

		} else if FindTokenIndex(tokens, "quote") != -1 {
			i := FindTokenIndex(tokens, "quote")
			newToken := token.NewToken("set", fromStringToSet(tokens[i+1].Lexema))

			tokens[i] = newToken
			tokens = removeFromSlice(tokens, i+2)
			tokens = removeFromSlice(tokens, i+1)
		} else if FindTokenIndex(tokens, "sum") != -1 {
			i := FindTokenIndex(tokens, "sum")
			// TODO: target to idents, get ident of another character
			newToken := token.NewToken("set", fmt.Sprintf("%v|%v", i-1, i+1))

			tokens[i-1] = newToken
			tokens = removeFromSlice(tokens, i+1)
			tokens = removeFromSlice(tokens, i)
		} else {
			return fmt.Sprintf("%v", tokens[0].Lexema)
		}
	}
}

// ParseCharacters parse the tokens inside characters
func ParseCharacters(tokens []token.Token) string {
	definitions := GetDefinitions(tokens)
	idents := []string{}
	descriptions := []string{}

	for _, def := range definitions {
		params := GetDefinitionParams(def)
		idents = append(idents, params[0][0].Lexema)

		charDesc := params[2]

		descriptions = append(descriptions, ParseCharactersDefinitions(charDesc))

	}
	text := ""
	for i := range definitions {
		text = text + idents[i] + " := " + "\"" + descriptions[i] + "\" \n"
	}

	return text

}
