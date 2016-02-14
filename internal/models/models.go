package models

import (
	"regexp"
	"sort"
	"strings"
)

type Expression struct {
	Value      string
	IsStar     bool
	IsVariadic bool
	IsWriter   bool
	Underlying string
}

func (e *Expression) String() string {
	value := e.Value
	if e.IsStar {
		value = "*" + value
	}
	if e.IsVariadic {
		return "[]" + value
	}
	return value
}

type Field struct {
	Name  string
	Type  *Expression
	Index int
}

func (f *Field) IsWriter() bool {
	return f.Type.IsWriter
}

func (f *Field) IsStruct() bool {
	return strings.HasPrefix(f.Type.Underlying, "struct")
}

func (f *Field) IsBasicType() bool {
	return isBasicType(f.Type.String()) || isBasicType(f.Type.Underlying)
}

func isBasicType(t string) bool {
	switch t {
	case "bool", "string", "int", "int8", "int16", "int32", "int64", "uint",
		"uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "rune",
		"float32", "float64", "complex64", "complex128":
		return true
	default:
		return false
	}
}

func (f *Field) IsNamed() bool {
	return f.Name != "" && f.Name != "_"
}

func (f *Field) ShortName() string {
	return strings.ToLower(string([]rune(f.Type.Value)[0]))
}

type Receiver struct {
	*Field
	Fields []*Field
}

type Function struct {
	Name         string
	IsExported   bool
	Receiver     *Receiver
	Parameters   []*Field
	Results      []*Field
	ReturnsError bool
}

func (f *Function) TestParameters() []*Field {
	var ps []*Field
	for _, p := range f.Parameters {
		if p.IsWriter() {
			continue
		}
		ps = append(ps, p)
	}
	return ps
}

func (f *Function) TestResults() []*Field {
	var ps []*Field
	ps = append(ps, f.Results...)
	for _, p := range f.Parameters {
		if !p.IsWriter() {
			continue
		}
		ps = append(ps, &Field{
			Name: p.Name,
			Type: &Expression{
				Value:      "string",
				IsWriter:   true,
				Underlying: "string",
			},
			Index: len(ps),
		})
	}
	return ps
}

func (f *Function) ReturnsMultiple() bool {
	return len(f.Results) > 1
}

func (f *Function) OnlyReturnsOneValue() bool {
	return len(f.Results) == 1 && !f.ReturnsError
}

func (f *Function) OnlyReturnsError() bool {
	return len(f.Results) == 0 && f.ReturnsError
}

func (f *Function) TestName() string {
	var r string
	if f.Receiver != nil {
		r = f.Receiver.Type.Value
	}
	return "Test" + strings.Title(r) + strings.Title(f.Name)
}

type Header struct {
	Package string
	Imports []*Import
	Code    []byte
}

type Import struct {
	Name, Path string
}

type SourceInfo struct {
	Header *Header
	Funcs  []*Function
}

func (i *SourceInfo) TestableFuncs(only, excl *regexp.Regexp, testFuncs []string) []*Function {
	sort.Strings(testFuncs)
	var fs []*Function
	for _, f := range i.Funcs {
		if f.Receiver == nil && len(f.Parameters) == 0 && len(f.Results) == 0 {
			continue
		}
		if len(testFuncs) > 0 && contains(testFuncs, f.TestName()) {
			continue
		}
		if excl != nil && excl.MatchString(f.Name) {
			continue
		}
		if only != nil && !only.MatchString(f.Name) {
			continue
		}
		fs = append(fs, f)
	}
	return fs
}

func (i *SourceInfo) UsesReflection() bool {
	for _, f := range i.Funcs {
		for _, r := range f.Results {
			if !r.IsBasicType() {
				return true
			}
		}
	}
	return false
}

func contains(ss []string, s string) bool {
	if i := sort.SearchStrings(ss, s); i < len(ss) && ss[i] == s {
		return true
	}
	return false
}

type Path string

func (p Path) TestPath() string {
	if p.IsTestPath() {
		return string(p)
	}
	return strings.TrimSuffix(string(p), ".go") + "_test.go"
}

func (p Path) IsTestPath() bool {
	return strings.HasSuffix(string(p), "_test.go")
}