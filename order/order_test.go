package order_test

import (
	"testing"

	"github.com/jamieabc/go-abc/order"
	"github.com/stretchr/testify/assert"

	"github.com/jamieabc/go-abc/report"
)

func TestOrderWhenOrderByScore(t *testing.T) {
	reports := []report.Report{
		{
			Path:  "abc",
			Score: 20,
			ABC:   report.ABC{},
		},
		{
			Path:  "def",
			Score: 10,
			ABC:   report.ABC{},
		},
	}

	o := order.NewOrderByScore(reports)
	o.Sort()
	result := o.Report()

	assert.Equal(t, 2, len(result), "wrong length")
	assert.Equal(t, 10, result[0].Score, "wrong first item")
	assert.Equal(t, 20, result[1].Score, "wrong second item")
}
