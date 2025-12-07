package message

type Role string

const (
	UserRole Role = "user"
	BotRole  Role = "bot"
)

// Represents a message from a user that should be handled.
type Request struct {
	RequestMessage Message
	History        []Message
	Channel        string
}

// Represents a response to a message.
type Response struct {
	ResponseMessage Message
}

// Represents a chat message.
type Message struct {
	// User who sent this message
	User string
	// Content of the message
	Content string
	// Role of the user
	Role Role
}
