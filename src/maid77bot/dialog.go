package maid77bot

import (
	"context"
	"github.com/pkg/errors"
	"log"
)

const (
	ConfidenceVeryHigh = 5
	ConfidenceHigh     = 4
	ConfidenceMedium   = 3
	ConfidenceLow      = 2
	ConfidenceVeryLow  = 1
)

// ClassifyResult is the result of Classification in Dialog
type ClassifyResult struct {
	// The level of confidence, like ConfidenceVeryHigh, ConfidenceLow, ...
	Confidence int

	// Score of the confidence
	ConfidenceValue float64

	// Some intermediate results from classifier
	Bundle map[string]interface{}
}

// Dialog is the interface to classify and reply user messages
type Dialog interface {
	// Classify predicts the confidence that the message belongs to this Dialog
	Classify(*UserMessage, context.Context) *ClassifyResult

	// Answer updates states and generates answers for the given messageg
	Reply(*UserMessage, *ClassifyResult, context.Context) (*ReplyMessage, error)
}

// Controller of each dialog module. It will create instances of each dialog
// according to the config. When user message arrives, it will call Classify()
// method of each dialog instance to choose the most related one. Then get the
// reply from it
type DialogController struct {
	botConfig  *Config
	dialogs    map[string]Dialog
	priorities map[string]int
}

// NewVeryLowConfidenceResult creates a ClassifyResult with very low confidence
func NewVeryLowConfidenceResult() *ClassifyResult {
	return &ClassifyResult{
		Confidence:      ConfidenceVeryLow,
		ConfidenceValue: 0.0,
	}
}

// NewDialogController creates a new instance of DialogController
func NewDialogController(config *Config) (*DialogController, error) {
	priorities := map[string]int{}
	for name, prior := range config.DialogPriority {
		priorities[name] = prior
	}

	// Create instance of Dialog according to priorities
	dialogs := map[string]Dialog{}
	for name := range priorities {
		dialog, err := createDialogByName(name, config)
		if err != nil {
			return nil, err
		}

		dialogs[name] = dialog
	}

	return &DialogController{
		botConfig:  config,
		dialogs:    dialogs,
		priorities: priorities,
	}, nil
}

// createDialogByName create dialog by its name
func createDialogByName(name string, botConfig *Config) (Dialog, error) {
	switch name {
	case DialogModuleChooseBook:
		return NewChooseBookDialog(botConfig), nil
	case DialogModuleNya:
		return NewNyaDialog(), nil
	case DialogModuleDefault:
		return NewDefaultDialog(), nil
	default:
		return nil, errors.Errorf("unexpected dialog module: %s", name)
	}
}

// classify calls the Classify method of all dialog instances. And return the
// result as a map from dialog name to *ClassifyResult
func (c *DialogController) classify(
	m *UserMessage,
	ctx context.Context) (map[string]*ClassifyResult, error) {
	type nameResultT struct {
		name   string
		result *ClassifyResult
	}
	resultCh := make(chan nameResultT, 100)

	// Run classifiers from each dialog
	for dialogName, dialog := range c.dialogs {
		go func(dialogName string, dialog Dialog) {
			log.Println("classify: ", dialogName, " start")
			r := dialog.Classify(m, ctx)
			resultCh <- nameResultT{dialogName, r}
		}(dialogName, dialog)
	}

	// Get the answers each classifier
	classifyResults := map[string]*ClassifyResult{}
	for _ = range c.dialogs {
		select {
		case nr := <-resultCh:
			log.Println("classify: ", nr.name, " finished")
			classifyResults[nr.name] = nr.result
		case <-ctx.Done():
			return nil, errors.WithStack(ErrTimeout)
		}
	}

	return classifyResults, nil
}

// rank choose the best dialog from classify results and return its name
func (c *DialogController) rankDialog(res map[string]*ClassifyResult) string {
	maxScore := 0.0
	maxDialog := ""
	for name, r := range res {
		priority := float64(c.priorities[name])
		score := 100000*float64(r.Confidence) +
			10*(1000-float64(priority)) +
			r.ConfidenceValue
		log.Printf(
			"%s: priority = %d, c = %d, cv = %.3f, s = %.3f",
			name,
			c.priorities[name],
			r.Confidence,
			r.ConfidenceValue,
			score)
		if score >= maxScore {
			maxScore = score
			maxDialog = name
		}
	}

	return maxDialog
}

// Reply replies the user message using the most confident dialog it belongs to.
// And updates session state for this dialog
func (c *DialogController) Reply(
	m *UserMessage,
	ctx context.Context) (*ReplyMessage, error) {
	result, err := c.classify(m, ctx)
	if err != nil {
		return nil, err
	}
	log.Println("result = ", result)

	bestName := c.rankDialog(result)
	bestDialog, ok := c.dialogs[bestName]
	log.Println("bestName = ", bestName)

	if !ok {
		return nil, errors.WithStack(ErrNotMatched)
	}
	bestResult := result[bestName]

	return bestDialog.Reply(m, bestResult, ctx)
}
