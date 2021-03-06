package covid

import "time"

type date struct {
	value time.Time
}

var _ Filter = &date{}

func (d *date) MetricName() string {
	return "date"
}

func (d *date) Value() string {
	return d.value.Format("2006-01-02")
}

// Date returns a Filter which will retrieve data for the date provided.
func Date(t time.Time) Filter {
	return &date{t}
}
