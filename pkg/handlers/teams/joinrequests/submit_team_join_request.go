package joinrequests

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) SubmitTeamJoinRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	joinRequest, err := mapSubmitTeamJoinRequestBodyIntoTeamJoinRequest(c.ShouldBindJSON, c.Param, c.Request.ContentLength, claims.ID)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	code, err := checkSubmitTeamJoinRequestConflicts(h.DB, *joinRequest, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	err = h.DB.Create(&joinRequest).Error

	if err != nil {
		code, err := parseSubmitTeamJoinRequestError(err)
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"success": "заявка на вступление в команду успешно отправлена"})

	var message string
	if joinRequest.IntroductoryMessage != nil {
		message = *joinRequest.IntroductoryMessage
	}

	if _, err := h.NotificationsCLient.NotifyAboutSubmittedTeamJoinRequest(
		c.Request.Context(), &pb.TeamJoinRequest{
			TeamID:      uint64(joinRequest.TeamID),
			CandidateID: uint64(joinRequest.CandidateID),
			Message:     message,
		},
	); err != nil {
		log.Println(err)
	}
}

func mapSubmitTeamJoinRequestBodyIntoTeamJoinRequest(
	bindFn func(any) error,
	paramFn func(string) string,
	contentLength int64,
	userID uint,
) (
	*models.TeamJoinRequest,
	error,
) {
	teamID, err := strconv.ParseUint(paramFn("id"), 10, 64)
	if err != nil {
		return nil, err
	}

	var requestBody dto.CreateTeamJoinRequestDTO

	if err := bindFn(&requestBody); err != nil && contentLength != 0 {
		return nil, err
	}

	res := requestBody.ToTeamJoinRequest(userID, uint(teamID))

	return &res, nil
}

func checkSubmitTeamJoinRequestConflicts(db *gorm.DB, joinRequest models.TeamJoinRequest, userID uint) (code int, err error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE id = ? AND team_id IS NOT NULL) AS user_have_team"

	if joinRequest.RoleID != nil {
		query += ",EXISTS(SELECT 1 FROM roles WHERE id = ? AND type = 'team') AS role_exists"
	}

	var check struct {
		UserHaveTeam bool
		RoleExists   bool
	}

	if joinRequest.RoleID != nil {
		err = db.Raw(query, userID, joinRequest.RoleID).Scan(&check).Error
	} else {
		err = db.Raw(query, userID).Scan(&check).Error
	}

	if err != nil {
		return 500, err
	}

	if check.UserHaveTeam {
		return 409, errors.New("вы уже состоите в команде перевода")
	}

	if joinRequest.RoleID != nil && !check.RoleExists {
		return 404, errors.New("роль не найдена")
	}

	return 0, nil
}

func parseSubmitTeamJoinRequestError(err error) (code int, parsedErr error) {
	if dbErrors.IsUniqueViolation(err, constraints.UniqTeamJoinRequest) {
		return 409, errors.New("вы уже оставили заявку на вступление в эту команду")
	}

	if dbErrors.IsForeignKeyViolation(err, constraints.FkTeamJoinRequestsTeam) {
		return 404, errors.New("команда не найдена")
	}

	return 500, err
}
