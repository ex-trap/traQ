package v1

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traQ/model"
	"github.com/traPtitech/traQ/service/rbac/permission"
	"github.com/traPtitech/traQ/utils/optional"
	"time"
)

type meResponse struct {
	UserID      uuid.UUID               `json:"userId"`
	Name        string                  `json:"name"`
	DisplayName string                  `json:"displayName"`
	IconID      uuid.UUID               `json:"iconFileId"`
	Bot         bool                    `json:"bot"`
	TwitterID   string                  `json:"twitterId"`
	LastOnline  optional.Time           `json:"lastOnline"`
	IsOnline    bool                    `json:"isOnline"`
	Suspended   bool                    `json:"suspended"`
	Status      int                     `json:"accountStatus"`
	Role        string                  `json:"role"`
	Permissions []permission.Permission `json:"permissions"`
}

func (h *Handlers) formatMe(user model.UserInfo) *meResponse {
	res := &meResponse{
		UserID:      user.GetID(),
		Name:        user.GetName(),
		DisplayName: user.GetResponseDisplayName(),
		IconID:      user.GetIconFileID(),
		Bot:         user.IsBot(),
		TwitterID:   user.GetTwitterID(),
		IsOnline:    h.OC.IsOnline(user.GetID()),
		Suspended:   user.GetState() != model.UserAccountStatusActive,
		Status:      user.GetState().Int(),
		Role:        user.GetRole(),
		Permissions: h.RBAC.GetGrantedPermissions(user.GetRole()),
	}

	if res.IsOnline {
		res.LastOnline = optional.TimeFrom(time.Now())
	} else {
		res.LastOnline = user.GetLastOnline()
	}
	return res
}

type userResponse struct {
	UserID      uuid.UUID     `json:"userId"`
	Name        string        `json:"name"`
	DisplayName string        `json:"displayName"`
	IconID      uuid.UUID     `json:"iconFileId"`
	Bot         bool          `json:"bot"`
	TwitterID   string        `json:"twitterId"`
	LastOnline  optional.Time `json:"lastOnline"`
	IsOnline    bool          `json:"isOnline"`
	Suspended   bool          `json:"suspended"`
	Status      int           `json:"accountStatus"`
}

func (h *Handlers) formatUser(user model.UserInfo) *userResponse {
	res := &userResponse{
		UserID:      user.GetID(),
		Name:        user.GetName(),
		DisplayName: user.GetResponseDisplayName(),
		IconID:      user.GetIconFileID(),
		Bot:         user.IsBot(),
		TwitterID:   user.GetTwitterID(),
		IsOnline:    h.OC.IsOnline(user.GetID()),
		Suspended:   user.GetState() != model.UserAccountStatusActive,
		Status:      user.GetState().Int(),
	}

	if res.IsOnline {
		res.LastOnline = optional.TimeFrom(time.Now())
	} else {
		res.LastOnline = user.GetLastOnline()
	}
	return res
}

func (h *Handlers) formatUsers(users []model.UserInfo) []*userResponse {
	res := make([]*userResponse, len(users))
	for i, user := range users {
		res[i] = h.formatUser(user)
	}
	return res
}

type userDetailResponse struct {
	UserID      uuid.UUID      `json:"userId"`
	Name        string         `json:"name"`
	DisplayName string         `json:"displayName"`
	IconID      uuid.UUID      `json:"iconFileId"`
	Bot         bool           `json:"bot"`
	TwitterID   string         `json:"twitterId"`
	LastOnline  optional.Time  `json:"lastOnline"`
	IsOnline    bool           `json:"isOnline"`
	Suspended   bool           `json:"suspended"`
	Status      int            `json:"accountStatus"`
	TagList     []*tagResponse `json:"tagList"`
}

