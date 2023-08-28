package builder

func newPaneSettings(resourceNames ResourceNames, gridPos PanelGridPos, orchestrator string) PanelSettings {
	return PanelSettings{
		resourceNames: resourceNames,
		gridPos:       gridPos,
		orchestrator:  orchestrator,
	}
}

type PanelSettings struct {
	resourceNames ResourceNames
	gridPos       PanelGridPos
	orchestrator  string
}

type PanelFactory func(settings PanelSettings) Panel

type Panel struct {
	Collapsed   bool             `json:"collapsed,omitempty"`
	Datasource  string           `json:"datasource"`
	FieldConfig PanelFieldConfig `json:"fieldConfig"`
	GridPos     PanelGridPos     `json:"gridPos"`
	Options     interface{}      `json:"options"`
	Targets     []interface{}    `json:"targets"`
	Title       string           `json:"title"`
	Type        string           `json:"type"`
	Panels      []Panel          `json:"panels"`
}

type PanelFieldConfig struct {
	Defaults  PanelFieldConfigDefaults    `json:"defaults"`
	Overrides []PanelFieldConfigOverwrite `json:"overrides"`
}

type PanelFieldConfigDefaults struct {
	Custom     PanelFieldConfigDefaultsCustom     `json:"custom"`
	Max        string                             `json:"max,omitempty"`
	Min        string                             `json:"min,omitempty"`
	Thresholds PanelFieldConfigDefaultsThresholds `json:"thresholds"`
	Unit       string                             `json:"unit,omitempty"`
}

type PanelFieldConfigDefaultsThresholds struct {
	Mode  string                                   `json:"mode"`
	Steps []PanelFieldConfigDefaultsThresholdsStep `json:"steps"`
}

type PanelFieldConfigDefaultsThresholdsStep struct {
	Color string `json:"color"`
	Value int    `json:"value"`
}

type PanelFieldConfigDefaultsCustom struct {
	AxisPlacement   string          `json:"axisPlacement"`
	LineWidth       int             `json:"lineWidth"`
	ThresholdsStyle ThresholdsStyle `json:"thresholdsStyle"`
	SpanNulls       bool            `json:"spanNulls"`
}

type ThresholdsStyle struct {
	Mode string `json:"mode"`
}

type PanelFieldConfigOverwrite struct {
	Matcher    OverwriteMatcher    `json:"matcher"`
	Properties []OverwriteProperty `json:"properties"`
}

type OverwriteMatcher struct {
	Id      string `json:"id"`
	Options string `json:"options"`
}

type OverwriteProperty struct {
	Id    string                 `json:"id"`
	Value OverwritePropertyValue `json:"value"`
}

type OverwritePropertyValue struct {
	FixedColor string `json:"fixedColor"`
	Mode       string `json:"mode"`
}

func NewColorPropertyOverwrite(metric string, color string) PanelFieldConfigOverwrite {
	return PanelFieldConfigOverwrite{
		Matcher: OverwriteMatcher{
			Id:      "byName",
			Options: metric,
		},
		Properties: []OverwriteProperty{
			{
				Id: "color",
				Value: OverwritePropertyValue{
					FixedColor: color,
					Mode:       "fixed",
				},
			},
		},
	}
}

type PanelGridPos struct {
	H int `json:"h"`
	W int `json:"w"`
	X int `json:"x"`
	Y int `json:"y"`
}

func NewPanelGridPos(h int, w int, x int, y int) PanelGridPos {
	return PanelGridPos{
		H: h,
		W: w,
		X: x,
		Y: y,
	}
}

type PanelOptionsCloudWatch struct {
	Tooltip PanelOptionsTooltip `json:"tooltip"`
}

type PanelOptionsTooltip struct {
	Mode string `json:"mode"`
}

type PanelTargetCloudWatch struct {
	Alias      string            `json:"alias"`
	Dimensions map[string]string `json:"dimensions"`
	Expression string            `json:"expression"`
	Id         string            `json:"id"`
	Hide       bool              `json:"hide"`
	MatchExact bool              `json:"matchExact"`
	MetricName string            `json:"metricName"`
	Namespace  string            `json:"namespace"`
	Period     string            `json:"period"`
	RefId      string            `json:"refId"`
	Region     string            `json:"region"`
	Statistics []string          `json:"statistics"`
}

type PanelTargetPrometheus struct {
	Exemplar     bool   `json:"exemplar"`
	Expression   string `json:"expr"`
	Hide         bool   `json:"hide"`
	Interval     string `json:"interval"`
	LegendFormat string `json:"legendFormat"`
	RefId        string `json:"refId"`
}

type PanelTargetElasticsearch struct {
	RefId     string                           `json:"refId"`
	Query     string                           `json:"query"`
	Metrics   []PanelTargetElasticsearchMetric `json:"metrics"`
	TimeField string                           `json:"timeField"`
}

type PanelTargetElasticsearchMetric struct {
	Id       string                                 `json:"id"`
	Type     string                                 `json:"type"`
	Settings PanelTargetElasticsearchMetricSettings `json:"settings"`
}

type PanelTargetElasticsearchMetricSettings struct {
	Limit string `json:"limit"`
}

type PanelOptionsElasticsearch struct {
	ShowTime           bool   `json:"showTime"`
	ShowLabels         bool   `json:"showLabels"`
	ShowCommonLabels   bool   `json:"showCommonLabels"`
	WrapLogMessage     bool   `json:"wrapLogMessage"`
	PrettifyLogMessage bool   `json:"prettifyLogMessage"`
	EnableLogDetails   bool   `json:"enableLogDetails"`
	DedupStrategy      string `json:"dedupStrategy"`
	SortOrder          string `json:"sortOrder"`
}
