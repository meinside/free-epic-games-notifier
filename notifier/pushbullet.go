package notifier

import (
	"fmt"

	"github.com/meinside/free-epic-games-notifier/extractor"

	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
)

const (
	pushbulletTitle = "Free Epic Games Notification"
)

// PushbulletNotifier struct
type PushbulletNotifier struct {
	Token string
}

// Notify notifies a new free game
func (n PushbulletNotifier) Notify(game extractor.FreeGame) (err error) {
	client := pushbullet.New(n.Token)

	note := requests.NewNote()
	note.Title = pushbulletTitle
	note.Body = fmt.Sprintf("Free Now: %s\n%s", game.Title, game.StoreURL)

	_, err = client.PostPushesNote(note)

	return err
}

// NotifyError notifies an error
func (n PushbulletNotifier) NotifyError(e error) (err error) {
	client := pushbullet.New(n.Token)

	note := requests.NewNote()
	note.Title = pushbulletTitle
	note.Body = fmt.Sprintf("Error: %s", e)

	_, err = client.PostPushesNote(note)

	return err
}
