package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v72"

    
    "backendgo/internal/db"
    "backendgo/internal/db/migration"
	
	
	"backendgo/internal/handler/auth"
    
	
	
	"backendgo/internal/handler/subscription"
    
)

func main() {
	slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	if err := godotenv.Load(); err != nil {
		slog.Warn("godotenv.Load failed", slog.Any("error", err))
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	// For sample support and debugging, not required for production:
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "backend",
		Version: "0.0.1",
	})

    
	migration.Migrate()
	db.Setup()
	

	mux := http.NewServeMux()

    
	mux.HandleFunc("/config", subscription.HandleGetListPrices)
	mux.HandleFunc("/webhook/stripe", subscription.HandleWebhook)
	mux.HandleFunc("/customers", subscription.HandleCreateCustomer)
    

    
	// recurring
	mux.HandleFunc("/subscriptions", subscription.HandleSubscriptions)
	mux.HandleFunc("/subscriptions/{id}", subscription.HandleSubscription)
	mux.HandleFunc("/subscriptions/invoice/preview", subscription.HandleInvoicePreview)
	mux.HandleFunc("/recurring/checkout/sessions", subscription.HandleCheckoutSessions)
    

    
    
	// auth
	mux.HandleFunc("/login", auth.HandleLoginWithEmailPassword)
	mux.HandleFunc("/register", auth.Register)
	

	addr := "0.0.0.0:4242"
	log.Printf("Listening on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, withCORS(mux)))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("received request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("remote_addr", r.RemoteAddr),
		)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
