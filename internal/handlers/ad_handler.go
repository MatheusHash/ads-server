package handlers

import (
	"ad-server/internal/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"log"
	"net/http"

	"github.com/gorilla/mux"
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

type RequestBody struct {
	AdId      int `json:"adId"`
	AccountId int `json:"accountId"`
}

func GetAllAds(w http.ResponseWriter, r *http.Request) {

	var ads []Ad

	rows, err := database.DB.Query("SELECT * FROM ads")

	if err != nil {
		log.Printf("Erro ao buscar anúncios: %v", err)
		http.Error(w, "Erro ao buscar anúncios", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // Garante que os recursos associados ao resultado da query sejam liberados

	for rows.Next() {
		var ad Ad
		if err := rows.Scan(&ad.ID, &ad.Category, &ad.Content, &ad.Impressions); err != nil {
			log.Printf("Erro ao escanear anúncio: %v", err)
			http.Error(w, "Erro ao processar anúncios", http.StatusInternalServerError)
			return
		}
		ads = append(ads, ad)
	}

	// log.Println(ads)

	// Verifica se houve algum erro durante a iteração
	if err := rows.Err(); err != nil {
		log.Printf("Erro durante a iteração: %v", err)
		http.Error(w, "Erro ao processar resultados", http.StatusInternalServerError)
		return
	}

	// Configurar o cabeçalho e enviar o anúncio em JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ads)
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

// Função para buscar um anúncio pelo ID (slug)
func GetAdById(w http.ResponseWriter, r *http.Request) {
	// Pegar o 'slug' (ID) da URL
	vars := mux.Vars(r)
	slug := vars["slug"]

	// Converter o 'slug' para um número inteiro (ID)
	adID, err := strconv.Atoi(slug)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var ad Ad

	// Buscar o anúncio pelo ID no banco de dados
	row := database.DB.QueryRow("SELECT id, category, content, impressions FROM ads WHERE id = ?", adID)
	if err := row.Scan(&ad.ID, &ad.Category, &ad.Content, &ad.Impressions); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Anúncio não encontrado", http.StatusNotFound)
			return
		}
		log.Printf("Erro ao buscar anúncio: %v", err)
		http.Error(w, "Erro ao buscar anúncio", http.StatusInternalServerError)
		return
	}

	// Retornar o anúncio em formato JSON
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

// Endpoint para registrar um click em um anúncio
func RegisterClick(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	log.Println("body", body)
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}
	adId := body.AdId
	accountId := body.AccountId
	// Incrementar as impressões do anúncio
	_, err := database.DB.Exec("UPDATE ads SET impressions = impressions + 1 WHERE id = ?", adId)
	if err != nil {
		http.Error(w, "Erro ao registrar clique", http.StatusInternalServerError)
		return
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05") // Formato padrão ISO 8601

	result, err := database.DB.Exec("INSERT INTO account_ads (account_id, ad_id, interaction_type, interaction_date) VALUES (?, ?, ?, ?)", accountId, adId, "click", currentTime)

	if err != nil {
		log.Printf("Erro ao relacionar click: %v", err)
		http.Error(w, "Erro!", http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	// Configurar o cabeçalho e enviar o anúncio em JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		"Clique registrado com sucesso",
	)

}
