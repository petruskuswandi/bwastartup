package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/petruskuswandi/bwastartup.git/request"
	"github.com/petruskuswandi/bwastartup.git/service"
)

type sessionHandler struct {
	userService service.ServiceUser
}

func NewSessionHandler(userService service.ServiceUser) *sessionHandler {
	return &sessionHandler{userService}
}

func (h *sessionHandler) New(c *gin.Context) {
	c.HTML(http.StatusOK, "session_new.html", nil)
}

func (h *sessionHandler) Create(c *gin.Context) {
	var input request.LoginInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	user, err := h.userService.Login(input)
	if err != nil || user.Role != "admin" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	session := sessions.Default(c)
	session.Set("userID", user.ID)
	session.Set("userName", user.Name)
	session.Save()

	c.Redirect(http.StatusFound, "/users")
}

func (h *sessionHandler) Destroy(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, "/login")
}