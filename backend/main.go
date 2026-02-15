package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	InitDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Backend is running! 🚀"))
	})

	http.HandleFunc("/api/auth/signup", CORSHandler(SignupAPI))
	http.HandleFunc("/api/auth/login", CORSHandler(LoginAPI))
	http.HandleFunc("/api/auth/logout", CORSHandler(LogoutAPI))
	http.HandleFunc("/api/auth/verify", CORSHandler(VerifyTokenAPI))

	http.HandleFunc("/api/todos", CORSHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			ListTodosAPI(w, r)
		} else if r.Method == http.MethodPost {
			CreateTodoAPI(w, r)
		}
	}))

	http.HandleFunc("/api/todos/", CORSHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			UpdateTodoAPI(w, r)
		} else if r.Method == http.MethodDelete {
			DeleteTodoAPI(w, r)
		}
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func CORSHandler(pHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		pHandler(w, r)
	}
}

func ReadBody(r *http.Request) string {
	lBody, lErr := ioutil.ReadAll(r.Body)
	if lErr != nil {
		return ""
	}
	return string(lBody)
}

func SendJSONResponse(w http.ResponseWriter, pResponse APIResponse, pStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(pStatus)

	lJSON, lErr := json.Marshal(pResponse)
	if lErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(lJSON)
}

func SendErrorResponse(w http.ResponseWriter, pMessage string, pStatus int) {
	lResponse := APIResponse{
		Status:  "e",
		Message: pMessage,
		Data:    nil,
	}
	SendJSONResponse(w, lResponse, pStatus)
}

