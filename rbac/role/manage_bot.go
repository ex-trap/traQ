package role

import (
	"github.com/traPtitech/traQ/rbac"
	"github.com/traPtitech/traQ/rbac/permission"
)

// ManageBot BotConsoleロール
const ManageBot = "manage_bot"

var manageBotPerms = []rbac.Permission{
	permission.GetChannel,
	permission.GetUser,
	permission.GetMe,
	permission.GetWebhook,
	permission.CreateWebhook,
	permission.EditWebhook,
	permission.DeleteWebhook,
	permission.GetBot,
	permission.CreateBot,
	permission.EditBot,
	permission.DeleteBot,
	permission.ReissueBotToken,
	permission.InstallBot,
	permission.UninstallBot,
}