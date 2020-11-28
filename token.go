package golas

// TokenType represents a lexical token type
type TokenType uint

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
}

// comment
const (
	TEndOfFile TokenType = iota
	TComment
	TVersionInformation
	TWellInformation
	TCurveInformation
	TASCIILogData
	TParameterInformation
	TOther
	TSectionCustom
	TMnemonic
	TUnits
	TData
	TDescription
)
