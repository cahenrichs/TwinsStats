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

func (r *Repository) GetPlayerByMLBID(mlbid int) (*models.Player, error) {
	var player models.Player
	err := r.db.Where("mlb.id = ?", mlbid).First(&player).Error
	if err != nil {
		return nil, fmt.Errorf("player milbid not found", err)
	}
	return &player, nil
}

func (r *Repository) SavePlayer(player *models.Player) error {
	return r.db.Save(player).Error
}

func (r *Repository) GetPlayerStats(playerID int, season int, statType string) (interface{}, error) {
	if statType == "hitting" {
		var stats models.HittingStats
		err := r.db.Where("player_id = ? AND season = ?", playerID, season).First(&stats).Error
		return stats, err
	}
	var stats models.PitchingStats
	err := r.db.Where("player_id = ? AND season = ?", playerID, season).First(&stats).Error
	return stats, err
}

func (r *Repository) SaveHittingStats(stats *models.HittingStats) error {
	return r.db.Save(stats).Error
}

func (r *Repository) SavePitchingStats(stats *models.PitchingStats) error {
	return r.db.Save(stats).Error
}

// Team
func (r *Repository) FindTeamByName(name string) (*models.Team, error) {
	var team models.Team
	name = strings.ToLower(name)
	err := r.db.Where("LOWER(name LIKE ? OR LOWER(nickname) LIKE ? OR LOWER(abbr) LIKE ?", "%"+name+"%", "%"+name+"%", "%"+name+"%").First(&team).Error
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}
	return &team, nil
}

func (r *Repository) GetTeamByMLBID(mlbid int) (*models.Team, error) {
	var team models.Team
	err := r.db.Where("milb_id = ?", mlbid).First(&team).Error
	if err != nil {
		return nil, fmt.Errorf("team not found by mlbid: %w", err)
	}
	return &team, nil
}

func (r *Repository) SaveTeam(team *models.Team) error {
	return r.db.Save(team).Error
}

func (r *Repository) SaveRoster(teamID uint, players []models.Player) error {
	return r.db.Where("current_team_id = ?", teamID).Delete(&models.Player{}).Error
}
