package main
import (
	"fmt"
	"io/ioutil"
	"strconv"

	scanner "github.com/ram16230/compis1/Scanner"
	token "github.com/ram16230/compis1/Token"
)

var list []token.Token = []token.Token{}
var index int = 0
var lastToken = token.NewToken("", "")

func Expect(item interface{}) bool {
	i := index
	next := index
	if next < len(list) {
		switch v := item.(type) {
		case string:
			{
				index = i
				return list[next].Lexema == v
			}
		case token.Token:
			{
				index = i
				return v.ID == list[next].ID
			}
		case bool:
			{
				return v
			}
		}
	}
	index = i
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
		lastToken = list[index]
		index++
	}

	return status
}

func Expr()bool {

status := true

curr := index
_ = curr
for true {
status = status && Stat()
status = status && Read(";")
if !status {
status = true 
break 
}
}
status = status && Read(".")

return status
}
func Stat()bool {

status := true

curr := index
_ = curr
var value int
status = status && Expression(&value)
 fmt.Printf("%v\n", value) 

return status
}
func Expression(result *int)bool {

status := true

curr := index
_ = curr
var result1,result2 int
status = status && Term(&result1)
for true {

curr = index
status = status && (Expect("+") || Expect("-"))
index = curr

curr = index
if Expect("+") {
index = curr
status = status && Read("+")
status = status && Term(&result2)
result1+=result2;

} else if Expect("-") {
index = curr
status = status && Read("-")
status = status && Term(&result2)
result1-=result2;

}
if !status {
status = true 
break 
}
}
*result = result1

return status
}
func Term(result *int)bool {

status := true

curr := index
_ = curr
var result1,result2 int
status = status && Factor(&result1)
for true {

curr = index
status = status && (Expect("*") || Expect("/"))
index = curr

curr = index
if Expect("*") {
index = curr
status = status && Read("*")
status = status && Factor(&result2)
result1*=result2;

} else if Expect("/") {
index = curr
status = status && Read("/")
status = status && Factor(&result2)
result1/=result2;

}
if !status {
status = true 
break 
}
}
*result = result1

return status
}
func Factor(result *int)bool {

status := true

curr := index
_ = curr
var signo int =1
curr = index
if Expect("-"){
index = curr
status = status && Read("-")
signo = -1
}

curr = index
status = status && (Expect(Number(result)) || Expect("("))
index = curr

curr = index
if Expect(Number(result)) {
index = curr
status = status && Number(result)

} else if Expect("(") {
index = curr
status = status && Read("(")
status = status && Expression(result)
status = status && Read(")")

}
*result*=signo

return status
}
func Number(result *int)bool {

status := true

curr := index
_ = curr
status = status && Read(token.NewToken("number", ""))
 r, _ := strconv.Atoi(lastToken.Lexema) 
 *result = r 

return status
}
func main() {
tab := `(	)` 
_ = tab
eol := `(
)` 
_ = eol
letter := `(A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)` 
_ = letter
digit := `(0|1|2|3|4|5|6|7|8|9)` 
_ = digit

var tkns []token.TokenDescriptor
tkns = append(tkns, token.NewKeywordTokenDescriptor("while", "while"))
tkns = append(tkns, token.NewKeywordTokenDescriptor("do", "do"))
tkns = append(tkns, token.NewKeywordTokenDescriptor("if", "if"))
tkns = append(tkns, token.NewKeywordTokenDescriptor("switch", "switch"))

tkns = append(tkns, token.NewTokenDescriptor("ident", letter+"("+letter+"|"+digit+")*"))
tkns = append(tkns, token.NewTokenDescriptor("number", digit+"("+digit+")*"))


scan := scanner.MakeAFNS(tkns)
data, _ := ioutil.ReadFile("./test/test.txt")

tokens := scan.Simulate(string(data))
list = tokens
Expr()
}