func (h *Handlers) formatUserDetail(user model.UserInfo, tagList []model.UserTag) (*userDetailResponse, error) {
	res := &userDetailResponse{
		UserID:      user.GetID(),
		Name:        user.GetName(),
		DisplayName: user.GetResponseDisplayName(),
		IconID:      user.GetIconFileID(),
		Bot:         user.IsBot(),
		TwitterID:   user.GetTwitterID(),
		IsOnline:    h.OC.IsOnline(user.GetID()),
		Suspended:   user.GetState() != model.UserAccountStatusActive,
		Status:      user.GetState().Int(),
		TagList:     formatTags(tagList),
	}

	if res.IsOnline {
		res.LastOnline = optional.TimeFrom(time.Now())
	} else {
		res.LastOnline = user.GetLastOnline()
	}
	return res, nil
}

type webhookResponse struct {
	WebhookID   string    `json:"webhookId"`
	BotUserID   string    `json:"botUserId"`
	DisplayName string    `json:"displayName"`
	Description string    `json:"description"`
	Secure      bool      `json:"secure"`
	ChannelID   string    `json:"channelId"`
	CreatorID   string    `json:"creatorId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func formatWebhook(w model.Webhook) *webhookResponse {
	return &webhookResponse{
		WebhookID:   w.GetID().String(),
		BotUserID:   w.GetBotUserID().String(),
		DisplayName: w.GetName(),
		Description: w.GetDescription(),
		Secure:      len(w.GetSecret()) > 0,
		ChannelID:   w.GetChannelID().String(),
		CreatorID:   w.GetCreatorID().String(),
		CreatedAt:   w.GetCreatedAt(),
		UpdatedAt:   w.GetUpdatedAt(),
	}
}

func formatWebhooks(ws []model.Webhook) []*webhookResponse {
	res := make([]*webhookResponse, len(ws))
	for i, w := range ws {
		res[i] = formatWebhook(w)
	}
	return res
}

type channelResponse struct {
	ChannelID  string      `json:"channelId"`
	Name       string      `json:"name"`
	Parent     string      `json:"parent"`
	Topic      string      `json:"topic"`
	Children   []uuid.UUID `json:"children"`
	Member     []uuid.UUID `json:"member"`
	Visibility bool        `json:"visibility"`
	Force      bool        `json:"force"`
	Private    bool        `json:"private"`
	DM         bool        `json:"dm"`
}

func (h *Handlers) formatChannel(channel *model.Channel) (response *channelResponse, err error) {
	response = &channelResponse{
		ChannelID:  channel.ID.String(),
		Name:       channel.Name,
		Topic:      channel.Topic,
		Children:   channel.ChildrenID,
		Visibility: channel.IsVisible,
		Force:      channel.IsForced,
		Private:    !channel.IsPublic,
		DM:         channel.IsDMChannel(),
		Member:     make([]uuid.UUID, 0),
	}
	if channel.ParentID != uuid.Nil {
		response.Parent = channel.ParentID.String()
	}

	if response.Private {
		response.Member, err = h.ChannelManager.GetDMChannelMembers(channel.ID)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

type botResponse struct {
	BotID           uuid.UUID           `json:"botId"`
	BotUserID       uuid.UUID           `json:"botUserId"`
	Description     string              `json:"description"`
	SubscribeEvents model.BotEventTypes `json:"subscribeEvents"`
	State           model.BotState      `json:"state"`
	CreatorID       uuid.UUID           `json:"creatorId"`
	CreatedAt       time.Time           `json:"createdAt"`
	UpdatedAt       time.Time           `json:"updatedAt"`
}

func formatBot(b *model.Bot) *botResponse {
	return &botResponse{
		BotID:           b.ID,
		BotUserID:       b.BotUserID,
		Description:     b.Description,
		SubscribeEvents: b.SubscribeEvents,
		State:           b.State,
		CreatorID:       b.CreatorID,
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       b.UpdatedAt,
	}
}

func formatBots(bs []*model.Bot) []*botResponse {
	res := make([]*botResponse, len(bs))
	for i, b := range bs {
		res[i] = formatBot(b)
	}
	return res
}

type botDetailResponse struct {
	BotID            uuid.UUID           `json:"botId"`
	BotUserID        uuid.UUID           `json:"botUserId"`
	Description      string              `json:"description"`
	SubscribeEvents  model.BotEventTypes `json:"subscribeEvents"`
	State            model.BotState      `json:"state"`
	CreatorID        uuid.UUID           `json:"creatorId"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
	VerificationCode string              `json:"verificationCode"`
	AccessToken      string              `json:"accessToken"`
	PostURL          string              `json:"postUrl"`
	Privileged       bool                `json:"privileged"`
	BotCode          string              `json:"botCode"`
}

