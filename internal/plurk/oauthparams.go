package plurk

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"
)

func oauthNonce() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(b)
}

func oauthUnix() int64 {
	return time.Now().Unix()
}

func oauthBase(consumerKey, token string) map[string]string {
	return map[string]string{
		"oauth_consumer_key":     consumerKey,
		"oauth_token":           token,
		"oauth_nonce":           oauthNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":       strconv.FormatInt(oauthUnix(), 10),
		"oauth_version":         "1.0",
	}
}

func jsonPlurkIDs(ids []int) string {
	b, _ := json.Marshal(ids)
	return string(b)
}

