// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package router

import (
	"github.com/jinzhu/gorm"
	"github.com/leandro-lugaresi/hub"
	"github.com/traPtitech/traQ/repository"
	"github.com/traPtitech/traQ/router/oauth2"
	"github.com/traPtitech/traQ/router/session"
	"github.com/traPtitech/traQ/router/utils"
	"github.com/traPtitech/traQ/router/v1"
	"github.com/traPtitech/traQ/router/v3"
	"github.com/traPtitech/traQ/service"
	"github.com/traPtitech/traQ/utils/message"
	"go.uber.org/zap"
)

// Injectors from router_wire.go:

func newRouter(hub2 *hub.Hub, db *gorm.DB, repo repository.Repository, ss *service.Services, logger *zap.Logger, config *Config) *Router {
	manager := ss.ChannelManager
	echo := newEcho(logger, config, repo, manager)
	store := session.NewGormStore(db)
	rbac := ss.RBAC
	onlineCounter := ss.OnlineCounter
	viewerManager := ss.ViewerManager
	processor := ss.Imaging
	messageManager := ss.MessageManager
	fileManager := ss.FileManager
	replaceMapper := utils.NewReplaceMapper(repo, manager)
	replacer := message.NewReplacer(replaceMapper)
	handlers := &v1.Handlers{
		RBAC:           rbac,
		Repo:           repo,
		Hub:            hub2,
		Logger:         logger,
		OC:             onlineCounter,
		VM:             viewerManager,
		Imaging:        processor,
		SessStore:      store,
		ChannelManager: manager,
		MessageManager: messageManager,
		FileManager:    fileManager,
		Replacer:       replacer,
	}
	streamer := ss.WS
	webrtcv3Manager := ss.WebRTCv3
	engine := ss.Search
	v3Config := provideV3Config(config)
	v3Handlers := &v3.Handlers{
		RBAC:           rbac,
		Repo:           repo,
		WS:             streamer,
		Hub:            hub2,
		Logger:         logger,
		OC:             onlineCounter,
		VM:             viewerManager,
		WebRTC:         webrtcv3Manager,
		Imaging:        processor,
		SessStore:      store,
		SearchEngine:   engine,
		ChannelManager: manager,
		MessageManager: messageManager,
		FileManager:    fileManager,
		Replacer:       replacer,
		Config:         v3Config,
	}
	oauth2Config := provideOAuth2Config(config)
	handler := &oauth2.Handler{
		RBAC:      rbac,
		Repo:      repo,
		Logger:    logger,
		SessStore: store,
		Config:    oauth2Config,
	}
	router := &Router{
		e:         echo,
		sessStore: store,
		v1:        handlers,
		v3:        v3Handlers,
		oauth2:    handler,
	}
	return router
}
