package builder

type MetadataApplication struct {
	ApiServer MetadataApiServer `json:"apiserver"`
	Cloud     MetadataCloud     `json:"cloud"`
}

type MetadataApiServer struct {
	Routes []MetadataApiServerRoute `json:"routes"`
}

type MetadataApiServerRoute struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type MetadataCloud struct {
	Aws MetadataCloudAws `json:"aws"`
}

type MetadataCloudAws struct {
	Dynamodb MetadataCloudAwsDynamodb `json:"dynamodb"`
	Sqs      MetadataCloudAwsSqs      `json:"sqs"`
}

type MetadataCloudAwsDynamodb struct {
	Tables []string `json:"tables"`
}

type MetadataCloudAwsSqs struct {
	Queues []string `json:"queues"`
}
