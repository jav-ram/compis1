package main

import (
	"fmt"

	scanner "github.com/ram16230/compis1/Scanner"
	token "github.com/ram16230/compis1/Token"
)

func main() {

	letter := `a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z`
	_ = letter
	digit := `(0|1|2|3|4|5|6|7|8|9)`
	_ = digit
	tab := `	`
	_ = tab
	eol := `
`
	_ = eol
	blanco := `
|
|	`
	_ = blanco

	var tkns []token.TokenDescriptor
	tkns = append(tkns, token.NewKeywordTokenDescriptor("while", "while"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("do", "do"))

	tkns = append(tkns, token.NewTokenDescriptor("number", digit+"("+digit+")*"))
	tkns = append(tkns, token.NewTokenDescriptor("decnumber", digit+"("+digit+")*."+digit+"("+digit+")*"))
	tkns = append(tkns, token.NewTokenDescriptor("white", blanco+"("+blanco+")*"))

	scan := scanner.MakeAFNS(tkns)

	tokens := scan.Simulate("1.12151")
	fmt.Printf("%v\n", tokens)
}
