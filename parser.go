package golas

import (
	"io"
	"strings"
)

// Parse parses a las file
func Parse(r io.Reader) *LAS {
	var section *Section
	var line *Line
	var token Token

	las := &LAS{}
	lexer := NewLexer(r)
	lexer.Start(handleBegin)

	for {
		token = lexer.NextToken()
		token.Value = strings.TrimSpace(token.Value)

		switch token.Type {
		case TEndOfFile:
			if line != nil {
				section.Lines = append(section.Lines, *line)
			}
			if section != nil {
				las.Sections = append(las.Sections, *section)
			}
			return las
		case TSection, TSectionCustom, TSectionLogs:
			if section != nil {
				las.Sections = append(las.Sections, *section)
				section = nil
			}
			if token.Type == TSectionLogs {
				las.Logs = append(las.Logs, strings.Fields(strings.TrimSpace(token.Value)))
			} else {
				section = &Section{Name: token.Value}
			}
		case TMnemonic:
			line = &Line{Mnem: token.Value}
		case TUnits:
			line.Units = token.Value
		case TData:
			line.Data = token.Value
		case TDescription:
			line.Description = token.Value
			section.Lines = append(section.Lines, *line)
			line = nil
		case TComment:
			section.Comments = append(section.Comments, token.Value)
		}
	}
}
