package golas

import (
	"io"
	"strings"
)

// Parse parses a las file
func Parse(r io.Reader) *LAS {
	var section *Section
	var line Line
	var token Token
	las := &LAS{}
	// logs := Logs{}
	lexer := NewLexer(r)

	lexer.Start(handleBegin)

	for {
		token = lexer.NextToken()
		token.Value = strings.TrimSpace(token.Value)
		if token.Type == TEndOfFile {
			break
		}

		switch token.Type {
		case TSection:
			if section != nil {
				las.Sections = append(las.Sections, *section)
			}
			section = &Section{Name: token.Value}
		case TSectionASCIILogs:
			lexASCIILogs(&lexer, las)
		case TMnemonic:
			line = Line{Mnem: token.Value}
		case TUnits:
			line.Units = token.Value
		case TData:
			line.Data = token.Value
		case TDescription:
			line.Description = token.Value
			section.Lines = append(section.Lines, line)
		case TComment:
			section.Comments = append(section.Comments, token.Value)
		}
	}

	return las
}

func lexASCIILogs(lexer *Lexer, las *LAS) {
	for {
		lexer.stepUntil('\n', -1, '~')
		if lexer.char != '\n' {
			break
		}

		las.ASCIILogs.Rows = append(
			las.ASCIILogs.Rows,
			// strings.Fields will not add empty strings to slice during Split
			// (it natively splits a string using a space as separater)
			strings.Fields(strings.TrimSpace(lexer.buffer.String())),
		)
	}
}
