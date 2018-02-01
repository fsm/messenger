package messenger

type messageReceivedCallback struct {
	Object string         `json:"object"`
	Entry  []messageEntry `json:"entry"`
}

type messageEntry struct {
	ID              string           `json:"id"`
	Time            int64            `json:"time"`
	MessagingEvents []messagingEvent `json:"messaging"`
}

type messagingEvent struct {
	Sender    sender    `json:"sender"`
	Recipient recipient `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Message   message   `json:"message"`
}

type sender struct {
	ID string `json:"id"`
}

type recipient struct {
	ID string `json:"id"`
}

type message struct {
	MessageID string `json:"mid"`
	Sequence  int64  `json:"seq"`
	Text      string `json:"text"`
}
