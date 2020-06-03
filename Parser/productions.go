package Parser

import (
	"fmt"
	"unicode"

	token "github.com/ram16230/compis1/Token"
)

var tgtg int = 0

var blanks = []string{
	"\n",
	"\t",
	" ",
	"",
}

func remove(s []interface{}, i int) []interface{} {
	return append(s[:i], s[i+1:]...)
}

func ContainsInBlank(text string) bool {
	for _, b := range blanks {
		if b == text {
			return true
		}
	}
	return false
}

func GetProductionsDefinitions(tokens []token.Token) [][]token.Token {
	definitions := [][]token.Token{}
	definition := []token.Token{}
	for i, tkn := range tokens {
		definition = append(definition, tkn)
		if tkn.ID == "finish" {
			if i+1 < len(tokens) && i-1 >= 0 {
				if tokens[i+1].ID == "any" && (tokens[i-1].ID != "quote" || tokens[i+1].ID != "quote") {
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

func SimplifyString(tokens []token.Token) []token.Token {
	newTokens := []token.Token{}
	opened := false

	tmp := ""

	for _, tkn := range tokens {
		if tkn.ID == "quote" {
			if opened {
				// terminate string
				ntkn := token.NewToken("string", tmp)
				tmp = ""
				newTokens = append(newTokens, ntkn)
				opened = false
			} else {
				opened = true
			}
		} else if opened {
			tmp = tmp + tkn.Lexema
		} else if !opened {
			newTokens = append(newTokens, tkn)
		}
	}

	return newTokens
}

func SimplifyTranslation(tokens []token.Token) []token.Token {
	newTokens := []token.Token{}
	opened := false

	tmp := ""

	for _, tkn := range tokens {
		if tkn.ID == "translation_start" {
			opened = true
		} else if tkn.ID == "translation_end" {
			opened = false
			ntkn := token.NewToken("translation", tmp)
			tmp = ""
			newTokens = append(newTokens, ntkn)
		} else if opened {
			tmp = tmp + tkn.Lexema
		} else if !opened {
			newTokens = append(newTokens, tkn)
		}
	}

	return newTokens
}

func SimplifyParams(tokens []token.Token) []token.Token {
	newTokens := []token.Token{}
	opened := false

	tmp := ""

	for _, tkn := range tokens {
		if tkn.ID == "params_start" {
			newTokens = append(newTokens, tkn)
			opened = true
		} else if tkn.ID == "params_end" {
			opened = false
			ntkn := token.NewToken("param_translation", tmp)
			tmp = ""
			newTokens = append(newTokens, ntkn)
			newTokens = append(newTokens, tkn)
		} else if opened {
			tmp = tmp + tkn.Lexema
		} else if !opened {
			newTokens = append(newTokens, tkn)
		}
	}

	return newTokens
}

func CleanAnys(tokens []token.Token) []token.Token {
	newTokens := []token.Token{}

	for _, tkn := range tokens {
		if tkn.ID != "any" {
			newTokens = append(newTokens, tkn)
		}
	}

	return newTokens
}

func SimplifyDeclarations(tokens []token.Token, vars []string) []token.Token {
	newTokens := []token.Token{}
	added := false

	for _, tkn := range tokens {
		for _, v := range vars {
			if tkn.Lexema == v {
				added = true
				if unicode.IsUpper(rune(v[0])) {
					newTokens = append(newTokens, token.NewToken("production", v))
				} else {
					newTokens = append(newTokens, token.NewToken("variable", v))
				}
				break
			}
		}
		if !added {
			newTokens = append(newTokens, tkn)
		}
		added = false
	}

	return newTokens
}

func Clean(tokens []token.Token, vars []string) []token.Token {
	tokens = SimplifyDeclarations(tokens, vars)
	return CleanAnys(SimplifyString(SimplifyParams(SimplifyTranslation(tokens))))
}

func CleanIds(tokens []token.Token, vars []string) []token.Token {
	tokens = SimplifyDeclarations(tokens, vars)
	return SimplifyString(SimplifyParams(SimplifyTranslation(tokens)))
}

func CleanBreaks(tokens []token.Token) []token.Token {
	newTokens := []token.Token{}
	for _, tkn := range tokens {
		if tkn.ID != "break" {
			newTokens = append(newTokens, tkn)
		}
	}
	return newTokens
}

func FindInTokenList(tokens []token.Token, id string) int {
	for i, tkn := range tokens {
		if tkn.ID == id {
			return i
		}
	}
	return -1
}

func GetList(tokens []token.Token, tknList []string) []interface{} {
	var list []interface{}
	count := 0

	var temp []token.Token

	for i, tkn := range tokens {
		if tkn.ID == "ident" {
			isIn := false

			for _, s := range tknList {
				if s == tkn.Lexema {
					isIn = true
					break
				}
			}

			if isIn {
				tokens[i].ID = "token"
			}
		}
	}

	for _, tkn := range tokens {
		if tkn.ID == "group_start" {
			count++
			continue
		} else if tkn.ID == "group_end" {
			count--
			if count == 0 {
				temp = append(temp, tkn)
				t := GetList(temp, tknList)
				list = append(list, t)
				temp = []token.Token{}
			}
			continue
		}
		if count == 0 {
			list = append(list, tkn)
		} else {
			temp = append(temp, tkn)
		}
	}

	/* fmt.Println("-------")

	fmt.Printf("i: %v \n%v\n", tgtg, list)
	fmt.Printf("i: %v \n%v\n", tgtg, GetOrs(list))
	tgtg++
	fmt.Println("-------") */

	return list
}

func GetOrs(list []interface{}) []interface{} {

	var newList []interface{}
	for i, item := range list {
		switch v := item.(type) {
		case token.Token:
			{
				if v.ID == "or" {
					var newList []interface{}
					// list left
					// check if is a list
					var left []interface{}
					switch k := list[i-1].(type) {
					case []token.Token:
						{
							for _, t := range k {
								left = append(left, t)
							}
						}
					case token.Token:
						{
							start := 0
							end := i - 1
							if k.ID == "break" {
								// find start
								for j := end - 1; j >= start; j-- {

									switch x := list[j].(type) {
									case token.Token:
										{
											if x.ID == "break" {
												for k := j; k < len(list); k++ {
													if list[k].(token.Token).ID == "string" || list[k].(token.Token).ID == "production" {
														start = k
														break
													}
												}
											}
										}
									}

								}

								for j := 0; j <= end; j++ {
									if j >= start {
										left = append(left, list[j])

										if j == end-1 {
											newList = append(newList, left)
											break
										}
									} else {
										newList = append(newList, list[j])
									}
								}

							}
						}
					}
					newList = append(newList, list[i])
					// list right
					// check if is a list
					var right []interface{}
					switch k := list[i+1].(type) {
					case []token.Token:
						{
							{
								for _, t := range k {
									right = append(left, t)
								}
							}
						}
					case token.Token:
						{
							start := i + 1
							end := len(list)
							// find start
							for j := start; j < end; j++ {

								switch x := list[j].(type) {
								case token.Token:
									{
										if x.ID == "break" {
											end = j
											break
										}
									}
								}

							}

							for j := start; j < len(list); j++ {
								if j <= end {
									right = append(right, list[j])
									if j == end-1 {
										newList = append(newList, right)
									}
								} else {
									newList = append(newList, list[j])
								}
							}

						}
					}
					return newList
				} else {
					// added to newList
					newList = append(newList, v)
				}
			}
		case []interface{}:
			{
				t := GetOrs(v)
				newList = append(newList, t)
			}
		}
	}

	return newList
}

func GetFirst(list []interface{}) []string {
	set := []string{}
	for i, l := range list {
		switch v := l.(type) {
		case []interface{}:
			{
				set = append(set, GetFirst(v)...)
			}
		case token.Token:
			{
				if v.ID == "production" {
					params := "()"
					if i+1 < len(list) && list[i+1].(token.Token).ID == "params_start" {
						params = "(" + list[i+2].(token.Token).Lexema + ")"
					}
					return append(set, v.Lexema+params)
				} else if v.ID == "string" || v.ID == "ident" {
					return append(set, `"`+v.Lexema+`"`)
				} else if v.ID == "token" {
					return append(set, `token.NewToken("`+v.Lexema+`", "")`)
				}
			}
		}
	}
	return set
}

func GetTokenTranslation(tkn token.Token, isID bool) string {
	linebreak := "\n"
	switch tkn.ID {
	case "iteration_start":
		{
			return "for status {" + linebreak
		}
	case "iteration_end":
		{
			return "if !status {\nstatus = true \nbreak \n}\n}" + linebreak
		}
	case "translation":
		{
			if isID {
				return tkn.Lexema
			} else {
				return tkn.Lexema + linebreak
			}
		}
	case "param_translation":
		{
			return tkn.Lexema
		}
	case "params_start":
		{
			return "("
		}
	case "params_end":
		{
			if isID {
				return ")"
			} else {
				return ")" + linebreak
			}
		}
	case "production":
		{
			return tkn.Lexema
		}
	case "string":
		{
			return `status = status && Read("` + tkn.Lexema + `")` + linebreak
		}
	case "ident":
		{
			return tkn.Lexema + " "
		}
	case "option_end":
		{
			return "}" + linebreak
		}
	case "token":
		{
			return `status = status && Read(token.NewToken("` + tkn.Lexema + `", ""))` + linebreak
		}
	}
	return ""
}

func Translate(list []interface{}, isID bool) string {
	newString := ""
	for i := 0; i < len(list); i++ {
		item := list[i]
		switch v := item.(type) {
		case []interface{}:
			{
				// check if next is a OR
				firstOr := false
				if i+1 < len(list) {
					switch k := list[i+1].(type) {
					case token.Token:
						{
							firstOr = k.ID == "or"
						}
					}
				}
				if firstOr {
					fmt.Printf("%v\n", v)
					f := GetFirst(v)[0]                         // get first of 1
					l := GetFirst(list[i+2].([]interface{}))[0] // get first of 2
					newString = newString + `
curr = index
status = status && (Expect(` + f + `) || Expect(` + l + `))
index = curr

curr = index
if Expect(` + f + `) {
index = curr
` + Translate(v, isID) + `
} else if Expect(` + l + `) {
index = curr
` + Translate((list[i+2]).([]interface{}), isID) + `
}
`
					i += 2
					continue
				} else {
					// list call again Translate
					newString = newString + Translate(v, isID)
				}
			}
		case token.Token:
			{
				if v.ID == "production" {
					if i+1 < len(list) {
						switch k := list[i+1].(type) {
						case []interface{}:
							{
								if isID {
									newString = newString + GetTokenTranslation(v, isID) + "()"
								} else {
									newString = newString + "status = status && " + GetTokenTranslation(v, isID) + "()\n"
								}
								continue
							}
						case token.Token:
							{
								if k.ID != "params_start" {
									if isID {
										newString = newString + GetTokenTranslation(v, isID) + "()"
									} else {
										newString = newString + "status = status && " + GetTokenTranslation(v, isID) + "()\n"
									}
									continue
								} else if !isID {
									newString = newString + "status = status && " + GetTokenTranslation(v, isID)
									continue
								}
							}
						}
					} else {
						if isID {
							newString = newString + GetTokenTranslation(v, isID) + "()"
						} else {
							newString = newString + "status = status && " + GetTokenTranslation(v, isID) + "()\n"
						}
						continue
					}
				}

				if v.ID == "option_start" {
					first := ""
					for j := i + 1; j < len(list); j++ {
						switch k := list[j].(type) {
						case []interface{}:
							{
								first = GetFirst(k)[0]
								break
							}
						case token.Token:
							{
								if k.ID == "string" || k.ID == "ident" {
									first = k.Lexema
									newString = newString + "curr = index\n" + `if Expect("` + first + `")` + "{\nindex = curr\n"
									j = len(list)
								} else if k.ID == "production" {
									params := "()"
									if j+1 < len(list) {
										if list[j+1].(token.Token).ID == "params_start" {
											params = "(" + list[j+2].(token.Token).Lexema + ")"
										}
									}
									newString = newString + "curr = index\n" + `if Expect(` + first + params + `)` + "{\nindex = curr\n"
									j = len(list)
								}
							}
						}
					}
					continue
				}

				newString = newString + GetTokenTranslation(v, isID)
			}

		}
	}
	return newString
}

func ParseProductions(tokens []token.Token, tokenList []string) (string, string) {
	definitions := GetProductionsDefinitions(tokens)
	productions := []string{}
	code := ""

	for _, def := range definitions {
		params := GetDefinitionParams(def)
		i := GetDefinitionName(params[0])
		productions = append(productions, i)
	}

	for i, def := range definitions {
		params := GetDefinitionParams(def)
		//fmt.Printf("list: %v\n", GetList(CleanBreaks(Clean(params[0], variables))))
		id := Translate(GetList(CleanBreaks(CleanIds(params[0], productions)), tokenList), true)

		definition := Clean(params[2], productions)

		if i == len(definitions)-3 {
			fmt.Println("-------")
			fmt.Printf("%v\n", GetOrs(GetList(definition, tokenList)))
			fmt.Println("-------")
		}

		code += "func " + id + `bool {

status := true

curr := index
_ = curr
` + Translate(GetOrs(GetList(definition, tokenList)), false) + "\nreturn status\n}\n"
	}

	return code, productions[0] + "()"
}
