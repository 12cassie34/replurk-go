package plurk

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"time"

	"replurk-go/internal/httpclient"
	"replurk-go/internal/oauth"
)

type API struct {
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	TokenSecret    string
}

type searchResp struct {
	Plurks []Plurk `json:"plurks"`
}

func cloneMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// SearchCronQuery is the hashtag search string used by /run-cron (index.ts searchPlurks).
const SearchCronQuery = "#卿卿我我 #恩恩愛愛 -J·M"

// SearchPlurks returns matching plurks that are not yet replurked.
func (c *API) SearchPlurks() ([]Plurk, error) {
	q := SearchCronQuery
	p := oauthBase(c.ConsumerKey, c.Token)
	p["query"] = q
	p["offset"] = "0"
	sig := oauth.SignGET(http.MethodGet, searchURL, p, c.ConsumerSecret, c.TokenSecret)
	all := cloneMap(p)
	all["oauth_signature"] = sig
	u := searchURL + "?query=" + oauth.EncodeURIComponent(q) + "&offset=0"

	var data searchResp
	if err := httpclient.GetJSON(u, map[string]string{"Authorization": oauth.AuthHeader(all)}, &data); err != nil {
		return nil, err
	}
	if len(data.Plurks) == 0 {
		log.Println("No new Plurks found.")
		return nil, nil
	}
	var out []Plurk
	for _, pl := range data.Plurks {
		if !pl.Replurked {
			out = append(out, pl)
		}
	}
	return out, nil
}

// Replurk GET /APP/Timeline/replurk with retry (index.ts replurk).
func (c *API) Replurk(plurkIDs []int) error {
	if len(plurkIDs) == 0 {
		return nil
	}
	idsJSON := jsonPlurkIDs(plurkIDs)
	maxRetries := 3

	for retry := 0; retry < maxRetries; retry++ {
		p := oauthBase(c.ConsumerKey, c.Token)
		p["ids"] = idsJSON
		sig := oauth.SignGET(http.MethodGet, replurkURL, p, c.ConsumerSecret, c.TokenSecret)
		all := cloneMap(p)
		all["oauth_signature"] = sig

		uv := url.Values{}
		for _, k := range sortedKeys(all) {
			uv.Set(k, all[k])
		}
		fullURL := replurkURL + "?" + uv.Encode()

		var raw json.RawMessage
		err := httpclient.GetJSON(fullURL, map[string]string{"Authorization": oauth.AuthHeader(all)}, &raw)
		if err == nil {
			log.Printf("Replurk response: %s", string(raw))
			return nil
		}
		msg := err.Error()
		log.Printf("Error replurking (attempt %d/%d): %s", retry+1, maxRetries, msg)
		if retry == maxRetries-1 {
			return fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
		}
		retryCount := retry + 1
		delay := time.Millisecond * time.Duration(minInt(1000*(1<<(retryCount-1)), 10000))
		log.Printf("Waiting %v before retry...", delay)
		time.Sleep(delay)
	}
	return nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
