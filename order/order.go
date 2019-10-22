package order

import (
	"sort"

	"github.com/jamieabc/go-abc/report"
)

// Sort - order interface
type Order interface {
	Sort()
	Report() []report.Report
}

func NewOrderByScore(reports []report.Report) Order {
	return &orderByScore{reports: reports}
}

type orderByScore struct {
	reports []report.Report
}

func (o *orderByScore) Len() int {
	return len(o.reports)
}

func (o *orderByScore) Less(i, j int) bool {
	return o.reports[i].Score < o.reports[j].Score
}

func (o *orderByScore) Swap(i, j int) {
	o.reports[i], o.reports[j] = o.reports[j], o.reports[i]
}

func (o *orderByScore) Sort() {
	sort.Sort(o)
}

func (o *orderByScore) Report() []report.Report {
	return o.reports
}
