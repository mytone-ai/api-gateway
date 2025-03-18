package main

import (
	"log"
	"net/http"
	"time"

	_ "api-gateway/docs" // This will be auto-generated

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title API Gateway
// @version 1.0
// @description This is an API Gateway server.
// @host localhost:8000
// @BasePath /api/v1

type Gateway struct {
	locationServiceURL string
	// Add other service URLs as needed
}

// @Summary Health check endpoint
// @Description Get the health status of the API
// @Tags health
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /health [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// @Summary Proxy location service requests
// @Description Forward requests to the location service
// @Tags locations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} interface{}
// @Failure 401 {string} string "Unauthorized"
// @Failure 502 {string} string "Bad Gateway"
// @Router /api/v1/locations [get]
func (g *Gateway) proxyLocationService(w http.ResponseWriter, r *http.Request) {
	// Create a new request to the location service
	proxyReq, err := http.NewRequest(r.Method, g.locationServiceURL+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	// Copy headers
	for header, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for header, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy body
	if _, err := http.MaxBytesReader(w, resp.Body, 1048576).WriteTo(w); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func main() {
	gateway := &Gateway{
		locationServiceURL: "http://location-service:8080", // or environment variable
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300,
	}))

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))

	// Auth middleware
	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// Verify token here
			next.ServeHTTP(w, r)
		})
	}

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/health", healthCheck)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		// Location service routes
		r.Route("/api/v1/locations", func(r chi.Router) {
			r.Get("/", gateway.proxyLocationService)
			r.Post("/", gateway.proxyLocationService)
			// Add other routes as needed
		})
	})

	log.Fatal(http.ListenAndServe(":8000", r))
}
