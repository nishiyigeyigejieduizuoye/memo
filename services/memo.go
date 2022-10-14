package services

import (
	"github.com/gin-gonic/gin"
	"memo/model"
	"net/http"
	"strconv"
)

type MemoForm struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func memoEndpoints(g *gin.RouterGroup) {
	ng := g.Group("/memo", AuthRequired())

	ng.GET("/", listMemo)
	ng.POST("/", createMemo)

	each := ng.Group("/:id", queryMemo)
	each.GET("", func(ctx *gin.Context) {
		memo := ctx.MustGet("memo").(model.Memo)
		ctx.JSON(http.StatusOK, memo.Detail())
	})
	each.DELETE("", func(ctx *gin.Context) {
		memo := ctx.MustGet("memo").(model.Memo)
		db.Delete(&memo)
		ctx.Status(http.StatusNoContent)
	})
	each.PATCH("", func(ctx *gin.Context) {
		memo := ctx.MustGet("memo").(model.Memo)
		form := MemoForm{}
		if err := ctx.BindJSON(&form); err != nil {
			return
		}
		memo.Title = form.Title
		memo.Content = form.Content
		db.Save(&memo)
		ctx.JSON(http.StatusOK, memo.Detail())
	})
}

func createMemo(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)
	form := MemoForm{}
	if err := ctx.BindJSON(&form); err != nil {
		return
	}
	memo := model.Memo{
		UserID:  user.ID,
		Title:   form.Title,
		Content: form.Content,
	}
	db.Save(&memo)
	ctx.Status(http.StatusCreated)
}

func listMemo(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)

	var memos []model.Memo
	db.Where("user_id=?", user.ID).Order("created_at DESC").First(&memos)

	res := make([]model.MemoInfo, 0, len(memos))
	for _, v := range memos {
		res = append(res, v.Info())
	}
	ctx.JSON(http.StatusOK, res)
}

func queryMemo(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		ctx.Abort()
		return
	}
	memo := model.Memo{}
	err = db.Where("user_id=?", user.ID).Where("id=?", id).First(&memo).Error
	if err != nil {
		ctx.Status(http.StatusNotFound)
		ctx.Abort()
		return
	}
	ctx.Set("memo", memo)
}
