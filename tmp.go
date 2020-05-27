package tmp

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
	next := index
	if next < len(list) {
		switch v := item.(type) {
		case string:
			{
				return list[next].Lexema == v
			}
		case token.Token:
			{
				return v.ID == list[next].ID
			}
		}
	}
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

func Expr() bool {

	status := true
	for true {
		status = status && Stat()
		if !status {
			break
		}
	}
	status = status && Read(".")

	return status
}
func Stat() bool {
	status := true
	var value int
	status = status && Expression(&value)
	fmt.Printf("%v\n", value)

	return status
}
func Expression(result *int) bool {
	status := true
	var result1 int
	var result2 int
	status = status && Term(&result1)
	for true {
		if Expect("+") {
			status = status && Read("+")
			status = status && Term(&result2)
			result1 += result2

		} else if Expect("-") {
			status = status && Read("-")
			status = status && Term(&result2)
			result1 -= result2

		} else {
			break
		}
	}
	*result = result1

	return status
}
func Term(result *int) bool {
	status := true
	var result1, result2 int
	status = status && Factor(&result1)
	for true {
		if Expect("*") {
			status = status && Read("*")
			status = status && Factor(&result2)
			result1 *= result2
		} else if Expect("/") {
			status = status && Read("/")
			status = status && Factor(&result2)
			result1 /= result2

		} else {
			break
		}
	}
	*result = result1

	return status
}

func Factor(result *int) bool {
	status := true
	var signo int = 1
	if Expect("-") {
		status = status && Read("-")
		signo = -1
	}
	status = status && Number(result)
	r := *result * signo

	*result = r

	return status
}

func Number(result *int) bool {
	status := true
	status = status && Read(token.NewToken("number", ""))
	r, _ := strconv.Atoi(list[index-1].Lexema)
	*result = r
	return status
}

func main() {
	letter := `(A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)`
	_ = letter
	digit := `(0|1|2|3|4|5|6|7|8|9)`
	_ = digit
	tab := `(	)`
	_ = tab
	eol := `(
)`
	_ = eol

	var tkns []token.TokenDescriptor
	tkns = append(tkns, token.NewKeywordTokenDescriptor("while", "while"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("do", "do"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("if", "if"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("switch", "switch"))

	tkns = append(tkns, token.NewTokenDescriptor("number", digit+"("+digit+")*"))
	tkns = append(tkns, token.NewTokenDescriptor("ident", letter+"("+letter+"|"+digit+")*"))

	scan := scanner.MakeAFNS(tkns)
	data, _ := ioutil.ReadFile("./test/test.txt")

	tokens := scan.Simulate(string(data))
	list = tokens

	fmt.Printf("%v\n", tokens)
	Expr()
}
