package subscription

import (
	"net/http"

	"backendgo/internal/response"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
)

func HandleListSubscriptions(w http.ResponseWriter, r *http.Request) {
	// Read customer from cookie to simulate auth
	cookie, _ := r.Cookie("customer")
	customerID := cookie.Value

	params := &stripe.SubscriptionListParams{
		Customer: customerID,
		Status:   "all",
	}
	params.AddExpand("data.default_payment_method")
	i := sub.List(params)

	response.WriteJSON(w, struct {
		Subscriptions *stripe.SubscriptionList `json:"subscriptions"`
	}{
		Subscriptions: i.SubscriptionList(),
	}, nil)
}
