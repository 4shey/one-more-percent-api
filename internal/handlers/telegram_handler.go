package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"one_more_percent/internal/services"
)

type TelegramWebhook struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

func TelegramWebhookHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var update TelegramWebhook

	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		fmt.Println("decode error:", err)
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	userMessage := update.Message.Text
	chatID := update.Message.Chat.ID

	fmt.Println("==========")
	fmt.Println("User:", userMessage)
	fmt.Println("ChatID:", chatID)

	var reply string

	// Check if there's a pending reminder awaiting user confirmation.
	pending := services.GetPendingReminder(chatID)
	if pending != nil {
		if services.DetectCompletion(userMessage) {
			// Mark the schedule as completed in DB.
			err := services.MarkCompleted(pending.ScheduleID, pending.ProgressDate)
			if err != nil {
				fmt.Println("MarkCompleted error:", err)
			}
			services.ClearPendingReminder(chatID)

			// Generate a short celebratory reply.
			reply = services.GenerateCompletionReply(pending.Activity)
			fmt.Printf("Marked completed: schedule_id=%d activity=%s\n",
				pending.ScheduleID, pending.Activity)
		} else {
			// User replied but didn't confirm completion — treat as normal chat.
			reply = services.AskAI(chatID, userMessage)
		}
	} else {
		// No pending reminder — normal conversation.
		reply = services.AskAI(chatID, userMessage)
	}

	fmt.Println("AI Reply:", reply)

	err = services.SendTelegramMessage(chatID, reply)
	if err != nil {
		fmt.Println("telegram send error:", err)
	} else {
		fmt.Println("message sent")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}