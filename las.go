package golas

import (
	"strings"
)

// LAS represents a .las file
type LAS struct {
	Sections  []Section
	ASCIILogs LogData
}

// IsWrapped returns whether or not the las file is wrapped
func (las *LAS) IsWrapped() bool {
	var wrapped bool
	for sectionIndex := range las.Sections {
		if strings.ToLower(las.Sections[sectionIndex].Name) == "version information" {
			for lineIndex := range las.Sections[sectionIndex].Lines {
				if strings.ToLower(las.Sections[sectionIndex].Lines[lineIndex].Mnem) == "wrap" {
					wrapped = strings.ToLower(las.Sections[sectionIndex].Lines[lineIndex].Data) == "yes"
					goto Done
				}
			}
		}
	}
Done:
	return wrapped
}

// Version returns the las file version
func (las *LAS) Version() string {
	var result string
	for sectionIndex := range las.Sections {
		if strings.ToLower(las.Sections[sectionIndex].Name) == "version information" {
			for lineIndex := range las.Sections[sectionIndex].Lines {
				if strings.ToLower(las.Sections[sectionIndex].Lines[lineIndex].Mnem) == "vers" {
					result = las.Sections[sectionIndex].Lines[lineIndex].Data
					goto Done
				}
			}
		}
	}
Done:
	return result
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
