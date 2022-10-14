package services

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"memo/model"
	"net/http"
	"time"
)

type UserForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (form *UserForm) GetPassword() []byte {
	pw, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return pw
}

func AuthRequired() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sid, err := ctx.Cookie("session")
		if err != nil {
			ctx.String(http.StatusUnauthorized, "unauthorized")
			ctx.Abort()
			return
		}
		session := model.Session{}
		err = db.Preload("User").Where("id=?", sid).First(&session).Error
		if err != nil {
			ctx.String(http.StatusUnauthorized, "unauthorized")
			ctx.Abort()
			return
		}
		session.UserAgent = ctx.GetHeader("User-Agent")
		session.LastAccessed = time.Now()
		go db.Save(&session)
		ctx.Set("user", session.User)
		ctx.Set("session", session)
	}
}

func userEndpoints(g *gin.RouterGroup) {
	ng := g.Group("/user")

	ng.GET("/info", AuthRequired(), func(ctx *gin.Context) {
		user := ctx.MustGet("user").(model.User)
		ctx.JSON(http.StatusOK, user.Info())
	})

	ng.POST("/register", func(ctx *gin.Context) {
		form := UserForm{}
		if err := ctx.BindJSON(&form); err != nil {
			return
		}
		user := model.User{
			Username: form.Username,
			Password: string(form.GetPassword()),
		}
		err := db.Save(&user).Error
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, user.Info())
	})

	ng.POST("/login", func(ctx *gin.Context) {
		form := UserForm{}
		if err := ctx.BindJSON(&form); err != nil {
			return
		}
		user := model.User{}
		err := db.Where("username=?", form.Username).First(&user).Error
		if err != nil {
			ctx.String(http.StatusNotFound, "user not found")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
		if err != nil {
			ctx.String(http.StatusUnauthorized, "password does not match")
			return
		}
		session := model.Session{
			ID:           uuid.NewV4().String(),
			UserID:       user.ID,
			UserAgent:    ctx.GetHeader("User-Agent"),
			LastAccessed: time.Now(),
		}
		db.Save(&session)
		ctx.SetCookie("session", session.ID, 86400*30, "/", "", false, true)
		ctx.Status(http.StatusNoContent)
	})

	ng.POST("/logout", AuthRequired(), func(ctx *gin.Context) {
		session := ctx.MustGet("session").(model.Session)
		db.Delete(&session)
		ctx.SetCookie("session", "", 0, "/", "", false, true)
		ctx.Status(http.StatusNoContent)
	})
}
