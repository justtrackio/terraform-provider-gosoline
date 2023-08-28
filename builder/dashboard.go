package builder

import "fmt"

const (
	DashboadWidth          = 24
	PanelWidth             = 12
	PanelHeight            = 8
	orchestratorEcs        = "ecs"
	orchestratorKubernetes = "kubernetes"
)

type Dashboard struct {
	Title  string  `json:"title"`
	Panels []Panel `json:"panels"`
}

type DashboardBuilder struct {
	resourceNames  ResourceNames
	panelFactories []PanelFactory
	orchestrator   string
}

func NewDashboardBuilder(resourceNames ResourceNames, orchestrator string) *DashboardBuilder {
	return &DashboardBuilder{
		resourceNames:  resourceNames,
		panelFactories: make([]PanelFactory, 0),
		orchestrator:   orchestrator,
	}
}

func (d *DashboardBuilder) AddServiceAndTask() {
	d.AddPanel(NewPanelRow("Service Resource Usage"))
	d.AddPanel(NewPanelServiceUtilization)
	d.AddPanel(NewPanelTaskDeployment)
	for i := range d.resourceNames.Containers {
		d.AddPanel(NewPanelContainerCpuFactory(i))
		d.AddPanel(NewPanelContainerMemoryFactory(i))
	}
}

func (d *DashboardBuilder) AddElbTargetGroup(targetGroupIndex int) {
	d.AddPanel(NewPanelRow("Load Balancer"))
	d.AddPanel(NewPanelElbRequestCount(targetGroupIndex))
	d.AddPanel(NewPanelElbResponseTime(targetGroupIndex))
	d.AddPanel(NewPanelElbHttpStatus(targetGroupIndex))
	d.AddPanel(NewPanelElbHealthyHosts(targetGroupIndex))
	d.AddPanel(NewPanelElbRequestCountPerTarget(targetGroupIndex))
}

func (d *DashboardBuilder) AddTraefikService() {
	if d.orchestrator != orchestratorKubernetes {
		return
	}

	d.AddPanel(NewPanelRow("Traefik"))
	d.AddPanel(NewPanelTraefikRequestCount)
	d.AddPanel(NewPanelTraefikResponseTime)
	d.AddPanel(NewPanelTraefikHttpStatus)
	d.AddPanel(NewPanelKubernetesHealthyPods)
	d.AddPanel(NewPanelTraefikRequestCountPerTarget)
}

func (d *DashboardBuilder) AddApiServerHandler(method string, path string) {
	rowTitle := fmt.Sprintf("ApiServer: %s %s", method, path)

	d.AddPanel(NewPanelRow(rowTitle))
	d.AddPanel(NewPanelApiServerRequestCount(path))
	d.AddPanel(NewPanelApiServerResponseTime(path))
	d.AddPanel(NewPanelApiServerHttpStatus(path))
}

func (d *DashboardBuilder) AddDynamoDbTable(table MetadataCloudAwsDynamodbTable) {
	rowTitle := fmt.Sprintf("Dynamodb: %s", table.TableName)

	d.AddPanel(NewPanelRow(rowTitle))
	d.AddPanel(NewPanelDdbReadUsage(table))
	d.AddPanel(NewPanelDdbReadThrottles(table))
	d.AddPanel(NewPanelDdbWriteUsage(table))
	d.AddPanel(NewPanelDdbWriteThrottles(table))
}

func (d *DashboardBuilder) AddCloudAwsKinesisKinsumer(stream MetadataCloudAwsKinesisKinsumer) {
	rowTitle := fmt.Sprintf("Kinsumer on Stream: %s (%d Shards)", stream.StreamNameFull, stream.OpenShardCount)

	d.AddPanel(NewPanelRow(rowTitle))
	d.AddPanel(NewPanelKinesisKinsumerMillisecondsBehind(stream))
	d.AddPanel(NewPanelKinesisKinsumerMessageCounts(stream))
	d.AddPanel(NewPanelKinesisKinsumerReadOperations(stream))
	d.AddPanel(NewPanelKinesisKinsumerProcessDuration(stream))
}

