package Parser

import (
	"fmt"
	"strconv"

	token "github.com/ram16230/compis1/Token"
)

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '\'' && s[len(s)-1] == '\'' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func removeFromSet(set string, minus string) string {
	newText := ""
	for i := 0; i < len(set); i++ {
		isIn := false
		s := rune(set[i])
		for _, m := range minus {
			if m == s || s == '|' {
				isIn = true
				break
			}
		}
		if !isIn {
			newText = newText + string(set[i])
		}
	}

	return fromStringToSet(newText)
}

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

func CleanCharactes(tokens []token.Token) []token.Token {
	newToken := []token.Token{}
	for _, tkn := range tokens {
		if tkn.ID != "any" {
			newToken = append(newToken, tkn)
		}
	}
	return newToken
}

func ParseCharactersDefinitions(tokens []token.Token, chars map[string]string) string {
	for {
		if FindTokenIndex(tokens, "chr") != -1 {
			i := FindTokenIndex(tokens, "chr")
			tmp, _ := strconv.Atoi(tokens[i+2].Lexema)
			//char := trimQuotes(strconv.QuoteRune(rune(tmp)))
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
		} else if FindTokenIndex(tokens, "complete") != -1 {
			i := FindTokenIndex(tokens, "complete")
			start := rune(tokens[i-1].Lexema[0])
			end := rune(tokens[i+1].Lexema[0])
			newText := ""

			for i := int(start); i <= int(end); i++ {
				if i == int(end) {
					newText = newText + string(rune(i))
				} else {
					newText = newText + string(rune(i)) + "|"
				}
			}

			newToken := token.NewToken("set", newText)

			tokens[i-1] = newToken
			tokens = removeFromSlice(tokens, i+1)
			tokens = removeFromSlice(tokens, i)

		} else if FindTokenIndex(tokens, "sum") != -1 {
			i := FindTokenIndex(tokens, "sum")
			// TODO: target to idents, get ident of another character
			l := ""
			r := ""
			if tokens[i-1].ID == "ident" {
				l = chars[tokens[i-1].Lexema]
			} else if tokens[i-1].ID == "set" {
				l = tokens[i-1].Lexema
			}

			if tokens[i+1].ID == "ident" {
				r = chars[tokens[i+1].Lexema]
			} else if tokens[i+1].ID == "set" {
				r = tokens[i+1].Lexema
			}
			newToken := token.NewToken("set", fmt.Sprintf("%v|%v", l, r))

			tokens[i-1] = newToken
			tokens = removeFromSlice(tokens, i+1)
			tokens = removeFromSlice(tokens, i)
		} else if FindTokenIndex(tokens, "subtract") != -1 {
			i := FindTokenIndex(tokens, "subtract")
			// TODO: target to idents, get ident of another character
			l := ""
			r := ""
			if tokens[i-1].ID == "ident" {
				l = chars[tokens[i-1].Lexema]
			} else if tokens[i-1].ID == "set" {
				l = tokens[i-1].Lexema
			}

			if tokens[i+1].ID == "ident" {
				r = chars[tokens[i+1].Lexema]
			} else if tokens[i+1].ID == "set" {
				r = tokens[i+1].Lexema
			}
			s := removeFromSet(l, r)
			newToken := token.NewToken("set", s)

			tokens[i-1] = newToken
			tokens = removeFromSlice(tokens, i+1)
			tokens = removeFromSlice(tokens, i)
		} else {
			return fmt.Sprintf("%v", tokens[0].Lexema)
		}
	}
}

// ParseCharacters parse the tokens inside characters
func ParseCharacters(tokens []token.Token) (string, map[string]string) {
	definitions := GetDefinitions(tokens)
	// id -> production
	chars := map[string]string{}

	for _, def := range definitions {
		params := GetDefinitionParams(def)
		id := GetDefinitionName(params[0])
		fmt.Printf("id %v\n", id)

		charDesc := CleanCharactes(params[2])

		fmt.Printf("def %v\n", charDesc)
		d := ParseCharactersDefinitions(charDesc, chars)
		chars[id] = d

	}
	text := ""
	for i, d := range chars {
		text = text + i + " := " + "`(" + d + ")` \n"
		text = text + "_ = " + i + "\n"
	}

	// fmt.Printf("%v\n", text)

	return text, chars

}
