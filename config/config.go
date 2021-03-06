package config

import (
	"errors"
	"os"
	"time"

	"github.com/satori/go.uuid"
)

// SaveLocation the location to save the config to.
var SaveLocation = "config.json"

// DefaultConfigSavedError an feedback returned if the default config is saved.
var DefaultConfigSavedError = errors.New("the default config has been saved, please edit it")

// Config the main configuration.
type Config struct {
	Data         `json:"-"`
	Discord      Discord      `json:"discord"`
	DiscordOAuth DiscordOAuth `json:"discord_oauth"`
	MySQL        MySQL        `json:"mysql"`
	Redis        Redis        `json:"redis"`
	Web          Web          `json:"web"`
	OwO          OwO          `json:"owo"`
}

// Redis configures redis.
type Redis struct {
	Network  string `json:"network"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Database int    `json:"database"`
	Enabled  bool   `json:"enabled"`
}

// MySQL configures MySQL.
type MySQL struct {
	DatabaseType string `json:"database_type"`
	URI          string `json:"uri"`
	Enabled      bool   `json:"enabled"`
}

// Discord configures the Discord bot.
type Discord struct {
	Token string `json:"token"`
}

// DscordOAuth configures oauth via discord
type DiscordOAuth struct {
	Key      string `json:"key"`
	Secret   string `json:"secret"`
	Callback string `json:"callback"`
}

// Web configures gin and other web elements
type Web struct {
	StaticFilePath   string   `json:"static_file_path"`
	ListenAddress    string   `json:"listen_address"`
	LogAuthKey       string   `json:"log_auth_key"`
	TemplateGlob     string   `json:"template_glob"`
	SentryDSN        string   `json:"sentry_dsn"`
	CSRFSecret       string   `json:"csrf_secret"`
	CSPReportWebHook string   `json:"csp_report_webhook"`
	DomainNames      []string `json:"domain_names"`
	AlexaAppID       string   `json:"alexa_app_id"`
	LogDirectory     string   `json:"log_directory"`
}

type OwO struct {
	UploadURL string        `json:"upload_url"`
	Timeout   time.Duration `json:"timeout"`
	URL       string        `json:"url"`
}

// DefaultConfig the default configuration to save.
var DefaultConfig = Config{
	Data:    Data{},
	Discord: Discord{"TOKEN"},
	MySQL: MySQL{
		DatabaseType: "mysql",
		URI:          "username:password@tcp(127.0.0.1:3306)/selfbot?charset=utf8&parseTime=True&loc=Local",
		Enabled:      false,
	},
	Redis: Redis{
		Enabled:  false,
		Database: 1,
		Address:  "127.0.0.1:6379",
		Network:  "tcp",
		Password: "password",
	},
	DiscordOAuth: DiscordOAuth{
		Callback: "https://sb.cory.red/",
		Key:      "key",
		Secret:   "secret",
	},
	Web: Web{
		LogDirectory:     "/var/log/selfbot",
		StaticFilePath:   "static/",
		ListenAddress:    "",
		LogAuthKey:       "memememememem",
		TemplateGlob:     "templates/**/*.tmpl",
		CSRFSecret:       uuid.Must(uuid.NewV4()).String() + "-ChangePls",
		CSPReportWebHook: "https://discordapp.com/webhook/slack",
		AlexaAppID:       "amzn1.ask.skill.UUIDHERE",
		DomainNames:      []string{"sb.cory.red"},
	},
	OwO: OwO{
		UploadURL: "https://api.awau.moe/upload/?key=uuid",
		Timeout:   time.Second * 10,
		URL:       "https://awau.moe/",
	},
}

// Save saves the config.
func (c *Config) Save() error {
	saveLoc, envThere := os.LookupEnv("CONFIG_LOC")
	if !envThere {
		saveLoc = SaveLocation
	}

	return c.save(saveLoc, c)
}

// Load loads the config.
func (c *Config) Load() error {

	saveLoc, envThere := os.LookupEnv("CONFIG_LOC")
	if !envThere {
		saveLoc = SaveLocation
	}

	if err := c.load(saveLoc, c); err == DefaultConfigSavedError {
		if err := DefaultConfig.Save(); err != nil {
			return err
		}
		return DefaultConfigSavedError
	} else if err != nil {
		return err
	}

	return nil
}
