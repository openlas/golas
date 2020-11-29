package golas

import (
	"fmt"
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
		case TVersionInformation, TWellInformation, TCurveInformation, TParameterInformation, TOther, TSectionCustom:
			if section != nil {
				las.Sections = append(las.Sections, *section)
			}
			section = &Section{Name: token.Value}
		case TASCIILogData:
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

// handleBegin is a state function
func handleBegin(lexer *Lexer) HandlerFunc {
	if lexer.char == '~' {
		return handleSection
	} else if lexer.char == '#' {
		return handleComment
	} else if lexer.char == '.' {
		return handleMnemonic
	} else {
		lexer.step()
		return handleBegin
	}
}

// handleSection lexes a section
func handleSection(lexer *Lexer) HandlerFunc {
	if lexer.position != 1 {
		panic(fmt.Errorf("invalid las file section : tilde not first character on line : line %d : position %d", lexer.line+1, lexer.position))
	}

	var t TokenType
	var s string
	lexer.step()

	switch lexer.char {
	case 'V':
		s = "Version Information"
		t = TVersionInformation
	case 'W':
		s = "Well Information"
		t = TWellInformation
	case 'C':
		s = "Curve Information"
		t = TCurveInformation
	case 'A':
		s = "ASCII Logs"
		t = TASCIILogData
	case 'P':
		s = "Parameter Information"
		t = TParameterInformation
	case 'O':
		s = "Other Information"
		t = TOther
	default:
		t = TSectionCustom
	}

	lexer.stepUntil('\n')
	// If not custom section, use hard coded string as name
	if t != TSectionCustom && s != "" {
		lexer.overwriteBuffer(s)
	}
	lexer.emit(t)
	return handleMnemonic
}

// handleComment lexes a comment within a line
func handleComment(lexer *Lexer) HandlerFunc {
	for lexer.char != '\n' {
		lexer.step()
	}
	lexer.emit(TComment)
	return handleBegin
}

// handleMnemonic lexes a mnemonic within a non-ascii log data line
func handleMnemonic(lexer *Lexer) HandlerFunc {
	// Mnemonic only valid if it is the first dot on a line
	if lexer.dots == 1 {
		lexer.truncate()
		lexer.emit(TMnemonic)
		return handleUnits
	}
	return handleBegin
}

// handleUnits lexes units within a non-ascii log data line
func handleUnits(lexer *Lexer) HandlerFunc {
	for lexer.char != ' ' {
		lexer.step()
	}
	lexer.truncate()
	lexer.emit(TUnits)
	return handleLineData
}

// handleLineData lexes data within a non-ascii log data line
func handleLineData(lexer *Lexer) HandlerFunc {
	for lexer.char != ':' {
		lexer.step()
	}
	lexer.truncate()
	lexer.emit(TData)
	return handleDescription
}

// handleDescription lexes a description within a non-ascii log data line
func handleDescription(lexer *Lexer) HandlerFunc {
	for lexer.char != '\n' {
		lexer.step()
	}
	lexer.emit(TDescription)
	return handleBegin
}

//func handleASCIILogData(lexer *Lexer) HandlerFunc {
//
//}
