package builder

func NewPanelDdbReadUsage(table MetadataCloudAwsDynamodbTable) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(table.AwsClientName),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						SpanNulls: true,
					},
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{
					NewColorPropertyOverwrite("Provisioned", "dark-red"),
					NewColorPropertyOverwrite("Consumed", "super-light-blue"),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Provisioned",
					Dimensions: map[string]string{
						"TableName": table.TableName,
					},
					Expression: "",
					Id:         "m2",
					MatchExact: true,
					MetricName: "ProvisionedReadCapacityUnits",
					Namespace:  "AWS/DynamoDB",
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
				PanelTargetCloudWatch{
					Alias: "",
					Dimensions: map[string]string{
						"TableName": table.TableName,
					},
					Expression: "",
					Id:         "m1",
					Hide:       true,
					MatchExact: true,
					MetricName: "ConsumedReadCapacityUnits",
					Namespace:  "AWS/DynamoDB",
					Period:     "",
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias:      "Consumed",
					Dimensions: map[string]string{},
					Expression: "m1/PERIOD(m1)",
					Id:         "",
					MatchExact: true,
					MetricName: "",
					Namespace:  "",
					Period:     "",
					RefId:      "C",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Read Usage (average units/second)",
			Type:    "timeseries",
		}
	}
}

func NewPanelDdbReadThrottles(table MetadataCloudAwsDynamodbTable) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(table.AwsClientName),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						AxisPlacement: "right",
						LineWidth:     2,
						SpanNulls:     true,
					},
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "GetItem",
					Dimensions: map[string]string{
						"TableName": table.TableName,
						"Operation": "GetItem",
					},
					MatchExact: true,
					MetricName: "ThrottledRequests",
					Namespace:  "AWS/DynamoDB",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "Scan",
					Dimensions: map[string]string{
						"TableName": table.TableName,
						"Operation": "Scan",
					},
					MatchExact: true,
					MetricName: "ThrottledRequests",
					Namespace:  "AWS/DynamoDB",
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "Query",
					Dimensions: map[string]string{
						"TableName": table.TableName,
						"Operation": "Query",
					},
					MatchExact: true,
					MetricName: "ThrottledRequests",
					Namespace:  "AWS/DynamoDB",
					RefId:      "C",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "BatchGetItem",
					Dimensions: map[string]string{
						"TableName": table.TableName,
						"Operation": "BatchGetItem",
					},
					MatchExact: true,
					MetricName: "ThrottledRequests",
					Namespace:  "AWS/DynamoDB",
					RefId:      "D",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Read throttled requests (count)",
			Type:    "timeseries",
		}
	}
}

func NewPanelDdbWriteUsage(table MetadataCloudAwsDynamodbTable) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(table.AwsClientName),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						AxisPlacement: "right",
						LineWidth:     2,
						SpanNulls:     true,
					},
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{
					NewColorPropertyOverwrite("Provisioned", "dark-red"),
					NewColorPropertyOverwrite("Consumed", "super-light-blue"),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Provisioned",
					Dimensions: map[string]string{
						"TableName": table.TableName,
					},
					Expression: "",
					Id:         "m2",
					MatchExact: true,
					MetricName: "ProvisionedWriteCapacityUnits",
					Namespace:  "AWS/DynamoDB",
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
				PanelTargetCloudWatch{
					Alias: "",
					Dimensions: map[string]string{
						"TableName": table.TableName,
					},
					Expression: "",
					Id:         "m1",
					Hide:       true,
					MatchExact: true,
					MetricName: "ConsumedWriteCapacityUnits",
					Namespace:  "AWS/DynamoDB",
					Period:     "",
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias:      "Consumed",
					Dimensions: map[string]string{},
					Expression: "m1/PERIOD(m1)",
					Id:         "",
					MatchExact: true,
					MetricName: "",
					Namespace:  "",
					Period:     "",
					RefId:      "C",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Write Usage (average units/second)",
			Type:    "timeseries",
		}
	}
}

func NewPanelDdbWriteThrottles(table MetadataCloudAwsDynamodbTable) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(table.AwsClientName),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						AxisPlacement: "right",
						LineWidth:     2,
						SpanNulls:     true,
					},
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "PutItem",
					Dimensions: map[string]string{
						"TableName": table.TableName,
						"Operation": "PutItem",
					},
					MatchExact: true,
					MetricName: "ThrottledRequests",
					Namespace:  "AWS/DynamoDB",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "UpdateItem",
					Dimensions: map[string]string{
						"TableName": table.TableName,
						"Operation": "UpdateItem",
					},
					MatchExact: true,
					MetricName: "ThrottledRequests",
					Namespace:  "AWS/DynamoDB",
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "DeleteItem",
					Dimensions: map[string]string{
						"TableName": table.TableName,
						"Operation": "DeleteItem",
					},
					MatchExact: true,
					MetricName: "ThrottledRequests",
					Namespace:  "AWS/DynamoDB",
					RefId:      "C",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "BatchWriteItem",
					Dimensions: map[string]string{
						"TableName": table.TableName,
						"Operation": "BatchWriteItem",
					},
					MatchExact: true,
					MetricName: "ThrottledRequests",
					Namespace:  "AWS/DynamoDB",
					RefId:      "D",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Write throttled requests (count)",
			Type:    "timeseries",
		}
	}
}
