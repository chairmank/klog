package report

import (
	"fmt"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/src/service/period"
)

type yearAggregator struct{}

func NewYearAggregator() Aggregator {
	return &yearAggregator{}
}

func (a *yearAggregator) NumberOfPrefixColumns() int {
	return 1
}

func (a *yearAggregator) DateHash(date Date) period.Hash {
	return period.Hash(period.NewYearFromDate(date).Hash())
}

func (a *yearAggregator) OnHeaderPrefix(table *terminalformat.Table) {
	table.
		CellL("    ") // 2020
}

func (a *yearAggregator) OnRowPrefix(table *terminalformat.Table, date Date) {
	// Year
	table.CellR(fmt.Sprint(date.Year()))
}
