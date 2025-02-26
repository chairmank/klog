package klog

import (
	"errors"
	"regexp"
)

// RecordSummary contains the summary lines of the overall summary that
// appears underneath the date of a record.
type RecordSummary []string

// EntrySummary contains the summary line that appears behind the time value
// of an entry.
type EntrySummary []string

// NewRecordSummary creates a new RecordSummary from individual lines of text.
// None of the lines can start with blank characters, and none of the lines
// can be empty or blank.
func NewRecordSummary(line ...string) (RecordSummary, error) {
	for _, l := range line {
		if len(l) == 0 || regexp.MustCompile(`^[\p{Zs}\t]`).MatchString(l) {
			return nil, errors.New("MALFORMED_SUMMARY")
		}
	}
	return line, nil
}

// NewEntrySummary creates an EntrySummary from individual lines of text.
// Except for the first line, none of the lines can be empty or blank.
func NewEntrySummary(line ...string) (EntrySummary, error) {
	for i, l := range line {
		if i == 0 {
			continue
		}
		if len(l) == 0 || regexp.MustCompile(`^[\p{Zs}\t]*$`).MatchString(l) {
			return nil, errors.New("MALFORMED_SUMMARY")
		}
	}
	return line, nil
}

func (s RecordSummary) Lines() []string {
	return s
}

func (s EntrySummary) Lines() []string {
	return RecordSummary(s).Lines()
}

func (s RecordSummary) Tags() TagSet {
	tags := NewEmptyTagSet()
	for _, l := range s {
		for _, m := range HashTagPattern.FindAllStringSubmatch(l, -1) {
			tag, _ := NewTagFromString(m[0])
			tags.Put(tag)
		}
	}
	return tags
}

// Tags returns the tags that the entry summary contains.
func (s EntrySummary) Tags() TagSet {
	return RecordSummary(s).Tags()
}

func (s RecordSummary) Equals(summary RecordSummary) bool {
	if len(s) != len(summary) {
		return false
	}
	for i, l := range s {
		if l != summary[i] {
			return false
		}
	}
	return true
}

func (s EntrySummary) Equals(summary EntrySummary) bool {
	if len(s) == 1 && s[0] == "" && summary == nil {
		// In the case of entry summary, an empty one matches nil.
		return true
	}
	return RecordSummary(s).Equals(RecordSummary(summary))
}
