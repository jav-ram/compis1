package Token

import "fmt"

// Token is a token data structure
type Token struct {
	ID         string
	Lexema     string
	Attributes []string
}

type TokenDescriptor struct {
	ID        string
	Rgx       string
	IsKeyword bool
}

// NewToken makes a new token
func NewToken(id string, lex string) Token {
	return Token{id, lex, nil}
}

// SetAttributes gives attributes to a token
func (token Token) SetAttributes(attributes ...string) {
	token.Attributes = attributes
}

// String custom print
func (token Token) String() string {
	return fmt.Sprintf("<'%v', '%v'>", token.ID, token.Lexema)
}

func NewTokenDescriptor(id string, rgs string) TokenDescriptor {
	return TokenDescriptor{id, rgs, false}
}

func NewKeywordTokenDescriptor(id string, rgs string) TokenDescriptor {
	return TokenDescriptor{id, rgs, true}
}
