package golas

import (
	"fmt"
)

// handleBegin is a state function
func handleBegin(lexer *Lexer) HandlerFunc {
	switch lexer.char {
	case '~':
		return handleSection
	case '#':
		return handleComment
	case '.':
		return handleMnemonic
	default:
		lexer.step()
		return handleBegin
	}
}

// handleSection lexes a section
func handleSection(lexer *Lexer) HandlerFunc {
	if lexer.position != 1 {
		panic(fmt.Errorf("invalid las file section : tilde not first character on line : line %d : position %d", lexer.line+1, lexer.position))
	}
	var s string
	var t = TSection
	var h = handleMnemonic
	lexer.step()

	switch lexer.char {
	case 'A':
		// When dealing with ASCII Log section we don't care about anything after the ~A on the same line
		lexer.stepUntil('\n')
		lexer.buffer.Reset()
		t = TSectionASCIILogs
		h = handleASCIILogs
	case 'V':
		s = "Version Information"
	case 'W':
		s = "Well Information"
	case 'C':
		s = "Curve Information"
	case 'P':
		s = "Parameter Information"
	case 'O':
		s = "Other Information"
	default:
		t = TSectionCustom
	}

	lexer.stepUntil('\n')
	if t == TSection {
		// Only overwrite buffer is using a reserved section (not a custom section or ASCII log section)
		lexer.overwriteBuffer(s)
	}
	lexer.emit(t)
	return h
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
	// Mnemonic ends at the first dot (period) on a line
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
	return handleData
}

// handleLineData lexes data within a non-ascii log data line
func handleData(lexer *Lexer) HandlerFunc {
	for lexer.char != ':' {
		lexer.step()
	}
	lexer.truncate()
	lexer.emit(TData)
	return handleDescription
}

func handleASCIILogs(lexer *Lexer) HandlerFunc {
	// Step once more to see if we should continue reading ASCII log data or not
	lexer.step()
	switch lexer.char {
	case '~':
		return handleSection
	case '#':
		return handleComment
	case -1:
		lexer.emit(TEndOfFile)
		return nil
	default:
		// Step until new line as this will fill our buffer
		// with the value of said line. After emitting the line as Token.Value
		// we clear the buffer to save on resources.
		lexer.stepUntil('\n')
		lexer.emit(TSectionASCIILogs)
		return handleASCIILogs
	}
}

// handleDescription lexes a description within a non-ascii log data line
func handleDescription(lexer *Lexer) HandlerFunc {
	for lexer.char != '\n' {
		if lexer.char == -1 {
			lexer.emit(TDescription)
			return nil
		}
		lexer.step()
	}
	lexer.emit(TDescription)
	return handleBegin
}
