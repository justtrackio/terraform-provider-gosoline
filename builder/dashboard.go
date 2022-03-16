package builder

import "fmt"

const (
	DashboadWidth = 24
	PanelWidth    = 12
	PanelHeight   = 8
)

type Dashboard struct {
	Title  string  `json:"title"`
	Panels []Panel `json:"panels"`
}

type DashboardBuilder struct {
	appId          AppId
	panelFactories []PanelFactory
}

func NewDashboardBuilder(appId AppId) *DashboardBuilder {
	return &DashboardBuilder{
		appId:          appId,
		panelFactories: make([]PanelFactory, 0),
	}
}

func (d *DashboardBuilder) AddEcs() {
	d.AddPanel(NewPanelRow("Service Resource Usage"))
	d.AddPanel(NewPanelEcsUtilization)
	d.AddPanel(NewPanelEcsDeployment)
	d.AddPanel(NewPanelEcsCpu)
	d.AddPanel(NewPanelEcsMemory)
}

func (d *DashboardBuilder) AddElbTargetGroup(targetGroup ElbTargetGroup) {
	d.AddPanel(NewPanelRow("Load Balancer"))
	d.AddPanel(NewPanelElbRequestCount(targetGroup))
	d.AddPanel(NewPanelElbResponseTime(targetGroup))
	d.AddPanel(NewPanelElbHttpStatus(targetGroup))
	d.AddPanel(NewPanelElbHealthyHosts(targetGroup))
	d.AddPanel(NewPanelElbRequestCountPerTarget(targetGroup))
}

func (d *DashboardBuilder) AddApiServerHandler(method string, path string) {
	rowTitle := fmt.Sprintf("ApiServer: %s %s", method, path)

	d.AddPanel(NewPanelRow(rowTitle))
	d.AddPanel(NewPanelApiServerRequestCount(path))
	d.AddPanel(NewPanelApiServerResponseTime(path))
	d.AddPanel(NewPanelApiServerHttpStatus(path))
}

func (d *DashboardBuilder) AddDynamoDbTable(table string) {
	rowTitle := fmt.Sprintf("Dynamodb: %s", table)

	d.AddPanel(NewPanelRow(rowTitle))
	d.AddPanel(NewPanelDdbReadUsage(table))
	d.AddPanel(NewPanelDdbReadThrottles(table))
	d.AddPanel(NewPanelDdbWriteUsage(table))
	d.AddPanel(NewPanelDdbWriteThrottles(table))
}

func (d *DashboardBuilder) AddCloudAwsKinesisKinsumer(stream string, shardCount int) {
	rowTitle := fmt.Sprintf("Kinsumer on Stream: %s (%d Shards)", stream, shardCount)

	d.AddPanel(NewPanelRow(rowTitle))
	d.AddPanel(NewPanelKinesisKinsumerMillisecondsBehind(stream))
	d.AddPanel(NewPanelKinesisKinsumerMessageCounts(stream))
	d.AddPanel(NewPanelKinesisKinsumerReadOperations(stream, shardCount))
	d.AddPanel(NewPanelKinesisKinsumerProcessDuration(stream))
}

func (d *DashboardBuilder) AddCloudAwsKinesisRecordWriter(stream string, shardCount int) {
	rowTitle := fmt.Sprintf("Kinesis RecordWriter on Stream: %s (%d Shards)", stream, shardCount)

	d.AddPanel(NewPanelRow(rowTitle))
	d.AddPanel(NewPanelKinesisRecordWriterPutRecordsCount(stream))
	d.AddPanel(NewPanelKinesisRecordWriterPutRecordsBatchSize(stream, shardCount))
}

func (d *DashboardBuilder) AddCloudAwsKinesisStream(stream string, shardCount int) {
	rowTitle := fmt.Sprintf("Kinesis Stream: %s (%d Shards)", stream, shardCount)

	d.AddPanel(NewPanelRow(rowTitle))

	d.AddPanel(NewPanelKinesisStreamSuccessRate(stream))
	d.AddPanel(NewPanelKinesisStreamGetRecordsBytes(stream, shardCount))
	d.AddPanel(NewPanelKinesisStreamIncomingDataBytes(stream, shardCount))
	d.AddPanel(NewPanelKinesisStreamIncomingDataCount(stream, shardCount))
	d.AddPanel(NewPanelKinesisStreamRecordSize(stream))
}

func (d *DashboardBuilder) AddCloudAwsSqsQueue(queue string) {
	rowTitle := fmt.Sprintf("SQS: %s", queue)

	d.AddPanel(NewPanelRow(rowTitle))
	d.AddPanel(NewPanelSqsMessagesVisible(queue))
	d.AddPanel(NewPanelSqsTraffic(queue))
	d.AddPanel(NewPanelSqsMessageSize(queue))
}

func (d *DashboardBuilder) AddStreamConsumer(consumer MetadataStreamConsumer) {
	rowTitle := fmt.Sprintf("Stream Consumer: %s", consumer.Name)

	d.AddPanel(NewPanelRow(rowTitle))
	d.AddPanel(NewPanelStreamConsumerProcessedCount(consumer.Name))
	d.AddPanel(NewPanelStreamConsumerProcessDuration(consumer.Name))

	if consumer.RetryEnabled {
		d.AddPanel(NewPanelStreamConsumerRetryActions(consumer.Name, consumer.RetryType))
	}
}

func (d *DashboardBuilder) AddStreamProducerDaemon(producer string) {
	rowTitle := fmt.Sprintf("Stream Producer Daemon: %s", producer)

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

	return Dashboard{
		Title:  d.appId.Application,
		Panels: panels,
	}
}

func (d *DashboardBuilder) buildPanel(factory PanelFactory, x int, y int) Panel {
	panel := factory(d.appId, NewPanelGridPos(PanelHeight, PanelWidth, x, y))

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
