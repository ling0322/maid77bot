package maid77bot

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// ConfigBook is the config for specific book
type ConfigBook struct {
	// Name of this book
	DisplayName string `yaml:"display_name"`

	// Path of this book
	Path string `yaml:"path"`
}

// ConfigBookLang is the language and book config for WordBooks field in config
type ConfigBookLang struct {
	// Name of this language
	DisplayName string `yaml:"display_name"`

	// Books for this language (bookName -> )
	Books map[string]*ConfigBook
}

// Config for the bot
type Config struct {
	// Config for each language and book. The yaml structure would be like
	//   word_books:
	//     en:
	//       display_name: English
	//       books:
	//         cet4:
	//           display_name: CET-4
	//           path: cet-4.txt
	//         gre:
	//           display_name: GRE
	//           path: gre.txt
	WordBooks map[string]*ConfigBookLang `yaml:"word_books"`

	// Priority of each dialog module. Modules that didn't exist in this map is
	// disabled.
	DialogPriority map[string]int `yaml:"dialog_priority"`

	// Socks5Proxy is the address:port of socks5 proxy for this bot, for example
	// localhost:1080. If Socks5Proxy is empty, use the default
	Socks5Proxy string `yaml:"socks5_proxy"`

	// Token for this bot in Telegram
	TelegramToken string `yaml:"telegram_token"`
}

// LoadConfig loads yaml format config file and returns as a Config instance
func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := Config{}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &c, nil
}
