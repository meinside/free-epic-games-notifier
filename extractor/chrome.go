package extractor

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const (
	// urls
	epicURL  = "https://www.epicgames.com"
	storeURL = epicURL + "/store/en-US/free-games"

	// status
	statusFreeNow = "Free Now"

	// xpaths and css selectors
	xpathStatusFreeNow = `//span[contains(text(), 'Free Now')]`
	selectorStatus     = `div[data-component=StatusBar] > span[data-component=Message]`
	selectorTitle      = `div > span[data-testid=offer-title-info-title][data-component=OfferTitleInfo]`
	selectorLink       = `div[data-component=CardGridDesktopBase] > div > div > a`
	selectorImage      = `div[data-component=OfferCardImagePortrait] > div[data-component=Picture] > img`

	// for debugging
	_debug = false
	//_debug = true
)

// FreeGame is a struct for a free game
type FreeGame struct {
	Title    string
	StoreURL string
	ImageURL string
}

// ExtractFreeGames extracts free games from the store url
func ExtractFreeGames() ([]FreeGame, error) {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), chromedp.Headless, chromedp.NoSandbox, chromedp.DisableGPU)
	defer cancel()

	runCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// extract nodes from the store url
	var titles []*cdp.Node
	var statuses []*cdp.Node
	var links []*cdp.Node
	var imgs []*cdp.Node
	err := chromedp.Run(runCtx,
		chromedp.Navigate(storeURL),
		chromedp.WaitVisible(xpathStatusFreeNow, chromedp.BySearch),
		chromedp.Nodes(selectorTitle, &titles, chromedp.ByQueryAll),
		chromedp.Nodes(selectorStatus, &statuses, chromedp.ByQueryAll),
		chromedp.Nodes(selectorLink, &links, chromedp.ByQueryAll),
		chromedp.Nodes(selectorImage, &imgs, chromedp.ByQueryAll),
	)
	if err != nil {
		return nil, err
	}

	return filterFreeGames(titles, statuses, links, imgs)
}

// filter free games from given nodes
func filterFreeGames(titles, statuses, links, imgs []*cdp.Node) (games []FreeGame, err error) {
	games = []FreeGame{}

	// extract values
	gameTitles := []string{}
	for _, title := range titles {
		if len(title.Children) > 0 {
			gameTitles = append(gameTitles, title.Children[0].NodeValue)
		}
	}
	gameURLs := []string{}
	for _, link := range links {
		gameURLs = append(gameURLs, epicURL+link.AttributeValue("href"))
	}
	gameImageURLs := []string{}
	for _, img := range imgs {
		gameImageURLs = append(gameImageURLs, img.AttributeValue("src"))
	}

	if _debug {
		log.Printf("[debug] game titles: %+v", gameTitles)
		log.Printf("[debug] game links: %+v", gameURLs)
		log.Printf("[debug] game images: %+v", gameImageURLs)
	}

	// check game statuses
	for i, s := range statuses {
		if s.ChildNodeCount > 0 {
			status := s.Children[0].NodeValue

			if _debug {
				log.Printf("[debug] game status[%d]: %s", i, status)
			}

			// filter 'Free Now'
			if strings.Contains(status, statusFreeNow) {
				if len(gameTitles) > i && len(gameURLs) > i && len(gameImageURLs) > i {
					games = append(games, FreeGame{
						Title:    gameTitles[i],
						StoreURL: gameURLs[i],
						ImageURL: gameImageURLs[i],
					})
				} else {
					err = fmt.Errorf("elements' counts do not match")
				}
			}
		}
	}

	if len(games) <= 0 && err == nil {
		err = fmt.Errorf("there is no free game in the store page (%s)", storeURL)
	}

	return games, err
}
