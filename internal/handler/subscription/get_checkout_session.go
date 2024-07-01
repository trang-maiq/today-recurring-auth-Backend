package subscription

import (
	"log/slog"
	"net/http"

	"github.com/stripe/stripe-go/v72/checkout/session"

	"backendgo/internal/response"
)

func HandleGetCheckoutSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	sessionID := r.URL.Query().Get("sessionId")
	s, err := session.Get(sessionID, nil)
	if err != nil {
		slog.Error("get checkout session failed",
			slog.Any("error", err),
			slog.String("session_id", sessionID),
		)
		response.WriteJSON(w, s, err)
	}
	response.WriteJSON(w, s, nil)
}
