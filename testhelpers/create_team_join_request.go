package testhelpers

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	"gorm.io/gorm"
)

func CreateTeamJoinRequest(db *gorm.DB, candidateID, teamID uint, role string) (uint, error) {
	request := models.TeamJoinRequest{
		CandidateID: candidateID,
		TeamID:      teamID,
	}

	if err := db.Raw("SELECT id FROM roles WHERE name = ?", role).Scan(&request.RoleID).Error; err != nil {
		return 0, err
	}

	if result := db.Create(&request); result.Error != nil {
		return 0, result.Error
	}

	return request.ID, nil
}
