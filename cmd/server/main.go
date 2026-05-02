package main

import (
	"log"
	"net/http"

	"replurk-go/internal/config"
	"replurk-go/internal/plurk"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	api := &plurk.API{
		ConsumerKey:    cfg.ConsumerKey,
		ConsumerSecret: cfg.ConsumerSecret,
		Token:          cfg.OAuthToken,
		TokenSecret:    cfg.OAuthTokenSecret,
	}

	http.HandleFunc("/run-cron", func(w http.ResponseWriter, _ *http.Request) {
		log.Println("Manually triggering Plurk Auto-Replurk Job...")
		plurks, err := api.SearchPlurks()
		if err != nil {
			writeErr(w, err)
			return
		}
		ids := make([]int, 0, len(plurks))
		for _, p := range plurks {
			ids = append(ids, p.PlurkID)
		}
		if len(plurks) == 0 {
			writeOK(w)
			return
		}
		if err := api.Replurk(ids); err != nil {
			writeErr(w, err)
			return
		}
		writeOK(w)
	})

	addr := ":" + cfg.Port
	log.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Tiny plaintext bodies: cron-job.org aborts jobs that read more than ~1 KiB of
// response data (see their FAQ). Avoid large HTML error pages from intermediaries
// by keeping success/error bodies minimal; full errors stay in server logs.
func writeOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func writeErr(w http.ResponseWriter, err error) {
	log.Printf("handler error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("ERR"))
}
