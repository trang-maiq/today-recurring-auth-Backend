package subscription

import (
	"backendgo/internal/response"
	"encoding/json"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
	"log"
	"net/http"
	"os"
	"strings"
)

func HandleUpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SubscriptionID    string `json:"-"`
		NewPriceLookupKey string `json:"newPriceLookupKey"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	req.SubscriptionID = r.PathValue("id")

	// This is the ID of the Stripe Price object to which the subscription
	// will be upgraded or downgraded.
	newPriceID := os.Getenv(strings.ToUpper(req.NewPriceLookupKey))

	// Fetch the subscription to access the related subscription item's ID
	// that will be updated. In practice, you might want to store the
	// Subscription Item ID in your database to avoid this API call.
	s, err := sub.Get(req.SubscriptionID, nil)
	if err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("sub.Get: %v", err)
		return
	}

	params := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{{
			ID:    stripe.String(s.Items.Data[0].ID),
			Price: stripe.String(newPriceID),
		}},
	}

	updatedSubscription, err := sub.Update(req.SubscriptionID, params)

	if err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("sub.Update: %v", err)
		return
	}

	response.WriteJSON(w, struct {
		Subscription *stripe.Subscription `json:"subscription"`
	}{
		Subscription: updatedSubscription,
	}, nil)
}
