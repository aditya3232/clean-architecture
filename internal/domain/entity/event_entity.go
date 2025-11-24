package entity

type PublishMessage struct {
	Email     string `json:"email"`
	Message   string `json:"message"`
	UserId    int64  `json:"user_id"`
	Subject   string `json:"subject"`
	QueueName string `json:"queue_name"`
}

type KafkaEvent struct {
	Name string `json:"name"`
}

type KafkaMetaData struct {
	Sender    string `json:"sender"`
	SendingAt string `json:"sendingAt"`
}

type KafkaData struct {
	ReceiverEmail    string `json:"receiver_email"`
	Message          string `json:"message"`
	ReceiverId       int64  `json:"receiver_id"`
	Subject          string `json:"subject"`
	NotificationType string `json:"notification_type"`
}

type KafkaBody struct {
	Type string     `json:"type"`
	Data *KafkaData `json:"data"`
}

type KafkaMessage struct {
	Event    KafkaEvent    `json:"event"`
	Metadata KafkaMetaData `json:"metadata"`
	Body     KafkaBody     `json:"body"`
}
