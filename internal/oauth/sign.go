package oauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

// EncodeURIComponent mirrors JavaScript encodeURIComponent for header values (RFC 3896-ish; unreserved excludes more than Go's URL helpers).
func EncodeURIComponent(s string) string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' || r == '!' || r == '~' || r == '*' || r == '\'' || r == '(' || r == ')' {
			b.WriteRune(r)
		} else {
			buf := make([]byte, 4)
			n := utf8.EncodeRune(buf, r)
			for _, c := range buf[:n] {
				fmt.Fprintf(&b, "%%%02X", c)
			}
		}
	}
	return b.String()
}

// PercentEncode matches the Node app: encodeURIComponent(s).replace(/[!'()*]/g, escape)
func PercentEncode(s string) string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' || r == '!' || r == '~' || r == '*' || r == '\'' || r == '(' || r == ')' {
			b.WriteRune(r)
		} else {
			buf := make([]byte, 4)
			n := utf8.EncodeRune(buf, r)
			for _, c := range buf[:n] {
				fmt.Fprintf(&b, "%%%02X", c)
			}
		}
	}
	out := b.String()
	repl := strings.NewReplacer(
		"!", "%21",
		"'", "%27",
		"(", "%28",
		")", "%29",
		"*", "%2A",
	)
	return repl.Replace(out)
}

// SignGET builds OAuth 1.0a HMAC-SHA1 signature for parameters (all string values).
func SignGET(method, requestURL string, params map[string]string, consumerSecret, tokenSecret string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, PercentEncode(k)+"="+PercentEncode(params[k]))
	}
	sorted := strings.Join(pairs, "&")
	baseString := strings.ToUpper(method) + "&" + PercentEncode(requestURL) + "&" + PercentEncode(sorted)
	signingKey := PercentEncode(consumerSecret) + "&" + PercentEncode(tokenSecret)
	mac := hmac.New(sha1.New, []byte(signingKey))
	mac.Write([]byte(baseString))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func AuthHeader(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		v := params[k]
		parts = append(parts, k+"=\""+EncodeURIComponent(v)+"\"")
	}
	return "OAuth " + strings.Join(parts, ", ")
}
