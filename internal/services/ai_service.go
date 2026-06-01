package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"one_more_percent/internal/models"
	"os"
	"strings"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type GroqResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// callGroq mengirim full messages list ke Groq API.
// Caller bertanggung jawab menyusun [system, ...history, user].
func callGroq(messages []Message) string {
	token := os.Getenv("MODEL_TOKEN")
	model := os.Getenv("MODEL_NAME")

	if token == "" {
		return "MODEL_TOKEN belum diisi"
	}
	if model == "" {
		model = "llama-3.3-70b-versatile"
	}

	body := GroqRequest{
		Model:    model,
		Messages: messages,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "marshal error"
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "request error"
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return "gagal connect ke model"
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "read response error"
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("groq error: %s", string(respBody))
	}

	var result GroqResponse
	if err = json.Unmarshal(respBody, &result); err != nil {
		return "json parse error"
	}

	if len(result.Choices) == 0 {
		return "model tidak memberi response"
	}

	return result.Choices[0].Message.Content
}

// buildSystemPrompt menyusun system prompt dengan profil user + konteks jadwal aktif.
func buildSystemPrompt(scheduleCtx string) string {
	base := `Kamu adalah One More Percent, asisten personal Ambatukam.

Profil user:
- Nama: Ambatukam
- Status: Mahasiswa yang sedang mengejar kerja remote
- Tujuan: Belajar rutin, bangun skill, dapet kerja remote

Gaya ngobrol:
- Santai seperti teman, bukan robot
- Singkat dan padat
- Suportif tapi ga lebay
- Bahasa campur Indonesia + sedikit gaul oke`

	if scheduleCtx != "" {
		base += "\n\n" + scheduleCtx
	}
	return base
}

// getActiveScheduleContext mengembalikan string konteks jadwal yang sedang berjalan sekarang.
func getActiveScheduleContext() string {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.UTC
	}
	now := time.Now().In(loc)
	dayName := now.Weekday().String()
	currentHHMM := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())

	schedule, err := GetActiveSchedule(dayName, currentHHMM)
	if err != nil || schedule == nil {
		return ""
	}
	return fmt.Sprintf(
		"[Konteks: Sekarang jadwal aktif adalah '%s' (%s - %s). User seharusnya sedang mengerjakan ini.]",
		schedule.Activity, schedule.StartTime, schedule.EndTime,
	)
}

// AskAI menangani percakapan bebas dengan history + konteks jadwal aktif.
func AskAI(chatID int64, message string) string {
	scheduleCtx := getActiveScheduleContext()
	systemPrompt := buildSystemPrompt(scheduleCtx)

	// Susun: [system] + [history] + [pesan baru user]
	msgs := []Message{{Role: "system", Content: systemPrompt}}
	msgs = append(msgs, GetHistory(chatID)...)
	msgs = append(msgs, Message{Role: "user", Content: message})

	reply := callGroq(msgs)

	// Simpan ke history
	AddMessage(chatID, "user", message)
	AddMessage(chatID, "assistant", reply)

	return reply
}

// DetectCompletion menggunakan AI untuk mendeteksi apakah user menyatakan tugas selesai.
func DetectCompletion(message string) bool {
	msgs := []Message{
		{
			Role:    "system",
			Content: "Kamu adalah intent detector. Jawab HANYA dengan kata YES atau NO. Tidak ada kata lain.",
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`Pesan user: "%s"

Apakah user menyatakan bahwa tugas atau jadwal sudah selesai dikerjakan?

Jawab hanya:
YES
atau
NO`, message),
		},
	}

	result := callGroq(msgs)
	result = strings.TrimSpace(strings.ToUpper(result))
	return strings.HasPrefix(result, "YES")
}

// GenerateReminder membuat pesan reminder singkat dan kasual untuk satu aktivitas.
func GenerateReminder(activity string) string {
	msgs := []Message{
		{
			Role: "system",
			Content: `Kamu adalah One More Percent, teman belajar Ambatukam yang santai dan suportif.
Buat reminder singkat dan kasual. Tidak lebay. Maksimal 2-3 kalimat pendek.`,
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Buat reminder untuk jadwal: %s", activity),
		},
	}
	return callGroq(msgs)
}

// GenerateRecap membuat recap harian yang dikirim saat tengah malam.
func GenerateRecap(completed []models.ScheduleProgress, missed []models.ScheduleProgress) string {
	var completedLines, missedLines []string
	for _, p := range completed {
		completedLines = append(completedLines, "- "+p.Activity)
	}
	for _, p := range missed {
		missedLines = append(missedLines, "- "+p.Activity)
	}

	completedStr := "tidak ada"
	if len(completedLines) > 0 {
		completedStr = strings.Join(completedLines, "\n")
	}
	missedStr := "tidak ada"
	if len(missedLines) > 0 {
		missedStr = strings.Join(missedLines, "\n")
	}

	msgs := []Message{
		{
			Role: "system",
			Content: `Kamu adalah One More Percent.
Buat recap harian singkat untuk Ambatukam. Santai seperti teman. Sedikit suportif. Tidak menghakimi.
Gunakan emoji ✅ untuk selesai dan ❌ untuk tidak selesai. Maksimal 6-8 baris.`,
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`Selesai:
%s

Tidak selesai (missed):
%s

Total: %d dari %d jadwal selesai.`,
				completedStr,
				missedStr,
				len(completed),
				len(completed)+len(missed),
			),
		},
	}
	return callGroq(msgs)
}

// GenerateCompletionReply membuat respon singkat setelah user konfirmasi selesai.
func GenerateCompletionReply(activity string) string {
	msgs := []Message{
		{
			Role:    "system",
			Content: "Kamu adalah One More Percent. Singkat, santai, kasual. Maksimal 1-2 kalimat.",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("User baru selesai mengerjakan: %s. Beri respon singkat yang suportif.", activity),
		},
	}
	return callGroq(msgs)
}