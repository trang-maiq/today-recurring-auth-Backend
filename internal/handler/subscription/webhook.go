package subscription

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"backendgo/internal/response"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/stripe/stripe-go/v72/webhook"
)

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("ioutil.ReadAll: %v", err)
		return
	}

	event, err := webhook.ConstructEvent(b, r.Header.Get("Stripe-Signature"), os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("webhook.ConstructEvent: %v", err)
		return
	}

	if event.Type == "invoice.payment_succeeded" {
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pi, _ := paymentintent.Get(
			invoice.PaymentIntent.ID,
			nil,
		)

		params := &stripe.SubscriptionParams{
			DefaultPaymentMethod: stripe.String(pi.PaymentMethod.ID),
		}
		sub.Update(invoice.Subscription.ID, params)
		fmt.Println("Default payment method set for subscription: ", pi.PaymentMethod)
	}
	fmt.Println("Payment succeeded for invoice: ", event.ID)
}
