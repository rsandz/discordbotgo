package message

type Role string

const (
	UserRole Role = "user"
	BotRole  Role = "bot"
)

// Represents a message from a user that should be handled.
type Request struct {
	// RequestMessage is the message that started this request.
	RequestMessage Message
	// History contains prior messages in this conversation context.
	History []Message
	// Channel identifies the medium or location of the request.
	Channel string
}

// Creates a new request.
func NewRequest(requestMessage Message, history []Message, channel string) *Request {
	return &Request{
		RequestMessage: requestMessage,
		History:        history,
		Channel:        channel,
	}
}

// Represents a response to a message.
type Response struct {
	// ResponseMessage is the content to return to the user.
	ResponseMessage Message
	// ShouldContinueHandling indicates whether the request requires further processing by other handlers.
	ShouldContinueHandling bool
}

// Creates a new response.
func NewResponse(responseMessage Message, shouldContinueHandling bool) *Response {
	return &Response{
		ResponseMessage:        responseMessage,
		ShouldContinueHandling: shouldContinueHandling,
	}
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

// Creates a new message.
func NewMessage(user string, content string, role Role) *Message {
	return &Message{
		User:    user,
		Content: content,
		Role:    role,
	}
}
