package notifier

import (
	"fmt"

	"github.com/meinside/free-epic-games-notifier/extractor"

	"github.com/meinside/jandi-webhook-go"
)

const (
	incomingWebhookTitle       = "Free Epic Games"
	incomingWebhookDescription = "Notification from [Free Epic Games Notifier](https://github.com/meinside/free-epic-games-notifier)"
	incomingWebhookColor       = "#0000FF"
	incomingWebhookErrorColor  = "#FF0000"
)

// JandiNotifier struct
type JandiNotifier struct {
	WebhookURL string
}

// Notify notifies a new free game
func (n JandiNotifier) Notify(game extractor.FreeGame) (err error) {
	client := jandi.NewIncomingClient(n.WebhookURL)

	_, err = client.SendIncoming(
		incomingWebhookTitle,
		incomingWebhookColor,
		[]jandi.ConnectInfo{
			jandi.ConnectInfoFrom(fmt.Sprintf("Free Now: %s", game.Title), game.StoreURL, game.ImageURL),
		},
	)

	return err
}

// NotifyError notifies an error
func (n JandiNotifier) NotifyError(e error) (err error) {
	client := jandi.NewIncomingClient(n.WebhookURL)

	_, err = client.SendIncomingWithTitle(
		incomingWebhookTitle,
		incomingWebhookDescription,
		incomingWebhookErrorColor,
		[]jandi.ConnectInfo{
			jandi.ConnectInfoFrom(fmt.Sprintf("Error: %s", e), "", ""),
		},
	)

	return err
}
