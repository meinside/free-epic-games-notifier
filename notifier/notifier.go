package notifier

import (
	"github.com/meinside/free-epic-games-notifier/extractor"
)

// Notifier interface
type Notifier interface {
	Notify(game extractor.FreeGame) error
}
