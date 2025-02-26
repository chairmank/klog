/*
Package reconciling contains logic to manipulate klog source text.
The idea of the reconciler generally is to add or modify serialised records
in a minimally invasive manner. Instead or re-serialising the record itself,
it tries to find the location in the original text and modify that directly.
While this approach might feel a little hacky, it avoids lots of other
complications and sources of bugs, that could potentially mess up user data.
*/
package reconciling

import (
	"errors"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/engine"
	"strings"
)

// Reconciler is a mechanism to manipulate record data in a file.
type Reconciler struct {
	record          Record
	style           *parser.Style
	lastLinePointer int // Line index of the last entry
	lines           []engine.Line
	recordPointer   int
}

// Result is the result of an applied reconciler.
type Result struct {
	Record        Record
	AllRecords    []parser.ParsedRecord
	AllSerialised string
}

// Reconcile is a function interface for applying a reconciler.
type Reconcile func(*Reconciler) (*Result, error)

func countLines(es []Entry) int {
	result := 0
	for _, e := range es {
		result += len(e.Summary())
	}
	return result
}

// MakeResult returns the reconciled data.
func (r *Reconciler) MakeResult() (*Result, error) {
	text := ""
	for _, l := range r.lines {
		text += l.Original()
	}

	// As a safeguard, make sure the result is parseable.
	newRecords, errs := parser.Parse(text)
	if errs != nil {
		return nil, errors.New("This operation wouldn’t result in a valid record")
	}

	return &Result{
		Record:        newRecords[r.recordPointer],
		AllRecords:    newRecords,
		AllSerialised: text,
	}, nil
}

// findOpenRangeIndex returns the index of the open range entry, or -1 if no open range.
func (r *Reconciler) findOpenRangeIndex() int {
	openRangeEntryIndex := -1
	for i, e := range r.record.Entries() {
		Unbox(&e,
			func(Range) any { return nil },
			func(Duration) any { return nil },
			func(OpenRange) any {
				openRangeEntryIndex = i
				return nil
			},
		)
	}
	return openRangeEntryIndex
}

var blankLine = insertableText{"", 0}

type insertableText struct {
	text        string
	indentation int
}

func (r *Reconciler) insert(lineIndex int, texts []insertableText) {
	result := make([]engine.Line, len(r.lines)+len(texts))
	offset := 0
	for i := range result {
		if i >= lineIndex && offset < len(texts) {
			line := strings.Repeat(r.style.Indentation.Get(), texts[offset].indentation)
			line += texts[offset].text + r.style.LineEnding.Get()
			result[i] = engine.NewLineFromString(line, -1)
			offset++
		} else {
			result[i] = r.lines[i-offset]
		}
		result[i].LineNumber = i + 1
	}
	if lineIndex > 0 && result[lineIndex-1].LineEnding == "" {
		result[lineIndex-1].LineEnding = r.style.LineEnding.Get()
	}
	r.lines = result
}

func toMultilineEntryTexts(entryValue string, entrySummary string) []insertableText {
	var result []insertableText
	// Normalize potentially redundant escaping.
	entrySummary = strings.ReplaceAll(entrySummary, "\\n", "\n")
	for i, summaryLine := range strings.Split(entrySummary, "\n") {
		indent := 1
		if i > 0 {
			indent = 2
		}
		text := ""
		if i == 0 {
			text += entryValue
			if text != "" && summaryLine != "" {
				text += " "
			}
		}
		text += summaryLine
		result = append(result, insertableText{text, indent})
	}
	return result
}
