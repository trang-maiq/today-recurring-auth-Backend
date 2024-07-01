package subscription

import (
	"log"
	"net/http"
	"os"
	"strings"

	"backendgo/internal/response"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/sub"
)

func HandleInvoicePreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	// Read customer from cookie to simulate auth
	cookie, _ := r.Cookie("customer")
	customerID := cookie.Value

	query := r.URL.Query()
	subscriptionID := query.Get("subscriptionId")
	newPriceLookupKey := strings.ToUpper(query.Get("newPriceLookupKey"))

	s, err := sub.Get(subscriptionID, nil)
	if err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("sub.Get: %v", err)
		return
	}
	params := &stripe.InvoiceParams{
		Customer:     stripe.String(customerID),
		Subscription: stripe.String(subscriptionID),
		SubscriptionItems: []*stripe.SubscriptionItemsParams{{
			ID:    stripe.String(s.Items.Data[0].ID),
			Price: stripe.String(os.Getenv(newPriceLookupKey)),
		}},
	}
	in, err := invoice.GetNext(params)

	if err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("invoice.GetNext: %v", err)
		return
	}

	response.WriteJSON(w, struct {
		Invoice *stripe.Invoice `json:"invoice"`
	}{
		Invoice: in,
	}, nil)
}
