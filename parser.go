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
	lexer.Start(handleEntry)

	for {
		token = lexer.NextToken()
		token.Value = strings.TrimSpace(token.Value)

		if token.Type == TEndOfFile {
			appendLine(section, line)
			appendSection(las, section)
			break
		}

		switch token.Type {
		case TSection:
			appendSection(las, section)
			section = &Section{Name: token.Value}

		case TSectionLogs:
			appendSection(las, section)
			las.Logs = append(las.Logs, strings.Fields(strings.TrimSpace(token.Value)))

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

	return las
}

func appendSection(las *LAS, section *Section) {
	if section != nil {
		las.Sections = append(las.Sections, *section)
		section = nil
	}
}

func appendLine(section *Section, line *Line) {
	if line != nil {
		section.Lines = append(section.Lines, *line)
	}
}
