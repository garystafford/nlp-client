package main

import (
	"gopkg.in/jdkato/prose.v2"
)

// A Token represents an individual Token of Text such as a word or punctuation symbol.
// IOB format (short for inside, outside, beginning) is a common tagging format
type Token struct {
	Tag   string `json:"Tag"`   // The Token's part-of-speech Tag.
	Text  string `json:"Text"`  // The Token's actual content.
	Label string `json:"Label"` // The Token's IOB Label.
}

// An Entity represents an individual named-entity.
type Entity struct {
	Text  string // The entity's actual content.
	Label string // The entity's label.
}

func getTokens(phrase string) (interface{}, error) {
	doc, err := prose.NewDocument(phrase)
	if err != nil {
		return nil, err
	}

	var tokens []Token
	for _, docToken := range doc.Tokens() {
		tokens = append(tokens, Token{
			Tag:   docToken.Tag,
			Text:  docToken.Text,
			Label: docToken.Label,
		})
	}

	return tokens, nil
}

func getEntities(phrase string) (interface{}, error) {
	doc, err := prose.NewDocument(phrase)
	if err != nil {
		return nil, err
	}

	var entities []Entity
	for _, docEntities := range doc.Entities() {
		entities = append(entities, Entity{
			Text:  docEntities.Text,
			Label: docEntities.Label,
		})
	}

	return entities, nil
}
