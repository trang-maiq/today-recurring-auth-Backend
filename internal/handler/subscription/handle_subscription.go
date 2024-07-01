package subscription

import "net/http"

func HandleSubscriptions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		HandleUpdateSubscription(w, r)
		return
	case http.MethodGet:
		HandleListSubscriptions(w, r)
		return
	case http.MethodPost:
		HandleCreateSubscription(w, r)
		return
	default:
		http.NotFound(w, r)
		return
	}
}

func HandleSubscription(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		HandleUpdateSubscription(w, r)
		return
	case http.MethodDelete:
		HandleCancelSubscription(w, r)
		return
	default:
		http.NotFound(w, r)
		return
	}
}
