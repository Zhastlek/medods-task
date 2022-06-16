package handlers

import (
	"log"
	"medods/internal/model"
	"medods/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	Register(router *gin.Engine)
}

type handler struct {
	service service.AuthServiceInterface
}

func NewHandlers(auth service.AuthServiceInterface) HandlerInterface {
	return &handler{
		service: auth,
	}
}

func (h *handler) Register(router *gin.Engine) {
	router.GET("/sign-in/:guid", h.GenerateJWT)
	router.GET("/refresh/:guid", h.Refresh)
}

func (h *handler) GenerateJWT(c *gin.Context) {
	userGUID := c.Param("guid")
	tokens, err := h.service.CreateToken(userGUID)
	if err != nil {
		log.Printf("Error don't generate pair of tokens: %v\n", err)
		c.AbortWithStatusJSON(401, nil)
		return
	}
	h.setCookie(c, tokens)

}

func (h *handler) Refresh(c *gin.Context) {
	userGUID := c.Param("guid")
	refreshToken, err := c.Request.Cookie("refresh-token")
	if err != nil {
		log.Printf("Error don't find refresh token: %v\n", err)
		c.AbortWithStatusJSON(401, gin.H{"error": "Invalid Api token"})
		return
	}
	accessToken, err := c.Request.Cookie("access-token")
	if err != nil {
		log.Printf("Error don't find access-token token: %v\n", err)
		c.AbortWithStatusJSON(401, gin.H{"error": "Invalid Api token"})
		return
	}
	oldTokens := &model.Jwt{
		UserGUID:     userGUID,
		RefreshToken: refreshToken.Value,
		AccsessToken: accessToken.Value,
	}
	tokens, err := h.service.UpdateToken(oldTokens)
	if err != nil {
		log.Printf("Error don't find refresh token in DB: %v\n", err)
		c.AbortWithStatusJSON(401, gin.H{"error": "error"})
		return
	}
	h.setCookie(c, tokens)
}

func (h *handler) setCookie(c *gin.Context, tokens *model.Jwt) {
	accessExpire := time.Now().Add(time.Minute * 10)
	refreshExpire := time.Now().AddDate(0, 0, 10)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access-token",
		Value:    tokens.AccsessToken,
		HttpOnly: true,
		MaxAge:   10 * 24 * 60 * 60,
		Path:     "/",
		Expires:  accessExpire,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh-token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		MaxAge:   10 * 24 * 60 * 60,
		Path:     "/refresh",
		Expires:  refreshExpire,
	})
}
