package subscription

import (
	"log"
	"net/http"

	"backendgo/internal/response"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
)

func HandleCancelSubscription(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SubscriptionID string `json:"subscriptionId"`
	}
	req.SubscriptionID = r.PathValue("id")

	s, err := sub.Cancel(req.SubscriptionID, nil)
	if err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("sub.Cancel: %v", err)
		return
	}

	response.WriteJSON(w, struct {
		Subscription *stripe.Subscription `json:"subscription"`
	}{
		Subscription: s,
	}, nil)
}
