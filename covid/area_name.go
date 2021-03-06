package covid

type areaName struct {
	value string
}

var _ Filter = &areaName{}

func (an *areaName) MetricName() string {
	return "areaName"
}

func (an *areaName) Value() string {
	return an.value
}

// AreaName returns a Filter which will retrieve data for the area name provided.
func AreaName(name string) Filter {
	return &areaName{name}
}
