package subscription

import (
	"net/http"
)

func HandleCheckoutSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		HandleCreateCheckoutSession(w, r)
		return
	case http.MethodGet:
		HandleGetCheckoutSessions(w, r)
		return
	default:
		http.NotFound(w, r)
		return
	}
}
