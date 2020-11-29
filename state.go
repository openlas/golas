package golas

import "fmt"

// handleBegin is a state function
func handleBegin(lexer *Lexer) HandlerFunc {
	if lexer.char == '~' {
		return handleSection
	} else if lexer.char == '#' {
		return handleComment
	} else if lexer.char == '\n' {
		return handleLine
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

	var s string
	var t = TSection
	var h = handleMnemonic
	lexer.step()

	switch lexer.char {
	case 'V':
		s = "Version Information"
	case 'W':
		s = "Well Information"
	case 'C':
		s = "Curve Information"
	case 'A':
		t = TSectionASCIILogs
		h = handleASCIILogs
	case 'P':
		s = "Parameter Information"
	case 'O':
		s = "Other Information"
	}

	lexer.stepUntil('\n')
	// If a regular header section (non ascii log data)
	if t == TSection {
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
	for {
		lexer.step()
		if lexer.char == '~' {
			return handleSection
		}
		if lexer.char == -1 {
			lexer.emit(TEndOfFile)
			return nil
		}
		// code to split row string (lexer.buffer) into string slice
	}
}

func handleLine(lexer *Lexer) HandlerFunc {
	for {
		lexer.step()
		switch lexer.char {
		case '.':
			lexer.emit(TMnemonic)
			return handleUnits
		case '#':
			lexer.emit(TComment)
			return handleComment
		}
	}
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
