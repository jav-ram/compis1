package main

import (
	"fmt"
	"io/ioutil"

	scanner "github.com/ram16230/compis1/Scanner"
	token "github.com/ram16230/compis1/Token"
)

func main() {

	digit := `(0|1|2|3|4|5|6|7|8|9)`
	_ = digit
	sign := `(+|-)`
	_ = sign
	hexdigit := `(0|1|2|3|4|5|6|7|8|9|A|B|C|D|E|F)`
	_ = hexdigit
	eol := `(
)`
	_ = eol
	letterLo := `(a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)`
	_ = letterLo
	letter := `(a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z|A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z)`
	_ = letter
	consonants := `(b|c|d|f|g|h|j|k|l|m|n|p|q|r|s|t|v|w|x|y|z|B|C|D|F|G|H|J|K|L|M|N|P|Q|R|S|T|V|W|X|Y|Z)`
	_ = consonants
	space := `( )`
	_ = space
	whitespace := `(
|
|	| )`
	_ = whitespace
	letterUp := `(A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z)`
	_ = letterUp
	vowels := `(a|e|i|o|u|A|E|I|O|U)`
	_ = vowels
	tab := `(	)`
	_ = tab

	var tkns []token.TokenDescriptor
	tkns = append(tkns, token.NewKeywordTokenDescriptor("if", "if"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("while", "while"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("switch", "switch"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("do", "do"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("for", "for"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("exit", "exit"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("class", "class"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("import", "import"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("from", "from"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("try", "try"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("except", "except"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("lambda", "lambda"))

	tkns = append(tkns, token.NewTokenDescriptor("var", letter+"("+letter+"|"+digit+")*"+digit+""))
	tkns = append(tkns, token.NewTokenDescriptor("signInt", "("+sign+"|`)"+digit+"("+digit+")*"))
	tkns = append(tkns, token.NewTokenDescriptor("int", digit+"("+digit+")*"))
	tkns = append(tkns, token.NewTokenDescriptor("float", digit+"("+digit+")*."+digit+"("+digit+")*"))
	tkns = append(tkns, token.NewTokenDescriptor("hexnumber", hexdigit+"("+hexdigit+")*(H)"))
	tkns = append(tkns, token.NewTokenDescriptor("space", whitespace+"("+whitespace+")*"))
	tkns = append(tkns, token.NewTokenDescriptor("string", letter+"("+letter+")*"))
	tkns = append(tkns, token.NewTokenDescriptor("name", letterUp+"("+letterLo+")("+letterLo+")*"))

	scan := scanner.MakeAFNS(tkns)
	data, _ := ioutil.ReadFile("./test/test.txt")

	tokens := scan.Simulate(string(data))
	fmt.Printf("%v\n", tokens)
}
