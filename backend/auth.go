package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SignupAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("SignupAPI(+)")
	
	if r.Method != http.MethodPost {
		SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("SignupAPI(-)")
		return
	}
	
	var lReq SignupRequest
	lErr := json.Unmarshal([]byte(ReadBody(r)), &lReq)
	if lErr != nil {
		SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		log.Println("SignupAPI(-) error:", lErr)
		return
	}
	
	lUser, lErr := Signup(lReq.Username, lReq.Email, lReq.Password)
	if lErr != nil {
		SendErrorResponse(w, lErr.Error(), http.StatusBadRequest)
		log.Println("SignupAPI(-) error:", lErr)
		return
	}
	
	lToken, lErr := CreateSession(lUser.ID)
	if lErr != nil {
		SendErrorResponse(w, "Failed to create session", http.StatusInternalServerError)
		log.Println("SignupAPI(-) error:", lErr)
		return
	}
	
	lResponse := APIResponse{
		Status:  "s",
		Message: "Signup successful",
		Data: map[string]interface{}{
			"user":  lUser,
			"token": lToken,
		},
	}
	
	SendJSONResponse(w, lResponse, http.StatusOK)
	log.Println("SignupAPI(-)")
}

func LoginAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("LoginAPI(+)")
	
	if r.Method != http.MethodPost {
		SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("LoginAPI(-)")
		return
	}
	
	var lReq LoginRequest
	lErr := json.Unmarshal([]byte(ReadBody(r)), &lReq)
	if lErr != nil {
		SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		log.Println("LoginAPI(-) error:", lErr)
		return
	}
	
	lUser, lErr := Login(lReq.Username, lReq.Password)
	if lErr != nil {
		SendErrorResponse(w, lErr.Error(), http.StatusUnauthorized)
		log.Println("LoginAPI(-) error:", lErr)
		return
	}
	
	lToken, lErr := CreateSession(lUser.ID)
	if lErr != nil {
		SendErrorResponse(w, "Failed to create session", http.StatusInternalServerError)
		log.Println("LoginAPI(-) error:", lErr)
		return
	}
	
	lResponse := APIResponse{
		Status:  "s",
		Message: "Login successful",
		Data: map[string]interface{}{
			"user":  lUser,
			"token": lToken,
		},
	}
	
	SendJSONResponse(w, lResponse, http.StatusOK)
	log.Println("LoginAPI(-)")
}

func LogoutAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("LogoutAPI(+)")
	
	if r.Method != http.MethodPost {
		SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("LogoutAPI(-)")
		return
	}
	
	lToken := r.Header.Get("Authorization")
	if lToken == "" {
		SendErrorResponse(w, "Missing authorization token", http.StatusUnauthorized)
		log.Println("LogoutAPI(-)")
		return
	}
	
	lErr := Logout(lToken)
	if lErr != nil {
		SendErrorResponse(w, lErr.Error(), http.StatusBadRequest)
		log.Println("LogoutAPI(-) error:", lErr)
		return
	}
	
	lResponse := APIResponse{
		Status:  "s",
		Message: "Logout successful",
		Data:    nil,
	}
	
	SendJSONResponse(w, lResponse, http.StatusOK)
	log.Println("LogoutAPI(-)")
}

func VerifyTokenAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("VerifyTokenAPI(+)")
	
	if r.Method != http.MethodGet {
		SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("VerifyTokenAPI(-)")
		return
	}
	
	lToken := r.Header.Get("Authorization")
	if lToken == "" {
		SendErrorResponse(w, "Missing authorization token", http.StatusUnauthorized)
		log.Println("VerifyTokenAPI(-)")
		return
	}
	
	lUser, lErr := VerifyToken(lToken)
	if lErr != nil {
		SendErrorResponse(w, lErr.Error(), http.StatusUnauthorized)
		log.Println("VerifyTokenAPI(-) error:", lErr)
		return
	}
	
	lResponse := APIResponse{
		Status:  "s",
		Message: "Token valid",
		Data:    lUser,
	}
	
	SendJSONResponse(w, lResponse, http.StatusOK)
	log.Println("VerifyTokenAPI(-)")
}

