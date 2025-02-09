package sqs

type Event struct {
	Records []Record
}

type Record struct {
	EventVersion      string            `json:"eventVersion"`
	EventSource       string            `json:"eventSource"`
	AWSRegion         string            `json:"awsRegion"`
	EventTime         string            `json:"eventTime"`
	EventName         string            `json:"eventName"`
	UserIdentity      UserIdentity      `json:"userIdentity"`
	RequestParameters RequestParameters `json:"requestParameters"`
	ResponseElements  ResponseElements  `json:"responseElements"`
	S3                S3                `json:"s3"`
}

type UserIdentity struct {
	PrincipalID string `json:"principalId"`
}

type RequestParameters struct {
	SourceIPAddress string `json:"sourceIPAddress"`
}

type ResponseElements struct {
	XAMZRequestID string `json:"x-amz-request-id"`
	XAMZID2       string `json:"x-amz-id-2"`
}

type S3 struct {
	S3SchemaVersion string `json:"s3SchemaVersion"`
	ConfigurationID string `json:"configurationId"`
	Bucket          Bucket `json:"bucket"`
	Object          Object `json:"object"`
}

type Bucket struct {
	Name          string        `json:"name"`
	OwnerIdentity OwnerIdentity `json:"ownerIdentity"`
	ARN           string
}

type OwnerIdentity struct {
	PrincipalID string `json:"principalId"`
}

type Object struct {
	Key       string `json:"key"`
	Size      int    `json:"size"`
	ETag      string `json:"eTag"`
	Sequencer string `json:"sequencer"`
}
