package joinrequests

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"
)

func (h handler) DeclineTeamJoinRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	joinRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id заявки"})
		return
	}

	var joinRequest models.TeamJoinRequest

	err = h.DB.Raw(
		"DELETE FROM team_join_requests WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?) RETURNING *",
		joinRequestID, claims.ID,
	).Scan(&joinRequest).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if joinRequest.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "заявка на вступление в вашу команду не найдена"})
		return
	}

	c.JSON(200, gin.H{"success": "заявка на вступление в вашу команду успешно отменена"})

	if _, err := h.NotificationsCLient.NotifyAboutTeamJoinRequestResponse(
		c.Request.Context(),
		&pb.TeamJoinRequestResponse{
			Result: enums.ResultOfTeamJoinRequest_RESULT_OF_TEAM_JOIN_REQUEST_DECLINED,
			TeamID: uint64(joinRequest.TeamID),
			UserID: uint64(joinRequest.CandidateID),
		},
	); err != nil {
		log.Println(err)
	}
}
