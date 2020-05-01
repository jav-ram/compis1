package Parser

import (
	token "github.com/ram16230/compis1/Token"
)

func parseTokenDefinition(tokens []token.Token) string {
	for {
		if FindTokenIndex(tokens, "except") != -1 {
			i := FindTokenIndex(tokens, "except")
			n := token.NewToken("string", "")
			tokens[i] = n
		} else if FindTokenIndex(tokens, "quote") != -1 {
			start := FindTokenIndex(tokens, "quote")
			text := ""
			end := -1

			for i := start + 1; i < len(tokens); i++ {
				if tokens[i].ID == "quote" {
					end = i
					break
				} else {
					text = text + tokens[i].Lexema
				}
			}

			newToken := token.NewToken("string", text)

			tokens[start] = newToken

			for i := end; i > start; i-- {
				tokens = removeFromSlice(tokens, i)
			}
		} else if FindTokenIndex(tokens, "ident") != -1 {
			i := FindTokenIndex(tokens, "ident")
			newToken := token.NewToken("var", tokens[i].Lexema)
			tokens[i] = newToken
		} else if FindTokenIndex(tokens, "or") != -1 {
			i := FindTokenIndex(tokens, "or")

			l := ""
			r := ""
			if tokens[i-1].ID == "ident" {
				l = tokens[i-1].Lexema
				n := token.NewToken("var", l)
				tokens[i-1] = n
			} else if tokens[i-1].ID == "string" {
				l = tokens[i-1].Lexema
				n := token.NewToken("string", l)
				tokens[i-1] = n
			}

			if tokens[i+1].ID == "ident" {
				r = tokens[i+1].Lexema
				n := token.NewToken("var", r)
				tokens[i+1] = n
			} else if tokens[i+1].ID == "string" {
				r = tokens[i+1].Lexema
				n := token.NewToken("string", r)
				tokens[i+1] = n
			}

			newToken := token.NewToken("string", "|")
			tokens[i] = newToken
		} else if FindTokenIndex(tokens, "group_start") != -1 {
			start := FindTokenIndex(tokens, "group_start")
			end := FindTokenIndex(tokens, "group_end")

			startToken := token.NewToken("string", tokens[start].Lexema)
			endToken := token.NewToken("string", tokens[end].Lexema)

			tokens[start] = startToken
			tokens[end] = endToken
		} else if FindTokenIndex(tokens, "option_start") != -1 {
			start := FindTokenIndex(tokens, "option_start")
			end := FindTokenIndex(tokens, "option_end")

			startToken := token.NewToken("string", "(")
			endToken := token.NewToken("string", "|`)") // FIXME: epsilon

			tokens[start] = startToken
			tokens[end] = endToken
		} else if FindTokenIndex(tokens, "iteration_start") != -1 {
			start := FindTokenIndex(tokens, "iteration_start")
			end := FindTokenIndex(tokens, "iteration_end")

			startToken := token.NewToken("string", "(")
			endToken := token.NewToken("string", ")*") // FIXME: epsilon

			tokens[start] = startToken
			tokens[end] = endToken
		} else {
			text := ""
			wasString := false

			if tokens[0].ID == "string" {
				wasString = true
				text = text + `"`
			} else {
				wasString = false
			}

			for _, t := range tokens {
				if t.ID == "string" {
					if wasString {
						text = text + t.Lexema
					} else {
						wasString = true
						text = text + `+"` + t.Lexema
					}
				} else if t.ID == "var" {
					if wasString {
						wasString = false
						text = text + `"+` + t.Lexema
					} else {
						text = text + t.Lexema
					}
				}
			}
			return text + `"`
		}
	}
}

func ParseTokens(tokens []token.Token) string {
	definitions := GetDefinitions(tokens)
	tk := map[string]string{}

	for _, def := range definitions {
		params := GetDefinitionParams(def)
		id := params[0][0].Lexema

		tokenDesc := params[2]

		d := parseTokenDefinition(tokenDesc)
		tk[id] = d
	}

	text := ""
	for i, d := range tk {
		text = text + "tkns = append(tkns, token.NewTokenDescriptor(\"" + i + "\", " + d + "))\n"
	}

	return text
}
