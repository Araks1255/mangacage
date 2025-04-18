package testhelpers

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	"gorm.io/gorm"
)

func CreateTeamJoinRequest(db *gorm.DB, candidateID, teamID uint) (uint, error) {
	request := models.TeamJoinRequest{
		CandidateID: candidateID,
		TeamID:      teamID,
	}

	if result := db.Create(&request); result.Error != nil {
		return 0, result.Error
	}

	return request.ID, nil
}
