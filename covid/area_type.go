package covid

// AreaType indicates the type of geographical area data is being retrieved for.
type AreaType struct {
	value string
}

var _ Filter = &AreaType{}

func (at *AreaType) MetricName() string {
	return "areaType"
}

func (at *AreaType) Value() string {
	return at.value
}

// AreaTypeOverview Overview data for the United Kingdom.
func AreaTypeOverview() *AreaType {
	return &AreaType{"overview"}
}

// AreaTypeNation Nation data (England, Northern Ireland, Scotland, and Wales).
func AreaTypeNation() *AreaType {
	return &AreaType{"nation"}
}

// AreaTypeRegion Region data.
func AreaTypeRegion() *AreaType {
	return &AreaType{"region"}
}

// AreaTypeNHSRegion NHS Region data.
func AreaTypeNHSRegion() *AreaType {
	return &AreaType{"nhsRegion"}
}

// AreaTypeUTLA Upper-tier local authority data.
func AreaTypeUTLA() *AreaType {
	return &AreaType{"utla"}
}

// AreaTypeLTLA Lower-tier local authority data.
func AreaTypeLTLA() *AreaType {
	return &AreaType{"ltla"}
}
