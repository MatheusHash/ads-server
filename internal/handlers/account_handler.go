package handlers

import (
	"ad-server/internal/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type AccountCreate struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type AccountResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {

	var account AccountCreate

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "invalid data!", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if account.Type != "company" && account.Type != "personal" {
		account.Type = "personal"
	}

	result, err := database.DB.Exec("INSERT INTO accounts (name, type) VALUES (?, ?)", account.Name, account.Type)

	if err != nil {
		log.Printf("Erro ao criar a conta: %v", err)
		http.Error(w, "Erro ao criar conta", http.StatusInternalServerError)
		return
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		log.Printf("Erro ao buscar o item inserido: %v", err)
		http.Error(w, " Erro ao buscar/salvar conta", http.StatusInternalServerError)
		return
	}

	accountResponse := AccountResponse{
		ID:   int(lastId),
		Name: account.Name,
		Type: account.Type,
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accountResponse)
}
func GetAccountById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	accountId, err := strconv.Atoi(slug)

	if err != nil {
		http.Error(w, "Erro ao conveter slug.", http.StatusInternalServerError)
		return
	}

	var account Account

	row := database.DB.QueryRow("SELECT * FROM accounts WHERE id = ?", accountId)

	if err := row.Scan(&account.ID, &account.Name, &account.Type); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Conta nao encontrada.", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)

}
