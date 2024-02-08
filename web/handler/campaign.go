package handler

import (
	"fmt"
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

func (h *campaignHandler) Create(c *gin.Context) {
	var input request.FormCreateCampaignInput

	err := c.ShouldBind(&input)
	if err != nil {
		users, e := h.userService.GetAllUsers()
		if e != nil {
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}

		input.Users = users
		input.Error = err

		c.HTML(http.StatusOK, "campaign_new.html", input)
		return
	}

	user, err := h.userService.GetUserByID(input.UserID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	createCampaignInput := request.CreateCampaignInput{}
	createCampaignInput.Name = input.Name
	createCampaignInput.ShortDescription = input.ShortDescription
	createCampaignInput.Description = input.Description
	createCampaignInput.GoalAmount = input.GoalAmount
	createCampaignInput.Perks = input.Perks
	createCampaignInput.User = user

	_, err = h.campaignService.CreateCampaign(createCampaignInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) NewImage(c *gin.Context) {
	idParam := c.Param("id")

	c.HTML(http.StatusOK, "campaign_image.html", gin.H{"ID": idParam})
}

func (h *campaignHandler) CreateImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	idParam := c.Param("id")

	existingCampaign, err := h.campaignService.GetCampaignByID(request.GetCampaignDetailInput{ID: idParam})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	userID := existingCampaign.UserID

	path := fmt.Sprintf("images/%s-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	createCampaignImageInput := request.CreateCampaignImageInput{}
	createCampaignImageInput.CampaignID = idParam
	createCampaignImageInput.IsPrimary = true

	userCampaign, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	createCampaignImageInput.User = userCampaign

	_, err = h.campaignService.SaveCampaignImage(createCampaignImageInput, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) Edit(c *gin.Context) {
	idParam := c.Param("id")

	existingCampaign, err := h.campaignService.GetCampaignByID(request.GetCampaignDetailInput{ID: idParam})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := request.FormUpdateCampaignInput{}
	input.ID = existingCampaign.ID
	input.Name = existingCampaign.Name
	input.ShortDescription = existingCampaign.ShortDescription
	input.Description = existingCampaign.Description
	input.GoalAmount = existingCampaign.GoalAmount
	input.Perks = existingCampaign.Perks

	c.HTML(http.StatusOK, "campaign_edit.html", input)
}

func (h *campaignHandler) Update(c *gin.Context) {
	idParam := c.Param("id")

	var input request.FormUpdateCampaignInput

	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		input.ID = idParam
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	existingCampaign, err := h.campaignService.GetCampaignByID(request.GetCampaignDetailInput{ID: idParam})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	userID := existingCampaign.UserID

	userCampaign, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	updateInput := request.CreateCampaignInput{}
	updateInput.Name = input.Name
	updateInput.ShortDescription = input.ShortDescription
	updateInput.Description = input.Description
	updateInput.GoalAmount = input.GoalAmount
	updateInput.Perks = input.Perks
	updateInput.User = userCampaign

	_, err = h.campaignService.UpdateCampaign(request.GetCampaignDetailInput{ID: idParam}, updateInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}
