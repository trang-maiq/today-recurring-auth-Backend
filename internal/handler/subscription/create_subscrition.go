package subscription

import (
	"encoding/json"
	"log"
	"net/http"

	"backendgo/internal/response"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
)

func HandleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID string `json:"customerId"`
		PriceID    string `json:"priceId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	// Create subscription
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(req.CustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(req.PriceID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	s, err := sub.New(subscriptionParams)

	if err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("sub.New: %v", err)
		return
	}

	response.WriteJSON(w, struct {
		SubscriptionID string `json:"subscriptionId"`
		ClientSecret   string `json:"clientSecret"`
	}{
		SubscriptionID: s.ID,
		ClientSecret:   s.LatestInvoice.PaymentIntent.ClientSecret,
	}, nil)
}
