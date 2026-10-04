// Package notify defines the LINE notification port. The real LINE
// implementation is added in the notifications slice (P11).
package notify

import (
	"context"
	"sync"
)

// Status is the delivery outcome. Only StatusSent means LINE accepted the message.
type Status string

const (
	StatusSent    Status = "sent"
	StatusSkipped Status = "skipped"
	StatusFailed  Status = "failed"
)

// Message is one push to one LINE user.
type Message struct {
	LineUserID string
	Text       string
}

// Notifier sends LINE messages.
type Notifier interface {
	Send(ctx context.Context, msg Message) (Status, error)
}

// Noop is a NON-DELIVERING notifier for tests and local runs. It records each
// message and always returns StatusSkipped; it never reports a message as sent.
type Noop struct {
	mu       sync.Mutex
	messages []Message
}

// Send records msg and reports it as skipped.
func (n *Noop) Send(_ context.Context, msg Message) (Status, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.messages = append(n.messages, msg)
	return StatusSkipped, nil
}

// Messages returns a copy of every message passed to Send.
func (n *Noop) Messages() []Message {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]Message(nil), n.messages...)
}
