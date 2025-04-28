package models

import "time"

type Train struct {
	ID               int
	Line             string
	LastUpdated      time.Time
	PreviousSeverity int
	Severity         int
}

func (t Train) IsDisrupted() bool {
	return t.Severity == 6 || t.Severity == 9
}

func (t Train) SeverityMessage() string {
	return severity[t.Severity]
}

func (t Train) HasSameSeverity() bool {
	return t.PreviousSeverity == t.Severity
}

var severity = map[int]string{
	1:  "Closed",
	2:  "Suspended",
	3:  "Part Suspended",
	4:  "Planned Closure",
	5:  "Part Closure",
	6:  "Severe Delays",
	7:  "Reduced service",
	8:  "Bus Service",
	9:  "Minor Delays",
	10: "Good Service",
	11: "Part Closed",
	12: "Exit Only",
	13: "No Step Free Access",
	14: "Change of Frequency",
	15: "Diverted",
	16: "Not Running",
	17: "Issues Reported",
	18: "No Issues",
	19: "Information",
	20: "Service Closed",
}
