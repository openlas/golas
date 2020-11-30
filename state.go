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
		lexer.validateSection("a")
		lexer.stepUntil('\n')
		lexer.buffer.Reset()
		t = TSectionLogs
		h = handleLogs
	case 'V':
		lexer.validateSection("v")
		s = "Version Information"
	case 'W':
		lexer.validateSection("w")
		s = "Well Information"
	case 'C':
		lexer.validateSection("c")
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
		lexer.overwriteBuffer(s)
	}
	lexer.emit(t)
	return h
}

// handleComment lexes a comment within a line
func handleComment(lexer *Lexer) HandlerFunc {
	lexer.stepUntil('\n')
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
	lexer.stepUntil(' ')
	lexer.truncate()
	lexer.emit(TUnits)
	return handleData
}

// handleLineData lexes data within a non-ascii log data line
func handleData(lexer *Lexer) HandlerFunc {
	lexer.stepUntil(':')
	lexer.truncate()
	lexer.emit(TData)
	return handleDescription
}

// handleLogs lexes ASCII logs section
func handleLogs(lexer *Lexer) HandlerFunc {
	lexer.step()
	switch lexer.char {
	case '~', '#':
		panic(fmt.Errorf("invalid las file : ascii logs must be the last section : found data after logs : line %d", lexer.line+1))
	case -1:
		lexer.emit(TEndOfFile)
		return nil
	default:
		lexer.stepUntil('\n', -1)
		if lexer.char == -1 {
			lexer.truncate()
		}
		lexer.emit(TSectionLogs)
		return handleLogs
	}
}

// handleDescription lexes a description within a non-ascii log data line
func handleDescription(lexer *Lexer) HandlerFunc {
	lexer.stepUntil('\n', -1)
	if lexer.char == -1 {
		lexer.emit(TDescription)
		return nil
	}
	lexer.emit(TDescription)
	return handleBegin
}