func (d *DashboardBuilder) AddCloudAwsKinesisRecordWriter(stream MetadataCloudAwsKinesisRecordWriter) {
	rowTitle := fmt.Sprintf("Kinesis RecordWriter on Stream: %s (%d Shards)", stream.StreamName, stream.OpenShardCount)

	d.AddPanel(NewPanelRow(rowTitle))
	d.AddPanel(NewPanelKinesisRecordWriterPutRecordsCount(stream))
	d.AddPanel(NewPanelKinesisRecordWriterPutRecordsBatchSize(stream))
}

func (d *DashboardBuilder) AddCloudAwsKinesisStream(stream KinesisStreamAware) {
	rowTitle := fmt.Sprintf("Kinesis Stream: %s (%d Shards)", stream.GetStreamNameFull(), stream.GetOpenShardCount())

	d.AddPanel(NewPanelRow(rowTitle))

	d.AddPanel(NewPanelKinesisStreamSuccessRate(stream))
	d.AddPanel(NewPanelKinesisStreamGetRecordsBytes(stream))
	d.AddPanel(NewPanelKinesisStreamIncomingDataBytes(stream))
	d.AddPanel(NewPanelKinesisStreamIncomingDataCount(stream))
	d.AddPanel(NewPanelKinesisStreamRecordSize(stream))
}

func (d *DashboardBuilder) AddCloudAwsSqsQueue(queue MetadataCloudAwsSqsQueue) {
	rowTitle := fmt.Sprintf("SQS: %s", queue.QueueNameFull)

	d.AddPanel(NewPanelRow(rowTitle))
	d.AddPanel(NewPanelSqsMessagesVisible(queue))
	d.AddPanel(NewPanelSqsTraffic(queue))
	d.AddPanel(NewPanelSqsMessageSize(queue))
}

func (d *DashboardBuilder) AddStreamConsumer(consumer MetadataStreamConsumer) {
	rowTitle := fmt.Sprintf("Stream Consumer: %s", consumer.Name)

	d.AddPanel(NewPanelRow(rowTitle))
	d.AddPanel(NewPanelStreamConsumerProcessedCount(consumer))
	d.AddPanel(NewPanelStreamConsumerProcessDuration(consumer))

	if consumer.RetryEnabled {
		d.AddPanel(NewPanelStreamConsumerRetryActions(consumer))
	}
}

func (d *DashboardBuilder) AddStreamProducerDaemon(producer MetadataStreamProducer) {
	rowTitle := fmt.Sprintf("Stream Producer Daemon: %s", producer.Name)

	d.AddPanel(NewPanelRow(rowTitle))
	d.AddPanel(NewPanelStreamProducerDaemonSizes(producer))
	d.AddPanel(NewPanelStreamProducerMessageCount(producer))
}

func (d *DashboardBuilder) AddPanel(panel PanelFactory) {
	d.panelFactories = append(d.panelFactories, panel)
}

func (d *DashboardBuilder) Build() Dashboard {
	var x, y int
	panels := make([]Panel, len(d.panelFactories))

	for i, factory := range d.panelFactories {
		panel := d.buildPanel(factory, x, y)
		panels[i] = panel

		x += panel.GridPos.W

		if x >= DashboadWidth {
			x = 0
			y += panel.GridPos.W
		}
	}

	var dashboardTitle string
	switch d.orchestrator {
	case orchestratorEcs:
		dashboardTitle = d.resourceNames.EcsTaskDefinition
	case orchestratorKubernetes:
		dashboardTitle = d.resourceNames.KubernetesPod
	}

	return Dashboard{
		Title:  dashboardTitle,
		Panels: panels,
	}
}

func (d *DashboardBuilder) buildPanel(factory PanelFactory, x int, y int) Panel {
	panelGridPos := NewPanelGridPos(PanelHeight, PanelWidth, x, y)
	settings := newPaneSettings(d.resourceNames, panelGridPos, d.orchestrator)
	panel := factory(settings)

	if panel.FieldConfig.Defaults.Custom.AxisPlacement == "" {
		panel.FieldConfig.Defaults.Custom.AxisPlacement = "right"
	}

	if panel.FieldConfig.Defaults.Custom.LineWidth == 0 {
		panel.FieldConfig.Defaults.Custom.LineWidth = 2
	}

	if option, ok := panel.Options.(*PanelOptionsCloudWatch); ok {
		if option.Tooltip.Mode == "" {
			option.Tooltip.Mode = "multi"
		}
	}

	return panel
}
