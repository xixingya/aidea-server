package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/mylxsw/aidea-server/internal/service"
	"github.com/mylxsw/asteria/log"

	"github.com/mylxsw/aidea-server/api/controllers/common"
	"github.com/mylxsw/aidea-server/config"
	"github.com/mylxsw/aidea-server/internal/repo"
	"github.com/mylxsw/aidea-server/internal/youdao"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/web"
)

type CreativeController struct {
	conf         *config.Config
	translater   youdao.Translater       `autowire:"@"`
	creativeRepo *repo.CreativeRepo      `autowire:"@"`
	gallerySrv   *service.GalleryService `autowire:"@"`
}

func NewCreativeController(resolver infra.Resolver, conf *config.Config) web.Controller {
	ctl := CreativeController{conf: conf}
	resolver.MustAutoWire(&ctl)
	return &ctl
}

func (ctl *CreativeController) Register(router web.Router) {
	router.Group("/creatives", func(router web.Router) {
		router.Get("/gallery", ctl.Gallery)
		router.Get("/gallery/{id}", ctl.GalleryItem)
	})
}

// Gallery 作品图库列表
func (ctl *CreativeController) Gallery(ctx context.Context, webCtx web.Context) web.Response {
	page := webCtx.Int64Input("page", 1)
	if page < 1 || page > 20 {
		page = 1
	}

	pageSize := webCtx.Int64Input("per_page", 20)
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	res, err := ctl.gallerySrv.Gallery(ctx, page, pageSize)
	if err != nil {
		log.WithFields(log.Fields{"page": page, "per_page": pageSize}).Errorf("get gallery list failed: %v", err)
		return webCtx.JSONError(common.Text(webCtx, ctl.translater, common.ErrInternalError), http.StatusInternalServerError)
	}

	return webCtx.JSON(res)
}

// GalleryItem 作品图库详情
func (ctl *CreativeController) GalleryItem(ctx context.Context, webCtx web.Context) web.Response {
	id, err := strconv.Atoi(webCtx.PathVar("id"))
	if err != nil {
		return webCtx.JSONError(common.Text(webCtx, ctl.translater, common.ErrInvalidRequest), http.StatusBadRequest)
	}

	item, err := ctl.creativeRepo.GalleryByID(ctx, int64(id))
	if err != nil {
		return webCtx.JSONError(common.Text(webCtx, ctl.translater, common.ErrInternalError), http.StatusInternalServerError)
	}

	return webCtx.JSON(item)
}
