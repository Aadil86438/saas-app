package main

import (
	"log"
	"net/http"
	"os"
)

// enableCORS is a security gate that allows your Vercel frontend to talk to Railway.
// It doesn't change your logic, just adds permission headers.
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 1. Initialize the Database
	InitDB()

	// 2. Setup your Routes (Cursor logic)
	http.HandleFunc("/api/auth/signup", SignupHandler)
	http.HandleFunc("/api/auth/login", LoginHandler)
	http.HandleFunc("/api/auth/logout", LogoutHandler)
	http.HandleFunc("/api/auth/verify", VerifyHandler)
	http.HandleFunc("/api/todos", TodoHandler)
    // Add a health check so Railway knows the app is alive
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Backend is running!"))
	})

	// 3. Get the Port from Railway
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback for local testing
	}

	log.Printf("Server starting on :%s", port)

	// 4. Start Server with CORS enabled
	err := http.ListenAndServe(":"+port, enableCORS(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}
