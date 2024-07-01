package subscription

import (
	"backendgo/internal/response"
	"encoding/json"
	"fmt"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"log"
	"net/http"
)

func HandleCreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SuccessURL string `json:"successUrl"`
		CancelURL  string `json:"cancelUrl"`
		PriceID    string `json:"priceId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	// Create new Checkout Session for the order
	// Other optional params include:
	// [billing_address_collection] - to display billing address details on the page
	// [customer] - if you have an existing Stripe Customer ID
	// [payment_intent_data] - lets capture the payment later
	// [customer_email] - lets you prefill the email input in the form
	// [automatic_tax] - to automatically calculate sales tax, VAT and GST in the checkout page
	// For full details see https://stripe.com/docs/api/checkout/sessions/create

	// ?session_id={CHECKOUT_SESSION_ID} means the redirect will have the session ID
	// set as a query param
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(req.SuccessURL),
		CancelURL:  stripe.String(req.CancelURL),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price: stripe.String(req.PriceID),
				// For metered billing, do not pass quantity
				Quantity: stripe.Int64(1),
			},
		},
	}
	s, err := session.New(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("error while creating session %v", err.Error()), http.StatusInternalServerError)
		return
	}

	response.WriteJSON(w, s, nil)
}
