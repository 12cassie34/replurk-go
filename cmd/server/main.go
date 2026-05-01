package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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

func writeOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func writeErr(w http.ResponseWriter, err error) {
	log.Printf("handler error: %v", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":    "error",
		"message":   err.Error(),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
