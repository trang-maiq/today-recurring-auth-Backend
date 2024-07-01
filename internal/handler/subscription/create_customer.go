package subscription

import (
	"encoding/json"
	"log"
	"net/http"

	"backendgo/internal/response"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
)

func HandleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(req.Email),
	}

	c, err := customer.New(params)
	if err != nil {
		response.WriteJSON(w, nil, err)
		log.Printf("customer.New: %v", err)
		return
	}

	// You should store the ID of the customer in your database alongside your
	// users. This sample uses cookies to simulate auth.
	http.SetCookie(w, &http.Cookie{
		Name:  "customer",
		Value: c.ID,
	})

	response.WriteJSON(w, struct {
		Customer *stripe.Customer `json:"customer"`
	}{
		Customer: c,
	}, nil)
}
