package golas

import (
	"strings"
)

// LAS represents a .las file
type LAS struct {
	Sections  []Section
	ASCIILogs LogData
	version   string
	wrap      string
}

// IsWrapped returns whether or not the las file is wrapped
func (las *LAS) IsWrapped() bool {
	if las.wrap != "" {
		return strings.ToLower(las.wrap) == "yes"
	}
	for s := range las.Sections {
		if strings.ToLower(las.Sections[s].Name) == "version information" {
			for l := range las.Sections[s].Lines {
				if strings.ToLower(las.Sections[s].Lines[l].Mnem) == "wrap" {
					las.wrap = las.Sections[s].Lines[l].Data
					break
				}
			}
		}
		if las.wrap != "" {
			break
		}
	}
	return strings.ToLower(las.wrap) == "yes"
}

// Version returns the las file version
func (las *LAS) Version() string {
	if las.version != "" {
		return las.version
	}
	for s := range las.Sections {
		if strings.ToLower(las.Sections[s].Name) == "version information" {
			for l := range las.Sections[s].Lines {
				if strings.ToLower(las.Sections[s].Lines[l].Mnem) == "vers" {
					las.version = las.Sections[s].Lines[l].Data
					break
				}
			}
		}
		if las.version != "" {
			break
		}
	}
	return las.version
}

// Line represents a header line in a .las file section
type Line struct {
	Mnem        string
	Units       string
	Data        string
	Description string
}

// LogData represents a row in the ASCII Log Data section ('~A')
type LogData struct {
	Headers []string
	Rows    [][]string
}

// Section represents a .las file section
type Section struct {
	Name     string
	Lines    []Line
	Comments []string
}
