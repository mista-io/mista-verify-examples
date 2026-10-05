// Sign-up / login verification API with net/http.
//
//	export MISTA_API_TOKEN="your_mista_api_token_here"
//	go run ./cmd/server
//
//	curl -X POST localhost:3000/api/verify/start -H 'Content-Type: application/json' \
//	     -d '{"phone": "+15555550100", "channel": "whatsapp_sms"}'
//	curl -X POST localhost:3000/api/verify/check -H 'Content-Type: application/json' \
//	     -d '{"sid": "SID_FROM_START", "code": "123456"}'
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/mista-io/mista-go"
)

var reasonMessages = map[string]string{
	"invalid_code": "That code is incorrect. Try again.",
	"expired":      "This code has expired. Request a new one.",
	"max_attempts": "Too many attempts. Request a new code.",
	"not_open":     "This verification is no longer active. Request a new code.",
}

type server struct {
	mista *mista.Client

	// In production, keep this in your session store or database.
	mu      sync.Mutex
	pending map[string]string // sid -> phone
}

func main() {
	if os.Getenv("MISTA_API_TOKEN") == "" {
		log.Fatal("Set MISTA_API_TOKEN first (Mista dashboard -> Settings -> API).")
	}
	s := &server{mista: mista.NewClient(""), pending: map[string]string{}}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/verify/start", s.start)
	mux.HandleFunc("POST /api/verify/check", s.check)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func (s *server) start(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Phone   string `json:"phone"`
		Channel string `json:"channel"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || !strings.HasPrefix(body.Phone, "+") {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "phone must be in E.164 format, e.g. +15555550100"})
		return
	}
	if body.Channel == "" {
		body.Channel = "auto"
	}

	verification, err := s.mista.Verify.Start(r.Context(), &mista.StartVerificationParams{To: body.Phone, Channel: body.Channel})
	if err != nil {
		writeMistaError(w, err)
		return
	}

	s.mu.Lock()
	s.pending[verification.SID] = body.Phone
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, map[string]any{
		"sid":        verification.SID,
		"channel":    verification.Channel,
		"expires_at": verification.ExpiresAt,
	})
}

func (s *server) check(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SID  string `json:"sid"`
		Code string `json:"code"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	s.mu.Lock()
	phone, ok := s.pending[body.SID]
	s.mu.Unlock()
	if !ok || body.Code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Unknown verification. Request a new code."})
		return
	}

	result, err := s.mista.Verify.Check(r.Context(), body.SID, strings.TrimSpace(body.Code))
	if err != nil {
		writeMistaError(w, err)
		return
	}
	if !result.Verified {
		message, known := reasonMessages[result.Reason]
		if !known {
			message = "Verification failed. Request a new code."
		}
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"verified": false, "reason": result.Reason, "error": message})
		return
	}

	s.mu.Lock()
	delete(s.pending, body.SID)
	s.mu.Unlock()

	// The phone number is now proven. Create the account / start the session here.
	writeJSON(w, http.StatusOK, map[string]any{"verified": true, "phone": phone})
}

func writeMistaError(w http.ResponseWriter, err error) {
	var apiErr *mista.APIError
	switch {
	case errors.Is(err, mista.ErrRateLimited) && errors.As(err, &apiErr):
		writeJSON(w, http.StatusTooManyRequests, map[string]any{"error": "Too many requests. Try again shortly.", "retry_after": apiErr.RetryAfter()})
	case errors.Is(err, mista.ErrValidation), errors.Is(err, mista.ErrBadRequest), errors.Is(err, mista.ErrInvalidParams):
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
	case errors.Is(err, mista.ErrUnauthorized):
		log.Print("Mista rejected the API token. Check MISTA_API_TOKEN.")
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": "Verification is temporarily unavailable."})
	default:
		log.Print(err)
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": "Verification is temporarily unavailable."})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
