package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/meinside/free-epic-games-notifier/database"
	"github.com/meinside/free-epic-games-notifier/extractor"
	"github.com/meinside/free-epic-games-notifier/notifier"
)

const (
	defaultConfigFilename = "epic_notifier.json"
	defaultCacheFilename  = "caches.db"

	timeoutSeconds = 60
)

type conf struct {
	JandiWebhookURL       string `json:"jandi_webhook_url,omitempty"`
	PushbulletAccessToken string `json:"pushbullet_access_token,omitempty"`
}

var _confPath string
var _cachePath string

var _conf conf
var _notifiers []notifier.Notifier

func init() {
	flag.StringVar(&_confPath, "conf", defaultConfigFilename, "config filepath")
	flag.StringVar(&_cachePath, "cache", defaultCacheFilename, "cache filepath")
	flag.Parse()

	loadConf()
}

func loadConf() {
	var err error
	var bytes []byte
	if bytes, err = ioutil.ReadFile(_confPath); err == nil {
		if err = json.Unmarshal(bytes, &_conf); err == nil {
			_notifiers = []notifier.Notifier{}

			if _conf.JandiWebhookURL != "" {
				_notifiers = append(_notifiers, notifier.JandiNotifier{WebhookURL: _conf.JandiWebhookURL})
			}

			if _conf.PushbulletAccessToken != "" {
				_notifiers = append(_notifiers, notifier.PushbulletNotifier{Token: _conf.PushbulletAccessToken})
			}
		} else {
			log.Printf("failed to load config file: %s", err)
		}
	} else {
		log.Printf("failed to open config file: %s", err)
	}
}

// notify game
func notifyGame(game extractor.FreeGame) bool {
	sent := false

	for _, notifier := range _notifiers {
		if err := notifier.Notify(game); err == nil {
			sent = true
		} else {
			log.Printf("failed to send notification: %s", err)
		}
	}

	return sent
}

// notify error
func notifyError(err error) {
	for _, notifier := range _notifiers {
		if err := notifier.NotifyError(err); err != nil {
			log.Printf("failed to send error notification: %s", err)
		}
	}
}

func main() {
	quit := make(chan bool, 1)

	go func(quit chan bool) {
		if games, err := extractor.ExtractFreeGames(); err != nil {
			log.Printf("failed to extract free games: %s", err)

			notifyError(err)
		} else {
			if db, err := database.Open(_cachePath); err != nil {
				log.Printf("failed to open cache database: %s", err)
			} else {
				defer db.Close()

				sentNotification := false

				for _, game := range games {
					if cached, err := db.IsCachedGame(game.Title); !cached {
						if notifyGame(game) {
							sentNotification = true
						}

						if err := db.CacheGame(game); err != nil {
							log.Printf("failed to cache free game: %s", err)
						}
					} else if err != nil {
						log.Printf("failed to check if it is cached: %s", err)
					}
				}

				if !sentNotification {
					log.Printf("no new games notified")
				}
			}
		}

		quit <- true
	}(quit)

	// wait...
	select {
	case <-quit:
		log.Printf("operation finished")
		break
	case <-time.After(timeoutSeconds * time.Second):
		err := fmt.Errorf("operation timed out: %d seconds", timeoutSeconds)
		log.Printf("error: %s", err)
		notifyError(err)
		break
	}
}
