package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cyberlab/backend/internal/models"
	"github.com/cyberlab/backend/internal/repository"
	"github.com/cyberlab/backend/pkg/utils"
	"go.uber.org/zap"
)

type AIService struct {
	config  AIServiceConfig
	aiRepo  *repository.AIRepository
	logRepo *repository.LogRepository
}

type AIServiceConfig struct {
	Provider    string
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float64
}

func NewAIService(cfg AIServiceConfig, aiRepo *repository.AIRepository, logRepo *repository.LogRepository) *AIService {
	return &AIService{
		config:  cfg,
		aiRepo:  aiRepo,
		logRepo: logRepo,
	}
}

type HintRequest struct {
	ChallengeID uint   `json:"challengeId"`
	Question    string `json:"question"`
}

type HintResponse struct {
	Answer string `json:"answer"`
}

type Vulnerability struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
}

type AuditResponse struct {
	RiskCount      int             `json:"riskCount"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

func (s *AIService) GetHint(userID uint, req *HintRequest) (*HintResponse, error) {
	challenge, err := s.aiRepo.GetChallengeDetail(req.ChallengeID)
	if err != nil {
		return nil, errors.New("challenge not found")
	}

	prompt := fmt.Sprintf(`You are a CTF challenge assistant. The user is working on a cybersecurity challenge.

Challenge: %s
Category: %s
Difficulty: %s

User's question: %s

Provide a helpful hint that guides the user toward the solution without giving away the flag directly. Focus on the methodology and thought process.`,
		challenge.Title, challenge.Category, challenge.Difficulty, req.Question)

	answer, err := s.callAI(prompt)
	if err != nil {
		utils.Log.Warn("AI hint service unavailable, using fallback", zap.Error(err))
		return &HintResponse{
			Answer: "AI hint service is currently unavailable. Please try again later or refer to the challenge description for more context.",
		}, nil
	}

	// Save hint record
	s.aiRepo.CreateHint(&models.AIHint{
		UserID:      userID,
		ChallengeID: req.ChallengeID,
		Question:    req.Question,
		Answer:      answer,
	})

	s.logRepo.Create(&models.SystemLog{
		UserID: userID,
		Action: models.ActionAIRequest,
	})

	return &HintResponse{Answer: answer}, nil
}

func (s *AIService) AuditCode(userID uint, fileContent string, filename string) (*AuditResponse, error) {
	prompt := fmt.Sprintf(`You are a code security auditor. Analyze the following source code for security vulnerabilities.

Filename: %s

Code:
%s

Identify all security vulnerabilities. For each, provide:
1. Type of vulnerability (SQLi, XSS, RCE, SSRF, File Upload, etc.)
2. Severity (CRITICAL, HIGH, MEDIUM, LOW)
3. File name
4. Line number (approximate)
5. Description of the issue
6. Fix suggestion

Return as a JSON array.`,
		filename, fileContent)

	result, err := s.callAI(prompt)
	if err != nil {
		return nil, errors.New("AI audit service unavailable")
	}

	// Try to parse AI response as JSON
	var vulnerabilities []Vulnerability
	if err := json.Unmarshal([]byte(result), &vulnerabilities); err != nil {
		// If AI didn't return valid JSON, return a generic response
		vulnerabilities = []Vulnerability{
			{
				Type:        "AI Analysis",
				Severity:    "INFO",
				File:        filename,
				Description: result,
				Suggestion:  "Please review the AI analysis above manually.",
			},
		}
	}

	return &AuditResponse{
		RiskCount:      len(vulnerabilities),
		Vulnerabilities: vulnerabilities,
	}, nil
}

func (s *AIService) callAI(prompt string) (string, error) {
	if s.config.APIKey == "" {
		return "", errors.New("AI API key not configured")
	}

	switch s.config.Provider {
	case "openai":
		return s.callOpenAI(prompt)
	default:
		return s.callOpenAI(prompt)
	}
}

func (s *AIService) callOpenAI(prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model": s.config.Model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a cybersecurity expert and CTF assistant. Provide clear, actionable guidance."},
			{"role": "user", "content": prompt},
		},
		"max_tokens":   s.config.MaxTokens,
		"temperature":  s.config.Temperature,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", errors.New("no AI response")
	}

	return result.Choices[0].Message.Content, nil
}
