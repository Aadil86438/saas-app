package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func CreateTodoAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateTodoAPI(+)")
	
	if r.Method != http.MethodPost {
		SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("CreateTodoAPI(-)")
		return
	}
	
	lToken := r.Header.Get("Authorization")
	if lToken == "" {
		SendErrorResponse(w, "Missing authorization token", http.StatusUnauthorized)
		log.Println("CreateTodoAPI(-)")
		return
	}
	
	lUser, lErr := GetUserFromToken(lToken)
	if lErr != nil {
		SendErrorResponse(w, "Invalid token", http.StatusUnauthorized)
		log.Println("CreateTodoAPI(-) error:", lErr)
		return
	}
	
	var lReq CreateTodoRequest
	lErr = json.Unmarshal([]byte(ReadBody(r)), &lReq)
	if lErr != nil {
		SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		log.Println("CreateTodoAPI(-) error:", lErr)
		return
	}
	
	lTodo, lErr := CreateTodo(lUser.ID, lReq.Title, lReq.Content)
	if lErr != nil {
		SendErrorResponse(w, lErr.Error(), http.StatusBadRequest)
		log.Println("CreateTodoAPI(-) error:", lErr)
		return
	}
	
	lResponse := APIResponse{
		Status:  "s",
		Message: "Todo created successfully",
		Data:    lTodo,
	}
	
	SendJSONResponse(w, lResponse, http.StatusOK)
	log.Println("CreateTodoAPI(-)")
}

func ListTodosAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("ListTodosAPI(+)")
	
	if r.Method != http.MethodGet {
		SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("ListTodosAPI(-)")
		return
	}
	
	lToken := r.Header.Get("Authorization")
	if lToken == "" {
		SendErrorResponse(w, "Missing authorization token", http.StatusUnauthorized)
		log.Println("ListTodosAPI(-)")
		return
	}
	
	lUser, lErr := GetUserFromToken(lToken)
	if lErr != nil {
		SendErrorResponse(w, "Invalid token", http.StatusUnauthorized)
		log.Println("ListTodosAPI(-) error:", lErr)
		return
	}
	
	lTodosArr, lErr := ListTodos(lUser.ID)
	if lErr != nil {
		SendErrorResponse(w, lErr.Error(), http.StatusInternalServerError)
		log.Println("ListTodosAPI(-) error:", lErr)
		return
	}
	
	lResponse := APIResponse{
		Status:  "s",
		Message: "Todos retrieved successfully",
		Data:    lTodosArr,
	}
	
	SendJSONResponse(w, lResponse, http.StatusOK)
	log.Println("ListTodosAPI(-)")
}

func UpdateTodoAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateTodoAPI(+)")
	
	if r.Method != http.MethodPut {
		SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("UpdateTodoAPI(-)")
		return
	}
	
	lToken := r.Header.Get("Authorization")
	if lToken == "" {
		SendErrorResponse(w, "Missing authorization token", http.StatusUnauthorized)
		log.Println("UpdateTodoAPI(-)")
		return
	}
	
	lUser, lErr := GetUserFromToken(lToken)
	if lErr != nil {
		SendErrorResponse(w, "Invalid token", http.StatusUnauthorized)
		log.Println("UpdateTodoAPI(-) error:", lErr)
		return
	}
	
	lPathParts := strings.Split(r.URL.Path, "/")
	if len(lPathParts) < 3 {
		SendErrorResponse(w, "Invalid todo ID", http.StatusBadRequest)
		log.Println("UpdateTodoAPI(-)")
		return
	}
	
	lTodoID, lErr := strconv.Atoi(lPathParts[len(lPathParts)-1])
	if lErr != nil {
		SendErrorResponse(w, "Invalid todo ID", http.StatusBadRequest)
		log.Println("UpdateTodoAPI(-) error:", lErr)
		return
	}
	
	var lReq UpdateTodoRequest
	lErr = json.Unmarshal([]byte(ReadBody(r)), &lReq)
	if lErr != nil {
		SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		log.Println("UpdateTodoAPI(-) error:", lErr)
		return
	}
	
	lTodo, lErr := UpdateTodo(lUser.ID, lTodoID, lReq.Title, lReq.Content, lReq.Completed)
	if lErr != nil {
		SendErrorResponse(w, lErr.Error(), http.StatusBadRequest)
		log.Println("UpdateTodoAPI(-) error:", lErr)
		return
	}
	
	lResponse := APIResponse{
		Status:  "s",
		Message: "Todo updated successfully",
		Data:    lTodo,
	}
	
	SendJSONResponse(w, lResponse, http.StatusOK)
	log.Println("UpdateTodoAPI(-)")
}

func DeleteTodoAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteTodoAPI(+)")
	
	if r.Method != http.MethodDelete {
		SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("DeleteTodoAPI(-)")
		return
	}
	
	lToken := r.Header.Get("Authorization")
	if lToken == "" {
		SendErrorResponse(w, "Missing authorization token", http.StatusUnauthorized)
		log.Println("DeleteTodoAPI(-)")
		return
	}
	
	lUser, lErr := GetUserFromToken(lToken)
	if lErr != nil {
		SendErrorResponse(w, "Invalid token", http.StatusUnauthorized)
		log.Println("DeleteTodoAPI(-) error:", lErr)
		return
	}
	
	lPathParts := strings.Split(r.URL.Path, "/")
	if len(lPathParts) < 3 {
		SendErrorResponse(w, "Invalid todo ID", http.StatusBadRequest)
		log.Println("DeleteTodoAPI(-)")
		return
	}
	
	lTodoID, lErr := strconv.Atoi(lPathParts[len(lPathParts)-1])
	if lErr != nil {
		SendErrorResponse(w, "Invalid todo ID", http.StatusBadRequest)
		log.Println("DeleteTodoAPI(-) error:", lErr)
		return
	}
	
	lErr = DeleteTodo(lUser.ID, lTodoID)
	if lErr != nil {
		SendErrorResponse(w, lErr.Error(), http.StatusBadRequest)
		log.Println("DeleteTodoAPI(-) error:", lErr)
		return
	}
	
	lResponse := APIResponse{
		Status:  "s",
		Message: "Todo deleted successfully",
		Data:    nil,
	}
	
	SendJSONResponse(w, lResponse, http.StatusOK)
	log.Println("DeleteTodoAPI(-)")
}

func CreateTodo(pUserID int, pTitle string, pContent string) (*Todo, error) {
	log.Println("CreateTodo(+)")
	
	lQuery := "INSERT INTO todos (user_id, title, content) VALUES ($1, $2, $3) RETURNING id, user_id, title, content, completed, created_at"
	lDB := GetDB()
	
	var lTodo Todo
	lErr := lDB.QueryRow(lQuery, pUserID, pTitle, pContent).Scan(&lTodo.ID, &lTodo.UserID, &lTodo.Title, &lTodo.Content, &lTodo.Completed, &lTodo.CreatedAt)
	if lErr != nil {
		log.Println("CreateTodo(-) error:", lErr)
		return nil, lErr
	}
	
	log.Println("CreateTodo(-)")
	return &lTodo, nil
}

func ListTodos(pUserID int) ([]Todo, error) {
	log.Println("ListTodos(+)")
	
	lQuery := "SELECT id, user_id, title, content, completed, created_at FROM todos WHERE user_id = $1 ORDER BY created_at DESC"
	lDB := GetDB()
	
	lRows, lErr := lDB.Query(lQuery, pUserID)
	if lErr != nil {
		log.Println("ListTodos(-) error:", lErr)
		return nil, lErr
	}
	defer lRows.Close()
	
	var lTodosArr []Todo
	for lRows.Next() {
		var lTodo Todo
		lErr := lRows.Scan(&lTodo.ID, &lTodo.UserID, &lTodo.Title, &lTodo.Content, &lTodo.Completed, &lTodo.CreatedAt)
		if lErr != nil {
			log.Println("ListTodos(-) error:", lErr)
			continue
		}
		lTodosArr = append(lTodosArr, lTodo)
	}
	
	log.Println("ListTodos(-)")
	return lTodosArr, nil
}

func UpdateTodo(pUserID int, pTodoID int, pTitle string, pContent string, pCompleted bool) (*Todo, error) {
	log.Println("UpdateTodo(+)")
	
	lQuery := "UPDATE todos SET title = $1, content = $2, completed = $3 WHERE id = $4 AND user_id = $5 RETURNING id, user_id, title, content, completed, created_at"
	lDB := GetDB()
	
	var lTodo Todo
	lErr := lDB.QueryRow(lQuery, pTitle, pContent, pCompleted, pTodoID, pUserID).Scan(&lTodo.ID, &lTodo.UserID, &lTodo.Title, &lTodo.Content, &lTodo.Completed, &lTodo.CreatedAt)
	if lErr != nil {
		log.Println("UpdateTodo(-) error:", lErr)
		return nil, lErr
	}
	
	log.Println("UpdateTodo(-)")
	return &lTodo, nil
}

func DeleteTodo(pUserID int, pTodoID int) error {
	log.Println("DeleteTodo(+)")
	
	lQuery := "DELETE FROM todos WHERE id = $1 AND user_id = $2"
	lDB := GetDB()
	
	lResult, lErr := lDB.Exec(lQuery, pTodoID, pUserID)
	if lErr != nil {
		log.Println("DeleteTodo(-) error:", lErr)
		return lErr
	}
	
	lRowsAffected, lErr := lResult.RowsAffected()
	if lErr != nil {
		log.Println("DeleteTodo(-) error:", lErr)
		return lErr
	}
	
	if lRowsAffected == 0 {
		log.Println("DeleteTodo(-) error: todo not found")
		return errors.New("todo not found")
	}
	
	log.Println("DeleteTodo(-)")
	return nil
}

