package main

import (
	"fmt"
	"io/ioutil"

	scanner "github.com/ram16230/compis1/Scanner"
	token "github.com/ram16230/compis1/Token"
)

func main() {

	hexterm := `(H)`
	_ = hexterm
	tab := `(	)`
	_ = tab
	eol := `(
)`
	_ = eol
	whitespace := `(
|
|	|
)`
	_ = whitespace
	sign := `(+|-)`
	_ = sign
	upletter := `(A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z)`
	_ = upletter
	digit := `(0|1|2|3|4|5|6|7|8|9)`
	_ = digit
	hexdigit := `(0|1|2|3|4|5|6|7|8|9|A|B|C|D|E|F)`
	_ = hexdigit
	downletter := `(a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)`
	_ = downletter
	letter := `(a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z|A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)`
	_ = letter

	var tkns []token.TokenDescriptor
	tkns = append(tkns, token.NewKeywordTokenDescriptor("while", "while"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("do", "do"))

	tkns = append(tkns, token.NewTokenDescriptor("hexnumber", hexdigit+"("+hexdigit+")*"+hexterm+""))
	tkns = append(tkns, token.NewTokenDescriptor("ident", letter+"("+letter+"|"+digit+")*"))
	tkns = append(tkns, token.NewTokenDescriptor("number", digit+"("+digit+")*"))
	tkns = append(tkns, token.NewTokenDescriptor("signnumber", "("+sign+"|`)"+digit+"("+digit+")*"))

	scan := scanner.MakeAFNS(tkns)
	data, _ := ioutil.ReadFile("./test/test.txt")

	tokens := scan.Simulate(string(data))
	fmt.Printf("%v\n", tokens)
}
