package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/f4ke-n0name/avito/internal/domain/entities"
	"github.com/f4ke-n0name/avito/internal/domain/errors"
	"github.com/f4ke-n0name/avito/internal/domain/services/interfaces"
)

type Server struct {
	pr    interfaces.PRService
	users interfaces.UserService
	teams interfaces.TeamService
}

func NewServer(pr interfaces.PRService, users interfaces.UserService, teams interfaces.TeamService) *Server {
	return &Server{pr: pr, users: users, teams: teams}
}

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.POST("/team/add", s.addTeam)
	r.GET("/team/get", s.getTeam)

	r.POST("/users/setIsActive", s.setIsActive)
	r.GET("/users/getReview", s.getReviewList)

	r.POST("/pullRequest/create", s.createPR)
	r.POST("/pullRequest/merge", s.mergePR)
	r.POST("/pullRequest/reassign", s.reassign)
}

type TeamAddRequest struct {
	TeamName string `json:"team_name" binding:"required"`
	Members  []struct {
		UserID   string `json:"user_id" binding:"required"`
		Username string `json:"username" binding:"required"`
		IsActive bool   `json:"is_active"`
	} `json:"members" binding:"required"`
}

func (s *Server) addTeam(c *gin.Context) {
	var req TeamAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newBadRequest(err))
		return
	}

	team := &entities.Team{
		TeamName: req.TeamName,
	}
	for _, m := range req.Members {
		team.Members = append(team.Members, entities.User{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
			TeamName: req.TeamName,
		})
	}

	created, err := s.teams.CreateTeam(c, team)
	if err != nil {
		if err == errors.ErrTeamExists {
			c.JSON(http.StatusBadRequest, errorResponse("TEAM_EXISTS", "team_name already exists"))
			return
		}
		c.JSON(http.StatusInternalServerError, newInternal(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"team": created})
}

func (s *Server) getTeam(c *gin.Context) {
	name := c.Query("team_name")
	t, err := s.teams.GetTeam(c, name)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse("NOT_FOUND", "team not found"))
		return
	}
	c.JSON(http.StatusOK, t)
}

func (s *Server) setIsActive(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		IsActive bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newBadRequest(err))
		return
	}

	u, err := s.users.SetIsActive(c, req.UserID, req.IsActive)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse("NOT_FOUND", "user not found"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": u})
}

func (s *Server) getReviewList(c *gin.Context) {
	uid := c.Query("user_id")

	list, err := s.pr.ListByReviewer(c, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, newInternal(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       uid,
		"pull_requests": list,
	})
}

func (s *Server) createPR(c *gin.Context) {
	var req struct {
		PRID   string `json:"pull_request_id" binding:"required"`
		Name   string `json:"pull_request_name" binding:"required"`
		Author string `json:"author_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newBadRequest(err))
		return
	}

	pr, err := s.pr.CreatePR(c, req.PRID, req.Name, req.Author)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			c.JSON(http.StatusNotFound, errorResponse("NOT_FOUND", "author not found"))
		case errors.ErrPRExists:
			c.JSON(http.StatusConflict, errorResponse("PR_EXISTS", "PR id already exists"))
		default:
			c.JSON(http.StatusInternalServerError, newInternal(err))
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"pr": pr})
}

func (s *Server) mergePR(c *gin.Context) {
	var req struct {
		PRID string `json:"pull_request_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newBadRequest(err))
		return
	}

	pr, err := s.pr.Merge(c, req.PRID)
	if err != nil {
		if err == errors.ErrPRNotFound {
			c.JSON(http.StatusNotFound, errorResponse("NOT_FOUND", "PR not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, newInternal(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"pr": pr})
}

func (s *Server) reassign(c *gin.Context) {
	var req struct {
		PRID      string `json:"pull_request_id" binding:"required"`
		OldUserID string `json:"old_user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newBadRequest(err))
		return
	}

	pr, newID, err := s.pr.ReplaceReviewer(c, req.PRID, req.OldUserID)
	if err != nil {
		switch err {
		case errors.ErrPRNotFound, errors.ErrUserNotFound:
			c.JSON(http.StatusNotFound, errorResponse("NOT_FOUND", "PR or user not found"))
		case errors.ErrPRAlreadyMerged:
			c.JSON(http.StatusConflict, errorResponse("PR_MERGED", "cannot reassign on merged PR"))
		case errors.ErrNoSuchReviewer:
			c.JSON(http.StatusConflict, errorResponse("NOT_ASSIGNED", "reviewer is not assigned to this PR"))
		case errors.ErrNoCandidates:
			c.JSON(http.StatusConflict, errorResponse("NO_CANDIDATE", "no active replacement candidate in team"))
		default:
			c.JSON(http.StatusInternalServerError, newInternal(err))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"pr": pr, "replaced_by": newID})
}

func errorResponse(code, msg string) gin.H {
	return gin.H{
		"error": gin.H{
			"code":    code,
			"message": msg,
		},
	}
}

func newBadRequest(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func newInternal(err error) gin.H {
	return gin.H{"error": err.Error()}
}
