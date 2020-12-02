package golas

import (
	"fmt"
)

// handleEntry is used as the start of our HandlerFunc chain
func handleEntry(lexer *Lexer) HandlerFunc {
	lexer.step()
	if lexer.char == '~' {
		return handleSection
	} else if lexer.char == '#' {
		return handleComment
	}
	panic("unknown first character : " + string(lexer.char))
}

// HandleNext ensures we are at the start of a line, then figures out which handler to run from there.
func handleNext(lexer *Lexer) HandlerFunc {
	if lexer.char != '\n' {
		// Make sure we are at the end of the line before continuing
		lexer.stepUntil('\n', -1)
	}
	lexer.step()
	switch lexer.char {
	case -1:
		return nil
	case '#':
		return handleComment
	case '~':
		return handleSection
	default:
		return handleMnemonic
	}
}

// handleSection figures out what type of section we have and which handler to run
func handleSection(lexer *Lexer) HandlerFunc {
	if lexer.position != 1 {
		panic(fmt.Errorf("invalid section : char '~' not first : line %d : position %d", lexer.line+1, lexer.position))
	}
	lexer.step()
	switch lexer.char {
	case 'V', 'W', 'C', 'A':
		return handleRequiredSection
	case 'P', 'O':
		return handleOptionalSection
	default:
		return handleCustomSection
	}
}

// handleLogsSection initiates lexing ASCII logs.
func handleLogsSection(lexer *Lexer) HandlerFunc {
	if lexer.char == 'A' {
		lexer.stepUntil('\n')
		// We don't care about line 0 (the line with ~A)
		lexer.buffer.Reset()
		lexer.emit(TSectionLogs)
		return handleLogsDataLine
	}
	panic(fmt.Errorf("expected 'A' : got '%s'", string(lexer.char)))
}

// handleCustomSection lexes any custom sections a las file may contain.
// The LAS 2.0 standard allows for custom sections, so long they do not start with a reserved char.
func handleCustomSection(lexer *Lexer) HandlerFunc {
	lexer.stepUntil('\n')
	return handleNext
}

// handleOptionalSection lexes any optional sections a las file may contain
// Since the LAS 2.0 standard does not require these sections, but does reserve the characters,
// we don't need to make sure our las file contains them.
func handleOptionalSection(lexer *Lexer) HandlerFunc {
	var name string
	switch lexer.char {
	case 'P':
		name = "Parameter Information"
	case 'O':
		name = "Optional Information"
	default:
		panic("unrecognized optinal section")
	}
	lexer.stepUntil('\n')
	lexer.overwriteBuffer(name)
	lexer.emit(TSection)
	return handleNext
}

// handleRequiredSection lexes required sections as defined by the LAS standard.
func handleRequiredSection(lexer *Lexer) HandlerFunc {
	var name string
	switch lexer.char {
	case 'A':
		return handleLogsSection
	case 'V':
		name = "Version Information"
	case 'W':
		name = "Well Information"
	case 'C':
		name = "Curve Information"
	default:
		panic("unrecognized required section")
	}
	lexer.stepUntil('\n')
	lexer.overwriteBuffer(name)
	lexer.emit(TSection)
	return handleNext
}

// handleComment lexes comments.
func handleComment(lexer *Lexer) HandlerFunc {
	lexer.stepUntil('\n')
	lexer.emit(TComment)
	return handleNext
}

// handleMnemonic lexes the mnemonic within a non-ascii log data line
func handleMnemonic(lexer *Lexer) HandlerFunc {
	lexer.stepUntil('.')
	lexer.truncate()
	lexer.emit(TMnemonic)
	return handleUnits
}

// handleUnits lexes the units (of measurement) within a non-ascii log data line.
func handleUnits(lexer *Lexer) HandlerFunc {
	lexer.stepUntil(' ')
	lexer.truncate()
	lexer.emit(TUnits)
	return handleData
}

// handleLineData lexes the data within a non-ascii log data line
func handleData(lexer *Lexer) HandlerFunc {
	lexer.stepUntil(':')
	lexer.truncate()
	lexer.emit(TData)
	return handleDescription
}

// handleLogsDataLine lexes an ASCII log data line.
func handleLogsDataLine(lexer *Lexer) HandlerFunc {
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
		return handleLogsDataLine
	}
}

// handleDescription lexes a description within a non-ascii log data line
func handleDescription(lexer *Lexer) HandlerFunc {
	lexer.stepUntil('\n', -1)
	lexer.emit(TDescription)
	return handleNext
}
