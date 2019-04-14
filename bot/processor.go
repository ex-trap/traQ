package bot

import (
	"bytes"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo"
	"github.com/leandro-lugaresi/hub"
	"github.com/traPtitech/traQ/event"
	"github.com/traPtitech/traQ/model"
	"github.com/traPtitech/traQ/repository"
	"github.com/traPtitech/traQ/utils/message"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

const (
	headerTRAQBotEvent             = "X-TRAQ-BOT-EVENT"
	headerTRAQBotRequestID         = "X-TRAQ-BOT-REQUEST-ID"
	headerTRAQBotVerificationToken = "X-TRAQ-BOT-TOKEN"
)

// Processor ボットプロセッサー
type Processor struct {
	repo    repository.Repository
	logger  *zap.Logger
	hub     *hub.Hub
	bufPool sync.Pool
	client  http.Client
}

// NewProcessor ボットプロセッサーを生成し、起動します
func NewProcessor(repo repository.Repository, hub *hub.Hub, logger *zap.Logger) *Processor {
	p := &Processor{
		repo:   repo,
		logger: logger,
		hub:    hub,
		bufPool: sync.Pool{
			New: func() interface{} { return &bytes.Buffer{} },
		},
		client: http.Client{
			Timeout:       5 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
		},
	}
	go func() {
		sub := hub.Subscribe(10, event.MessageCreated)
		for ev := range sub.Receiver {
			m := ev.Fields["message"].(*model.Message)
			e := ev.Fields["embedded"].([]*message.EmbeddedInfo)
			plain := ev.Fields["plain"].(string)
			go p.createMessageHandler(m, e, plain)
		}
	}()
	go func() {
		sub := hub.Subscribe(1, event.BotPingRequest)
		for ev := range sub.Receiver {
			botID := ev.Fields["bot_id"].(uuid.UUID)
			bot, err := repo.GetBotByID(botID)
			if err != nil {
				logger.Error("failed to GetBotByID", zap.Error(err), zap.Stringer("bot_id", botID))
				continue
			}
			p.pingHandler(bot)
		}
	}()
	go func() {
		sub := hub.Subscribe(10, event.BotJoined, event.BotLeft)
		for ev := range sub.Receiver {
			botID := ev.Fields["bot_id"].(uuid.UUID)
			chId := ev.Fields["channel_id"].(uuid.UUID)
			switch ev.Name {
			case event.BotJoined:
				go p.joinedAndLeftHandler(botID, chId, Joined)
			case event.BotLeft:
				go p.joinedAndLeftHandler(botID, chId, Left)
			}
		}
	}()
	return p
}

func (p *Processor) sendEvent(b *model.Bot, event model.BotEvent, body []byte) (ok bool) {
	reqID := uuid.Must(uuid.NewV4())

	req, _ := http.NewRequest(http.MethodPost, b.PostURL, bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	req.Header.Set(headerTRAQBotEvent, event.String())
	req.Header.Set(headerTRAQBotRequestID, reqID.String())
	req.Header.Set(headerTRAQBotVerificationToken, b.VerificationToken)

	res, err := p.client.Do(req)
	if err != nil {
		return false
	}
	_ = res.Body.Close()
	return res.StatusCode == http.StatusNoContent
}

func (p *Processor) makePayloadJSON(payload interface{}) (b []byte, releaseFunc func(), err error) {
	buf := p.bufPool.Get().(*bytes.Buffer)
	releaseFunc = func() {
		buf.Reset()
		p.bufPool.Put(buf)
	}

	if err := json.NewEncoder(buf).Encode(&payload); err != nil {
		releaseFunc()
		return nil, nil, err
	}

	return buf.Bytes(), releaseFunc, nil
}
