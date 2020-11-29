package golas

import "fmt"

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
