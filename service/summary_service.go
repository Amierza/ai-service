package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Amierza/ai-service/jwt"
	pb "github.com/Amierza/ai-service/proto"
	"github.com/Amierza/ai-service/repository"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type (
	ISummaryService interface {
		GenerateSummary(ctx context.Context, req *pb.SummaryRequest) (string, error)
	}

	summaryService struct {
		summaryRepo repository.ISummaryRepository
		logger      *zap.Logger
		jwt         jwt.IJWT
		openaiKey   string
		openai      *openai.Client
	}
)

// NewSummaryService membuat instance summaryService dan inisialisasi klien OpenAI
func NewSummaryService(summaryRepo repository.ISummaryRepository, logger *zap.Logger, jwt jwt.IJWT, openaiKey string) *summaryService {
	client := openai.NewClient(openaiKey)

	return &summaryService{
		summaryRepo: summaryRepo,
		logger:      logger,
		jwt:         jwt,
		openaiKey:   openaiKey,
		openai:      client,
	}
}

// GenerateSummary membuat ringkasan dari percakapan bimbingan
func (ss *summaryService) GenerateSummary(ctx context.Context, req *pb.SummaryRequest) (string, error) {
	if req == nil || req.Task == nil {
		ss.logger.Error("invalid request: task is nil")
		return "", errors.New("invalid request: task is nil")
	}

	task := req.Task
	ss.logger.Info("Generating summary for thesis",
		zap.String("session_id", task.SessionId),
		zap.String("title", task.ThesisInfo.Title),
	)

	// 1️⃣ Ambil semua pesan dari task.Messages
	if len(task.Messages) == 0 {
		return "", errors.New("no messages found to summarize")
	}

	var messages []string
	for _, msg := range task.Messages {
		if msg.IsText {
			sender := msg.Sender.Name
			messages = append(messages, fmt.Sprintf("%s: %s", sender, msg.Text))
		}
	}

	// 2️⃣ Gabungkan semua percakapan menjadi satu string
	joinedMessages := strings.Join(messages, "\n")

	// 3️⃣ Buat prompt untuk GPT
	prompt := fmt.Sprintf(`
Buatlah ringkasan akademik dari percakapan bimbingan berikut ini dalam bentuk poin-poin yang ringkas dan jelas.

Fokuskan pada hal-hal berikut:
1. Inti pembahasan utama.
2. Arahan atau masukan dari dosen.
3. Tindak lanjut atau progres mahasiswa.

Judul Tugas Akhir: %s
Deskripsi: %s

Percakapan:
%s

Tuliskan hasil ringkasan dalam Bahasa Indonesia yang formal, dalam format poin-poin (gunakan tanda '-' atau '•' di awal setiap poin).
`, task.ThesisInfo.Title, task.ThesisInfo.Description, joinedMessages)

	// 4️⃣ Panggil API OpenAI
	resp, err := ss.openai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini, // model cepat dan hemat
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "Kamu adalah asisten akademik yang membantu merangkum diskusi mahasiswa dan dosen.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 500,
	})
	if err != nil {
		ss.logger.Error("failed to call OpenAI API", zap.Error(err))
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no response from OpenAI")
	}

	summary := strings.TrimSpace(resp.Choices[0].Message.Content)
	if summary == "" {
		return "", errors.New("empty summary generated")
	}

	err = ss.summaryRepo.SaveSummary(ctx, nil, task.SessionId, summary)
	if err != nil {
		ss.logger.Error("failed to save summary to repository", zap.Error(err))
		return "", err
	}

	err = ss.summaryRepo.UpdateStatusSessionFinished(ctx, nil, req.Task.SessionId)
	if err != nil {
		ss.logger.Error("failed update session status to finished", zap.Error(err))
		return "", err
	}

	ss.logger.Info("Summary generated successfully", zap.String("session_id", task.SessionId))
	return summary, nil
}
