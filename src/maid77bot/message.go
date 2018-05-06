package maid77bot

// MessageButton represents a button in reply message
type MessageButton struct {
	DisplayText  string
	CallbackText string
}

// UserMessage is the input messages of dialog
type UserMessage struct {
	UserId string
	Text   string
}

// ReplyMessage is the output messages of dialog
type ReplyMessage struct {
	Text       string
	Buttons    []*MessageButton
	Confidence float64
}
