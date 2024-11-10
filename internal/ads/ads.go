package ads

import "ad-server/internal/models"

// Função básica para selecionar anúncio com base na categoria
func SelectAdByCategory(category string, ads []models.Ad) *models.Ad {
	for _, ad := range ads {
		if ad.Category == category {
			return &ad
		}
	}
	return nil
}
