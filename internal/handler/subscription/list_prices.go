package subscription

import (
	"net/http"
	"os"

	"backendgo/internal/response"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/price"
)

func HandleGetListPrices(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	params := &stripe.PriceListParams{
		LookupKeys: stripe.StringSlice([]string{"sample_basic", "sample_premium"}),
	}

	prices := make([]*stripe.Price, 0)

	i := price.List(params)
	for i.Next() {
		prices = append(prices, i.Price())
	}

	response.WriteJSON(w, struct {
		PublishableKey string          `json:"publishableKey"`
		Prices         []*stripe.Price `json:"prices"`
	}{
		PublishableKey: os.Getenv("STRIPE_PUBLISHABLE_KEY"),
		Prices:         prices,
	}, nil)
}
