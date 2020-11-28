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
	lexer := NewLexer(r)

	lexer.Start(HandleBegin)

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
			las.Sections = append(las.Sections, *section)
			return las
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

// HandleBegin is a state function
func HandleBegin(lexer *Lexer) HandlerFunc {
	if lexer.char == '~' {
		return HandleSection
	} else if lexer.char == '#' {
		return HandleComment
	} else if lexer.char == '.' {
		return HandleMnemonic
	} else {
		lexer.step()
		return HandleBegin
	}
}

// HandleSection lexes a section
func HandleSection(lexer *Lexer) HandlerFunc {
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
		s = "ASCII Log Data"
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
	if t != TSectionCustom {
		lexer.overwriteBuffer(s)
	}
	lexer.emit(t)
	return HandleMnemonic
}

// HandleComment lexes a comment within a line
func HandleComment(lexer *Lexer) HandlerFunc {
	for lexer.char != '\n' {
		lexer.step()
	}
	lexer.emit(TComment)
	return HandleBegin
}

// HandleMnemonic lexes a mnemonic within a non-ascii log data line
func HandleMnemonic(lexer *Lexer) HandlerFunc {
	// Mnemonic only valid if it is the first dot on a line
	if lexer.dots == 1 {
		lexer.truncate()
		lexer.emit(TMnemonic)
		return HandleUnits
	}
	return HandleBegin
}

// HandleUnits lexes units within a non-ascii log data line
func HandleUnits(lexer *Lexer) HandlerFunc {
	for lexer.char != ' ' {
		lexer.step()
	}
	lexer.truncate()
	lexer.emit(TUnits)
	return HandleLineData
}

// HandleLineData lexes data within a non-ascii log data line
func HandleLineData(lexer *Lexer) HandlerFunc {
	for lexer.char != ':' {
		lexer.step()
	}
	lexer.truncate()
	lexer.emit(TData)
	return HandleDescription
}

// HandleDescription lexes a description within a non-ascii log data line
func HandleDescription(lexer *Lexer) HandlerFunc {
	for lexer.char != '\n' {
		lexer.step()
	}
	lexer.emit(TDescription)
	return HandleBegin
}