func formatBotDetail(b *model.Bot, t *model.OAuth2Token) *botDetailResponse {
	return &botDetailResponse{
		BotID:            b.ID,
		BotUserID:        b.BotUserID,
		Description:      b.Description,
		SubscribeEvents:  b.SubscribeEvents,
		State:            b.State,
		CreatorID:        b.CreatorID,
		CreatedAt:        b.CreatedAt,
		UpdatedAt:        b.UpdatedAt,
		VerificationCode: b.VerificationToken,
		AccessToken:      t.AccessToken,
		PostURL:          b.PostURL,
		Privileged:       b.Privileged,
		BotCode:          b.BotCode,
	}
}

type tagResponse struct {
	ID        uuid.UUID `json:"tagId"`
	Tag       string    `json:"tag"`
	IsLocked  bool      `json:"isLocked"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func formatTag(ut model.UserTag) *tagResponse {
	return &tagResponse{
		ID:        ut.GetTagID(),
		Tag:       ut.GetTag(),
		IsLocked:  ut.GetIsLocked(),
		CreatedAt: ut.GetCreatedAt(),
		UpdatedAt: ut.GetUpdatedAt(),
	}
}

func formatTags(uts []model.UserTag) []*tagResponse {
	res := make([]*tagResponse, len(uts))
	for i, ut := range uts {
		res[i] = formatTag(ut)
	}
	return res
}

type userGroupResponse struct {
	GroupID     uuid.UUID   `json:"groupId"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	AdminUserID uuid.UUID   `json:"adminUserId"`
	Members     []uuid.UUID `json:"members"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

func formatUserGroup(g *model.UserGroup) *userGroupResponse {
	r := &userGroupResponse{
		GroupID:     g.ID,
		Name:        g.Name,
		Description: g.Description,
		Type:        g.Type,
		AdminUserID: g.Admins[0].UserID,
		Members:     make([]uuid.UUID, 0),
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
	for _, member := range g.Members {
		r.Members = append(r.Members, member.UserID)
	}
	return r
}

func (h *Handlers) formatUserGroups(gs []*model.UserGroup) ([]*userGroupResponse, error) {
	arr := make([]*userGroupResponse, len(gs))
	for i, g := range gs {
		arr[i] = formatUserGroup(g)
	}
	return arr, nil
}

type fileResponse struct {
	FileID      uuid.UUID `json:"fileId"`
	Name        string    `json:"name"`
	Mime        string    `json:"mime"`
	Size        int64     `json:"size"`
	MD5         string    `json:"md5"`
	HasThumb    bool      `json:"hasThumb"`
	ThumbWidth  int       `json:"thumbWidth,omitempty"`
	ThumbHeight int       `json:"thumbHeight,omitempty"`
	Datetime    time.Time `json:"datetime"`
}

func formatFile(f model.File) *fileResponse {
	hasThumb, t := f.GetThumbnail(model.ThumbnailTypeImage)
	return &fileResponse{
		FileID:      f.GetID(),
		Name:        f.GetFileName(),
		Mime:        f.GetMIMEType(),
		Size:        f.GetFileSize(),
		MD5:         f.GetMD5Hash(),
		HasThumb:    hasThumb,
		ThumbWidth:  t.Width,
		ThumbHeight: t.Height,
		Datetime:    f.GetCreatedAt(),
	}
}
