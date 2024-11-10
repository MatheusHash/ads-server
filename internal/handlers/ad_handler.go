package handlers

import (
	"ad-server/internal/database"
	"database/sql"
	"encoding/json"

	"log"
	"net/http"
)

// Estrutura para receber os dados do anúncio
type CreateAdRequest struct {
	Category string `json:"category"`
	Content  string `json:"content"`
}

// Estrutura para responder com os dados do anúncio criado
type AdResponse struct {
	ID          int    `json:"id"`
	Category    string `json:"category"`
	Content     string `json:"content"`
	Impressions int    `json:"impressions"`
}

// Estrutura para exibir um anúncio
type Ad struct {
	ID          int    `json:"id"`
	Category    string `json:"category"`
	Content     string `json:"content"`
	Impressions int    `json:"impressions"`
}

// Endpoint para obter um anúncio aleatório e incrementar suas impressões
func GetAd(w http.ResponseWriter, r *http.Request) {
	var ad Ad

	// Buscar um anúncio aleatório (ou baseado em filtros, caso aplicável)
	row := database.DB.QueryRow("SELECT id, category, content, impressions FROM ads ORDER BY RANDOM() LIMIT 1")
	if err := row.Scan(&ad.ID, &ad.Category, &ad.Content, &ad.Impressions); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Nenhum anúncio encontrado", http.StatusNotFound)
			return
		}
		log.Printf("Erro ao buscar anúncio: %v", err)
		http.Error(w, "Erro ao buscar anúncio", http.StatusInternalServerError)
		return
	}

	// Incrementar as impressões do anúncio
	_, err := database.DB.Exec("UPDATE ads SET impressions = impressions + 1 WHERE id = ?", ad.ID)
	if err != nil {
		log.Printf("Erro ao incrementar impressões: %v", err)
		http.Error(w, "Erro ao atualizar impressões", http.StatusInternalServerError)
		return
	}

	// Incrementar localmente para refletir o valor atualizado
	ad.Impressions++

	// Configurar o cabeçalho e enviar o anúncio em JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ad)
}

// Endpoint para criar um novo anúncio
func CreateAd(w http.ResponseWriter, r *http.Request) {
	var adRequest CreateAdRequest
	if err := json.NewDecoder(r.Body).Decode(&adRequest); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Inserir o anúncio no banco de dados
	result, err := database.DB.Exec("INSERT INTO ads (category, content, impressions) VALUES (?, ?, ?)", adRequest.Category, adRequest.Content, 0)
	if err != nil {
		log.Printf("Erro ao inserir anúncio: %v", err)
		http.Error(w, "Erro ao salvar anúncio", http.StatusInternalServerError)
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Erro ao obter ID do anúncio: %v", err)
		http.Error(w, "Erro ao salvar anúncio", http.StatusInternalServerError)
		return
	}

	// Construir a resposta com os dados do anúncio criado
	adResponse := AdResponse{
		ID:          int(lastID),
		Category:    adRequest.Category,
		Content:     adRequest.Content,
		Impressions: 0,
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adResponse)

}
