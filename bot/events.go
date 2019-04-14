package bot

import "github.com/traPtitech/traQ/model"

const (
	// Ping Pingイベント
	Ping model.BotEvent = "PING"
	// Joined チャンネル参加イベント
	Joined model.BotEvent = "JOINED"
	// Left チャンネル退出イベント
	Left model.BotEvent = "LEFT"
	// MessageCreated メッセージ作成イベント
	MessageCreated model.BotEvent = "MESSAGE_CREATED"
)

var eventSet = map[model.BotEvent]bool{
	Ping:           true,
	Joined:         true,
	Left:           true,
	MessageCreated: true,
}

// IsEvent 引数の文字列がボットイベントかどうか
func IsEvent(str string) bool {
	return eventSet[model.BotEvent(str)]
}
