package db

import (
	"fmt"
	"strings"

	"github.com/cahenrichs/TwinsStats/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Player
func (r *Repository) FindPlayerByName(name string) (*models.Player, error) {
	var player models.Player
	name = strings.ToLower(name)
	err := r.db.Where("LOWER(full_name) LIKE ?", "%"+name+"%").First(&player).Error
	if err != nil {
		return nil, fmt.Errorf("player not found %w", err)
	}
	return &player, nil
}
