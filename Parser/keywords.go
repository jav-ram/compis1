package Parser

import (
	token "github.com/ram16230/compis1/Token"
)

func ParseKeywords(tokens []token.Token) string {
	definitions := GetDefinitions(tokens)
	idents := []string{}
	descriptions := []string{}

	text := "var tkns []token.TokenDescriptor\n"

	for _, def := range definitions {
		params := GetDefinitionParams(def)
		ident := GetDefinitionName(params[0])
		desc := GetDefinitionName(params[2])
		idents = append(idents, ident)
		descriptions = append(descriptions, desc)

		text = text + "tkns = append(tkns, token.NewKeywordTokenDescriptor(\"" + ident + "\", \"" + desc + "\"))\n"
	}

	return text
}
