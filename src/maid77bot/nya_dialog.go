package maid77bot

import (
	"context"
	"github.com/pkg/errors"
	"strings"
)

// NyaDialog is just a echo service
type NyaDialog struct {
}

// NewNyaDialog creates a new instance of NyaDialog
func NewNyaDialog() *NyaDialog {
	return &NyaDialog{}
}

// Classify implements the Dialog interface
func (d *NyaDialog) Classify(
	message *UserMessage,
	ctx context.Context) *ClassifyResult {
	fields := strings.Fields(message.Text)
	if len(fields) == 0 {
		return NewVeryLowConfidenceResult()
	}

	cmd := fields[0]
	if cmd == CmdNya {
		return &ClassifyResult{
			Confidence:      ConfidenceHigh,
			ConfidenceValue: 1.0,
		}
	} else {
		return NewVeryLowConfidenceResult()
	}
}

// Reply implements the Dialog interface
func (d *NyaDialog) Reply(
	message *UserMessage,
	_ *ClassifyResult,
	ctx context.Context) (*ReplyMessage, error) {
	if message.Text == CmdNya {
		return &ReplyMessage{
			Text:       ReplyNya,
			Confidence: 1.0,
		}, nil
	}

	return nil, errors.WithStack(ErrNotTriggered)
}
