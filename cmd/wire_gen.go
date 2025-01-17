// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package cmd

import (
	"github.com/jinzhu/gorm"
	"github.com/leandro-lugaresi/hub"
	"github.com/traPtitech/traQ/repository"
	"github.com/traPtitech/traQ/router"
	"github.com/traPtitech/traQ/service"
	"github.com/traPtitech/traQ/service/bot"
	"github.com/traPtitech/traQ/service/channel"
	"github.com/traPtitech/traQ/service/counter"
	"github.com/traPtitech/traQ/service/exevent"
	"github.com/traPtitech/traQ/service/file"
	"github.com/traPtitech/traQ/service/imaging"
	"github.com/traPtitech/traQ/service/message"
	"github.com/traPtitech/traQ/service/notification"
	"github.com/traPtitech/traQ/service/rbac"
	"github.com/traPtitech/traQ/service/viewer"
	"github.com/traPtitech/traQ/service/webrtcv3"
	"github.com/traPtitech/traQ/service/ws"
	"github.com/traPtitech/traQ/utils/storage"
	"go.uber.org/zap"
)

import (
	_ "net/http/pprof"
)

// Injectors from serve_wire.go:

func newServer(hub2 *hub.Hub, db *gorm.DB, repo repository.Repository, fs storage.FileStorage, logger *zap.Logger, c2 *Config) (*Server, error) {
	manager, err := channel.InitChannelManager(repo, logger)
	if err != nil {
		return nil, err
	}
	botService := bot.NewService(repo, manager, hub2, logger)
	onlineCounter := counter.NewOnlineCounter(hub2)
	unreadMessageCounter, err := counter.NewUnreadMessageCounter(db, hub2)
	if err != nil {
		return nil, err
	}
	messageCounter, err := counter.NewMessageCounter(db, hub2)
	if err != nil {
		return nil, err
	}
	channelCounter, err := counter.NewChannelCounter(db, hub2)
	if err != nil {
		return nil, err
	}
	messageManager, err := message.NewMessageManager(repo, manager, logger)
	if err != nil {
		return nil, err
	}
	stampThrottler := exevent.NewStampThrottler(hub2, messageManager)
	firebaseCredentialsFilePathString := provideFirebaseCredentialsFilePathString(c2)
	client, err := newFCMClientIfAvailable(repo, logger, unreadMessageCounter, firebaseCredentialsFilePathString)
	if err != nil {
		return nil, err
	}
	config := provideImageProcessorConfig(c2)
	processor := imaging.NewProcessor(config)
	fileManager, err := file.InitFileManager(repo, fs, processor, logger)
	if err != nil {
		return nil, err
	}
	viewerManager := viewer.NewManager(hub2)
	webrtcv3Manager := webrtcv3.NewManager(hub2)
	streamer := ws.NewStreamer(hub2, viewerManager, webrtcv3Manager, logger)
	serverOriginString := provideServerOriginString(c2)
	notificationService := notification.NewService(repo, manager, messageManager, fileManager, hub2, logger, client, streamer, viewerManager, serverOriginString)
	rbacRBAC, err := rbac.New(db)
	if err != nil {
		return nil, err
	}
	esEngineConfig := provideESEngineConfig(c2)
	engine, err := initSearchServiceIfAvailable(messageManager, manager, repo, logger, esEngineConfig)
	if err != nil {
		return nil, err
	}
	services := &service.Services{
		BOT:                  botService,
		ChannelManager:       manager,
		OnlineCounter:        onlineCounter,
		UnreadMessageCounter: unreadMessageCounter,
		MessageCounter:       messageCounter,
		ChannelCounter:       channelCounter,
		StampThrottler:       stampThrottler,
		FCM:                  client,
		FileManager:          fileManager,
		Imaging:              processor,
		MessageManager:       messageManager,
		Notification:         notificationService,
		RBAC:                 rbacRBAC,
		Search:               engine,
		ViewerManager:        viewerManager,
		WebRTCv3:             webrtcv3Manager,
		WS:                   streamer,
	}
	routerConfig := provideRouterConfig(c2)
	echo := router.Setup(hub2, db, repo, services, logger, routerConfig)
	server := &Server{
		L:      logger,
		SS:     services,
		Router: echo,
		Hub:    hub2,
		Repo:   repo,
	}
	return server, nil
}
