package golas

import (
	"fmt"
)

// handleEntry is a state function
func handleEntry(lexer *Lexer) HandlerFunc {
	lexer.step()
	switch lexer.char {
	case '~':
		return handleSection
	case '#':
		return handleComment
	default:
		panic("unknown first character")
	}
}

// handleSection lexes a section
func handleSection(lexer *Lexer) HandlerFunc {
	if lexer.position != 1 {
		panic(fmt.Errorf("invalid las file section : tilde not first character on line : line %d : position %d", lexer.line+1, lexer.position))
	}

	lexer.step()

	switch lexer.char {
	case 'A':
		lexer.validateSection("a")
		lexer.stepUntil('\n')
		lexer.buffer.Reset()
		lexer.emit(TSectionLogs)
		return handleLogs
	case 'V':
		lexer.validateSection("v")
		lexer.overwriteBuffer("Version Information")
	case 'W':
		lexer.validateSection("w")
		lexer.overwriteBuffer("Well Information")
	case 'C':
		lexer.validateSection("c")
		lexer.overwriteBuffer("Curve Information")
	case 'P':
		lexer.overwriteBuffer("Parameter Information")
	case 'O':
		lexer.overwriteBuffer("Other Information")
	default:
		// TODO : prob a better way to do this for custom sections
		lexer.stepUntil('\n')
	}

	lexer.emit(TSection)
	lexer.stepUntil('\n')
	return getNextHandler(lexer)
}

// handleComment lexes a comment within a line
func handleComment(lexer *Lexer) HandlerFunc {
	lexer.stepUntil('\n')
	lexer.emit(TComment)
	return getNextHandler(lexer)
}

// handleMnemonic lexes a mnemonic within a non-ascii log data line
func handleMnemonic(lexer *Lexer) HandlerFunc {
	lexer.stepUntil('.')
	lexer.truncate()
	lexer.emit(TMnemonic)
	return handleUnits
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
	lexer.emit(TDescription)
	return getNextHandler(lexer)
}

func getNextHandler(lexer *Lexer) HandlerFunc {
	if lexer.char != '\n' {
		// Make sure we are at the end of the line
		lexer.stepUntil('\n', -1)
	}

	switch lexer.peekNext() {
	case -1:
		return nil
	case '#':
		lexer.step()
		return handleComment
	case '~':
		lexer.step()
		return handleSection
	default:
		return handleMnemonic
	}
}
