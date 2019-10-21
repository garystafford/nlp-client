package main

import (
	"gopkg.in/jdkato/prose.v2"
)

// A Token represents an individual token of text such as a word or punctuation symbol.
// IOB format (short for inside, outside, beginning) is a common tagging format
type Token struct {
	Tag   string `json:"tag"`   // The token's part-of-speech tag.
	Text  string `json:"text"`  // The token's actual content.
	Label string `json:"label"` // The token's IOB label.
}

func getTokens(phrase string) (interface{}, error) {
	doc, err := prose.NewDocument(phrase)
	if err != nil {
		return nil, err
	}

	var tokens []Token

	for _, token := range doc.Tokens() {
		tokens = append(tokens, Token{Tag: token.Tag, Text: token.Text, Label: token.Label})
	}

	return tokens, nil
}
