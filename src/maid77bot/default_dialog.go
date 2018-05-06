package maid77bot

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

// DefaultDialog is the dialog to handle out-of-domain user messages
type DefaultDialog struct {
}

// NewDefaultDialog creates a new instance of DefaultDialog
func NewDefaultDialog() *DefaultDialog {
	return &DefaultDialog{}
}

// Classify implements the Dialog interface
func (d *DefaultDialog) Classify(
	_ *UserMessage,
	_ context.Context) *ClassifyResult {
	return &ClassifyResult{
		Confidence:      ConfidenceHigh,
		ConfidenceValue: 1.0,
	}
}

// Reply implements the Dialog interface
func (d *DefaultDialog) Reply(
	_ *UserMessage,
	_ *ClassifyResult,
	_ context.Context) (*ReplyMessage, error) {
	replyText := fmt.Sprintf("%s\n%s\n%s", ReplyDefault, ReplyBar, ReplyUsage)
	return &ReplyMessage{
		Text:       replyText,
		Confidence: 1.0,
	}, nil

	return nil, errors.WithStack(ErrNotTriggered)
}
