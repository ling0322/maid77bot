package maid77bot

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// Dialog to choose books of wordlist
type ChooseBookDialog struct {
	// Map from language to it's display name
	languages map[string]string
}

// NewChooseBookDialog creates a new instance of ChooseBookDialog
func NewChooseBookDialog(c *Config) *ChooseBookDialog {
	languages := map[string]string{}
	for lang, langBook := range c.WordBooks {
		languages[lang] = langBook.DisplayName
	}

	return &ChooseBookDialog{
		languages: languages,
	}
}

// Classify implements the Dialog interface
func (d *ChooseBookDialog) Classify(
	message *UserMessage,
	ctx context.Context) *ClassifyResult {
	fields := strings.Fields(message.Text)
	if len(fields) == 0 {
		return NewVeryLowConfidenceResult()
	}

	cmd := fields[0]
	if cmd == CmdStartWords || cmd == CmdChooseBook {
		return &ClassifyResult{
			Confidence:      ConfidenceHigh,
			ConfidenceValue: 1.0,
		}
	}

	return NewVeryLowConfidenceResult()
}

// Reply implements the Dialog interface
func (d *ChooseBookDialog) Reply(
	message *UserMessage,
	_ *ClassifyResult,
	ctx context.Context) (*ReplyMessage, error) {
	if message.Text == CmdStartWords || message.Text == CmdChooseBook {
		// Choose language
		buttons := []*MessageButton{}
		for language, displayName := range d.languages {
			buttons = append(buttons, &MessageButton{
				DisplayText:  displayName,
				CallbackText: fmt.Sprintf("%s %s", CmdChooseBook, language),
			})
		}

		return &ReplyMessage{
			Text:       ReplyChooseLang,
			Buttons:    buttons,
			Confidence: 1.0,
		}, nil
	}

	return nil, errors.WithStack(ErrNotTriggered)
}
