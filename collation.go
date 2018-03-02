package quotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	slackAPIURL                 = "https://slack.com/api/"
	slackChannelHistoryEndpoint = "channels.history?token=%s&channel=%s"
	slackUsersEndpoint          = "users.list?token=%s"
)

var errDisplay = errors.New("failed :/ sorry, not sorry")

// GenerateSlackQuoteDisplay generates display
func GenerateSlackQuoteDisplay(channelID, slackToken string) (string, error) {
	// Retrieve author names
	n, err := GetSlackNames(slackToken)
	if err != nil {
		return "", errDisplay
	}

	// Retrieve latest quotes
	q, err := GetSlackQuotes(n, channelID, slackToken)
	if err != nil {
		return "", errDisplay
	}

	// Shuffle results
	rand.Seed(time.Now().UnixNano())
	Shuffle(q)

	// Take 30 or all if less than 30
	result := []Quote{}
	if len(q) >= 30 {
		result = q[0:30]
	} else {
		result = q
	}

	// Create front-end output
	j, err := json.Marshal(result)
	if err != nil {
		return "", errDisplay
	}

	t, err := getTemplate(string(j))
	if err != nil {
		return "", errDisplay
	}
	return t, nil
}

// Shuffle a slice
func Shuffle(vals []Quote) []Quote {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]Quote, len(vals))
	n := len(vals)
	for i := 0; i < n; i++ {
		randIndex := r.Intn(len(vals))
		ret[i] = vals[randIndex]
		vals = append(vals[:randIndex], vals[randIndex+1:]...)
	}
	return ret
}

// GetSlackQuotes ; pass in a list of possible authors with the 'quote' channel ID and this will return complete quotes
// will only return quotes of format `@bob "i am bob and i totally said this"`
func GetSlackQuotes(authors map[string]string, channelID, token string) ([]Quote, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(slackAPIURL+slackChannelHistoryEndpoint, token, channelID), nil)
	if err != nil {
		return []Quote{}, errors.New("Failed creating request to get channel")
	}

	cli := http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return []Quote{}, errors.New("Failed request to get channel")
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	var rb ChannelResp
	json.Unmarshal(b, &rb)
	result := []Quote{}

	// Create quotes
	for _, m := range rb.Messages {
		s := strings.Split(m.Text, " ")

		if len(s) <= 0 {
			continue
		}

		id := s[0]

		rn := authors[id]
		if rn == "" {
			continue
		}

		q := strings.Replace(m.Text, id, "", -1)
		q = strings.Replace(q, "“", "\"", -1)
		q = strings.Replace(q, "”", "\"", -1)
		q = strings.Replace(q, "’", "'", -1)

		// Ingore when links / images are posted
		if strings.Contains(q, "<http") {
			continue
		}

		// Replaces names in quotes
		if strings.Contains(q, "<@") {
			for ida, naa := range authors {
				q = strings.Replace(q, ida, naa, -1)
			}
		}

		qb64 := base64.StdEncoding.EncodeToString([]byte(html.EscapeString(q)))
		result = append(result, Quote{
			Speaker: rn,
			Speech:  qb64,
		})
	}

	return result, nil
}

// GetSlackNames returns a map of map[Slack_ID]Human_Name
func GetSlackNames(slackToken string) (map[string]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(slackAPIURL+slackUsersEndpoint, slackToken), nil)
	if err != nil {
		return map[string]string{}, errors.New("Failed creating request to get users")
	}

	cli := http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return map[string]string{}, errors.New("Failed request to get users")
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	var rb UsersResp
	json.Unmarshal(b, &rb)

	// Map IDs to names with internal ID structure
	for _, i := range rb.Members {
		result["<@"+i.ID+">"] = i.Name
	}

	return result, nil
}

func getTemplate(arr string) (string, error) {
	s, err := Asset("resources/template.html")
	if err != nil {
		return "", err
	}
	return strings.Replace(string(s), "[{\"speaker\": \"Realman Notfakeson\", \"speech\": \"IlRoaXMgaXMgYSByZWFsbHkgcmVhbCBxdW90ZSI=\"}]", arr, 1), nil
}
