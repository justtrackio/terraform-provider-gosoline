package builder

import "fmt"

func NewPanelTraefikRequestCount(settings PanelSettings) Panel {
	labelFilter := getTraefikServiceLabelFilter(settings.resourceNames.TraefikServiceName)

	return Panel{
		Datasource: datasourcePrometheus,
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Min: "0",
			},
			Overrides: []PanelFieldConfigOverride{
				NewColorPropertyOverride("Requests", "semi-dark-blue", ""),
			},
		},
		GridPos: settings.gridPos,
		Targets: []interface{}{
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`sum(irate(traefik_service_requests_total{%s}[1m])) * 60`, labelFilter),
				LegendFormat: "Requests",
				RefId:        "Requests",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Request Count",
		Type:    "timeseries",
	}
}

func NewPanelTraefikResponseTime(settings PanelSettings) Panel {
	labelFilter := getTraefikServiceLabelFilter(settings.resourceNames.TraefikServiceName)

	return Panel{
		Datasource: datasourcePrometheus,
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Min:  "0",
				Unit: "s",
			},
			Overrides: []PanelFieldConfigOverride{
				NewColorPropertyOverride("Response Time", "semi-dark-blue", ""),
			},
		},
		GridPos: settings.gridPos,
		Targets: []interface{}{
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`sum(irate(traefik_service_request_duration_seconds_sum{%s}[$__rate_interval]))`, labelFilter),
				LegendFormat: "Response Time",
				RefId:        "Requests",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Response Time",
		Type:    "timeseries",
	}
}

func NewPanelTraefikHttpStatus(settings PanelSettings) Panel {
	labelFilter := getTraefikServiceLabelFilter(settings.resourceNames.TraefikServiceName)

	return Panel{
		Datasource: datasourcePrometheus,
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Min: "0",
			},
			Overrides: []PanelFieldConfigOverride{
				NewColorPropertyOverride("HTTP 2XX", "semi-dark-green", ""),
				NewColorPropertyOverride("HTTP 3XX", "semi-dark-yellow", ""),
				NewColorPropertyOverride("HTTP 4XX", "semi-dark-orange", ""),
				NewColorPropertyOverride("HTTP 5XX", "dark-red", ""),
			},
		},
		GridPos: settings.gridPos,
		Targets: []interface{}{
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`sum(irate(traefik_service_requests_total{code=~"2.*",%s}[1m])) * 60 or vector(0)`, labelFilter),
				LegendFormat: "HTTP 2XX",
				RefId:        "A",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`sum(irate(traefik_service_requests_total{code=~"3.*",%s}[1m])) * 60 or vector(0)`, labelFilter),
				LegendFormat: "HTTP 3XX",
				RefId:        "B",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`sum(irate(traefik_service_requests_total{code=~"4.*",%s}[1m])) * 60 or vector(0)`, labelFilter),
				LegendFormat: "HTTP 4XX",
				RefId:        "C",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`sum(irate(traefik_service_requests_total{code=~"5.*",%s}[1m])) * 60 or vector(0)`, labelFilter),
				LegendFormat: "HTTP 5XX",
				RefId:        "D",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "HTTP Status Overview",
		Type:    "timeseries",
	}
}

func NewPanelKubernetesHealthyPods(settings PanelSettings) Panel {
	labelFilter := getKubernetesPodLabelFilter(settings.resourceNames.KubernetesNamespace, settings.resourceNames.KubernetesPod)

	return Panel{
		Datasource: datasourcePrometheus,
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Min: "0",
			},
			Overrides: []PanelFieldConfigOverride{},
		},
		GridPos: settings.gridPos,
		Targets: []interface{}{
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`count(kube_pod_status_ready{condition="true",%s})`, labelFilter),
				LegendFormat: "Healthy Endpoints",
				RefId:        "A",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Healthy Endpoints",
		Type:    "timeseries",
	}
}

func NewPanelTraefikRequestCountPerTarget(settings PanelSettings) Panel {
	labelFilterTraefik := getTraefikServiceLabelFilter(settings.resourceNames.TraefikServiceName)
	labelFilterPod := getKubernetesPodLabelFilter(settings.resourceNames.KubernetesNamespace, settings.resourceNames.KubernetesPod)

	return Panel{
		Datasource: datasourcePrometheus,
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Min: "0",
			},
			Overrides: []PanelFieldConfigOverride{},
		},
		GridPos: settings.gridPos,
		Targets: []interface{}{
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`sum(irate(traefik_service_requests_total{%s}[1m])) by () * 60/count(kube_pod_status_ready{condition="true",%s})`, labelFilterTraefik, labelFilterPod),
				LegendFormat: "Requests",
				RefId:        "A",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Requests Per Healthy Target",
		Type:    "timeseries",
	}
}
