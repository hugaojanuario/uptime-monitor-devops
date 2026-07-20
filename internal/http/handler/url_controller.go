package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/uptime-monitor-devops/internal/models"
	"github.com/hugaojanuario/uptime-monitor-devops/internal/services"
	"github.com/hugaojanuario/uptime-monitor-devops/internal/utils"
)

type URLController struct {
	s *services.Service
}

func NewURLController(s *services.Service) *URLController {
	return &URLController{s: s}
}

// CreateURL godoc
//
//	@Summary		Cadastra uma url
//	@Description	Cadastra uma url para ser monitorada
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CreateURLRequest	true	"Dados da url"
//	@Success		201		{object}	models.URL
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		409		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/urls [post]
func (u *URLController) CreateURL(c *gin.Context) {
	var req models.CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body invalid."})
		return
	}

	url, err := u.s.CreateURL(req)
	if err != nil {
		utils.RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, url)
}

// CheckAllURLs godoc
//
//	@Summary		Verifica todas as urls
//	@Description	Faz um GET em todas as urls cadastradas, retorna o status http de cada uma e grava os resultados no arquivo de resultados
//	@Tags			urls
//	@Produce		json
//	@Success		200	{array}		models.CheckResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/urls [get]
func (u *URLController) CheckAllURLs(c *gin.Context) {
	results, err := u.s.CheckAllURLs(c.Request.Context())
	if err != nil {
		utils.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, results)
}

// CheckURLByID godoc
//
//	@Summary		Verifica uma url
//	@Description	Faz um GET apenas na url do id informado, retorna o status http e grava o resultado no arquivo de resultados
//	@Tags			urls
//	@Produce		json
//	@Param			id	path		string	true	"ID da url (uuid)"
//	@Success		200	{object}	models.CheckResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/urls/{id} [get]
func (u *URLController) CheckURLByID(c *gin.Context) {
	id := c.Param("id")

	result, err := u.s.CheckURLByID(c.Request.Context(), id)
	if err != nil {
		utils.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
