package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/petruskuswandi/bwastartup.git/request"
	"github.com/petruskuswandi/bwastartup.git/service"
)

type campaignHandler struct {
	campaignService service.ServiceCampaign
	userService     service.ServiceUser
}

func NewCampaignHandler(campaignService service.ServiceCampaign, userService service.ServiceUser) *campaignHandler {
	return &campaignHandler{campaignService, userService}
}

func (h *campaignHandler) Index(c *gin.Context) {
	campaigns, err := h.campaignService.GetCampaigns("")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campaigns})
}

func (h *campaignHandler) New(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := request.FormCreateCampaignInput{}
	input.Users = users

	c.HTML(http.StatusOK, "campaign_new.html", input)
}
