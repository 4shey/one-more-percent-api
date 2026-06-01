package services

import "sync"

const maxHistoryMessages = 10 // simpan 10 pesan terakhir (5 exchange)

var (
	histMu  sync.Mutex
	history = map[int64][]Message{}
)

// AddMessage menyimpan pesan ke history chat user.
func AddMessage(chatID int64, role, content string) {
	histMu.Lock()
	defer histMu.Unlock()

	msgs := history[chatID]
	msgs = append(msgs, Message{Role: role, Content: content})

	// Buang pesan paling lama jika sudah melebihi batas.
	if len(msgs) > maxHistoryMessages {
		msgs = msgs[len(msgs)-maxHistoryMessages:]
	}

	history[chatID] = msgs
}

// GetHistory mengembalikan salinan history chat untuk satu chatID.
func GetHistory(chatID int64) []Message {
	histMu.Lock()
	defer histMu.Unlock()

	src := history[chatID]
	result := make([]Message, len(src))
	copy(result, src)
	return result
}