func Signup(pUsername string, pEmail string, pPassword string) (*User, error) {
	log.Println("Signup(+)")
	
	lHashedPassword, lErr := bcrypt.GenerateFromPassword([]byte(pPassword), bcrypt.DefaultCost)
	if lErr != nil {
		log.Println("Signup(-) error:", lErr)
		return nil, lErr
	}
	
	lQuery := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, username, email"
	lDB := GetDB()
	
	var lUser User
	lErr = lDB.QueryRow(lQuery, pUsername, pEmail, string(lHashedPassword)).Scan(&lUser.ID, &lUser.Username, &lUser.Email)
	if lErr != nil {
		log.Println("Signup(-) error:", lErr)
		return nil, lErr
	}
	
	log.Println("Signup(-)")
	return &lUser, nil
}

func Login(pUsername string, pPassword string) (*User, error) {
	log.Println("Login(+)")
	
	lQuery := "SELECT id, username, email, password FROM users WHERE username = $1"
	lDB := GetDB()
	
	var lUser User
	lErr := lDB.QueryRow(lQuery, pUsername).Scan(&lUser.ID, &lUser.Username, &lUser.Email, &lUser.Password)
	if lErr != nil {
		log.Println("Login(-) error:", lErr)
		return nil, lErr
	}
	
	lErr = bcrypt.CompareHashAndPassword([]byte(lUser.Password), []byte(pPassword))
	if lErr != nil {
		log.Println("Login(-) error:", lErr)
		return nil, lErr
	}
	
	lUser.Password = ""
	log.Println("Login(-)")
	return &lUser, nil
}

func CreateSession(pUserID int) (string, error) {
	log.Println("CreateSession(+)")
	
	lTokenBytes := make([]byte, 32)
	_, lErr := rand.Read(lTokenBytes)
	if lErr != nil {
		log.Println("CreateSession(-) error:", lErr)
		return "", lErr
	}
	
	lToken := hex.EncodeToString(lTokenBytes)
	lExpiresAt := time.Now().Add(24 * time.Hour)
	
	lQuery := "INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)"
	lDB := GetDB()
	
	_, lErr = lDB.Exec(lQuery, pUserID, lToken, lExpiresAt)
	if lErr != nil {
		log.Println("CreateSession(-) error:", lErr)
		return "", lErr
	}
	
	log.Println("CreateSession(-)")
	return lToken, nil
}

func VerifyToken(pToken string) (*User, error) {
	log.Println("VerifyToken(+)")
	
	lQuery := "SELECT s.user_id, u.id, u.username, u.email FROM sessions s JOIN users u ON s.user_id = u.id WHERE s.token = $1 AND s.expires_at > NOW()"
	lDB := GetDB()
	
	var lUser User
	lErr := lDB.QueryRow(lQuery, pToken).Scan(&lUser.ID, &lUser.ID, &lUser.Username, &lUser.Email)
	if lErr != nil {
		log.Println("VerifyToken(-) error:", lErr)
		return nil, lErr
	}
	
	log.Println("VerifyToken(-)")
	return &lUser, nil
}

func Logout(pToken string) error {
	log.Println("Logout(+)")
	
	lQuery := "DELETE FROM sessions WHERE token = $1"
	lDB := GetDB()
	
	_, lErr := lDB.Exec(lQuery, pToken)
	if lErr != nil {
		log.Println("Logout(-) error:", lErr)
		return lErr
	}
	
	log.Println("Logout(-)")
	return nil
}

func GetUserFromToken(pToken string) (*User, error) {
	log.Println("GetUserFromToken(+)")
	
	lUser, lErr := VerifyToken(pToken)
	if lErr != nil {
		log.Println("GetUserFromToken(-) error:", lErr)
		return nil, lErr
	}
	
	log.Println("GetUserFromToken(-)")
	return lUser, nil
}

