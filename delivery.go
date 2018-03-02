package quotify

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
)

const (
	displayName = "'Quote Fortifi' Daily Inspiration BOT"
	displayIcon = ":party_parrot:"
)

// DeliverInspiration adds new inspiration to team slack channel
// it also preloads the next image in the chain ready for
// consecutive execution of this process
func DeliverInspiration(slackWebhook, generateEndpoint, slackChannel string) {
	now, next := getURLS(generateEndpoint)

	// Make sure current image is in cache
	http.Get(now)
	fmt.Println(time.Now().UTC(), "pre-loaded current inspiration")

	// Preload the 'next' image so that it is cached in URL grab or refreshes
	// old cache using the same entropy;
	http.Get(next)
	fmt.Println(time.Now().UTC(), "pre-loaded next inspiration")

	// Post the 'now' image to channel
	payload := slack.Payload{
		Username:  displayName,
		Channel:   slackChannel,
		IconEmoji: displayIcon,
		Attachments: []slack.Attachment{
			slack.Attachment{
				ImageUrl: &now,
			},
		},
	}
	slack.Send(slackWebhook, "", payload)
	fmt.Println(time.Now().UTC(), "posted team inspiration")
}

// Deterministically retrieve current and next image url using entropy provided
func getURLS(generateEndpoint string) (current, next string) {
	// random seed used for image on each day
	current = generateEndpoint
	next = generateEndpoint
	now := time.Now()
	morrow := now.AddDate(0, 0, 1)
	entToday := fmt.Sprintf("%d%d%d", now.Year(), now.Month(), now.Day())
	entTomorrow := fmt.Sprintf("%d%d%d", morrow.Year(), morrow.Month(), morrow.Day())
	currentSuff, nextSuff := getSuffixForPeriods()

	// Current url with current rotation code
	current = current + entToday + currentSuff

	if nextSuff != "" {
		// Next url with current rotation code + suffix
		// (image was generated for the same day as the prior image)
		next = current + entToday + nextSuff
	} else {
		// Next url with tomorrow's rotation code
		next = next + entTomorrow
	}

	return
}

// getSuffixForPeriods gets the link suffix for
// current period of the day (since two images per day)
// and the period following the current one
func getSuffixForPeriods() (now, next string) {
	const (
		amSuffix   = ""
		pmSuffix   = "B"
		stdPMTempl = "PM"
	)
	now = time.Now().UTC().Format(stdPMTempl)
	if now == stdPMTempl {
		now = pmSuffix
		next = amSuffix
		return
	}

	now = amSuffix
	next = pmSuffix
	return
}
