package maid77bot

import (
	"errors"
)

// Name of each commands for bot
const (
	CmdStartWords = "/start_words"
	CmdChooseBook = "/choose_book"
	CmdNya        = "/nya"
)

// Reply text and template of reply text
const (
	ReplyChooseLang = "想选择什么语言的呢?"
	ReplyNya        = "喵喵喵喵喵喵喵! ~77现在的版本是v1"
	ReplyDefault    = "77不知道你在说什么呢（现在只是半成品阶段喵）"
	ReplyBar        = "-------------"
)

// Name for each dialog modules
const (
	DialogModuleChooseBook = "choose_book"
	DialogModuleNya        = "nya"
	DialogModuleDefault    = "default"
)

var ErrNotTriggered = errors.New("dialog is not triggered")
var ErrTimeout = errors.New("timeout")
var ErrNotMatched = errors.New("not matched")

const ReplyUsage = `女仆77的使用方法:
/nya: 喵喵喵!`
