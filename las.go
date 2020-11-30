package golas

import (
	"strings"
)

// LAS represents a .las file
type LAS struct {
	Sections    []Section
	Logs        LogData
	curveInfo   *Section
	otherInfo   *Section
	paramInfo   *Section
	version     string
	versionInfo *Section
	wellInfo    *Section
	wrap        string
}

// Section represents a .las file section
type Section struct {
	Name     string
	Lines    []Line
	Comments []string
}

// Line represents a header line in a .las file section
type Line struct {
	Mnem        string
	Units       string
	Data        string
	Description string
}

// LogData represents a row in the ASCII Log Data section ('~A')
type LogData [][]string

// VersionInformation returns the Version Information Section
func (las *LAS) VersionInformation() Section {
	if las.versionInfo == nil {
		las.versionInfo = las.findSection("Version Information")
	}
	return *las.versionInfo
}

// WellInformation returns the Version Information Section
func (las *LAS) WellInformation() Section {
	if las.wellInfo == nil {
		las.wellInfo = las.findSection("Well Information")
	}
	return *las.wellInfo
}

// CurveInformation returns the Version Information Section
func (las *LAS) CurveInformation() Section {
	if las.curveInfo == nil {
		las.curveInfo = las.findSection("Curve Information")
	}
	return *las.curveInfo
}

// OtherInformation returns the Version Information Section
func (las *LAS) OtherInformation() Section {
	if las.otherInfo == nil {
		las.otherInfo = las.findSection("Other Information")
	}
	return *las.otherInfo
}

// ParameterInformation returns the Version Information Section
func (las *LAS) ParameterInformation() Section {
	if las.paramInfo == nil {
		las.paramInfo = las.findSection("Parameter Information")
	}
	return *las.paramInfo
}

// Wrap returns whether or not the las file is wrapped
func (las *LAS) Wrap() Line {
	var l Line
	vi := las.VersionInformation()
	for i := 0; i < len(vi.Lines); i++ {
		if strings.ToLower(vi.Lines[i].Mnem) == "wrap" {
			l = vi.Lines[i]
			break
		}
	}
	return l
}

// Version returns whether or not the las file is wrapped
func (las *LAS) Version() Line {
	var l Line
	vi := las.VersionInformation()
	for i := 0; i < len(vi.Lines); i++ {
		if strings.ToLower(vi.Lines[i].Mnem) == "vers" {
			l = vi.Lines[i]
			break
		}
	}
	return l
}

func (las *LAS) findSection(sectionName string) *Section {
	var s *Section
	for i := 0; i < len(las.Sections); i++ {
		if strings.ToLower(las.Sections[i].Name) == strings.ToLower(sectionName) {
			s = &las.Sections[i]
			break
		}
	}
	return s
}
