package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"grading-api/config"
	"grading-api/types"

	"github.com/sashabaranov/go-openai"
)

type QuestionAnalyzer struct {
	client *openai.Client
}

func NewQuestionAnalyzer(client *openai.Client) *QuestionAnalyzer {
	return &QuestionAnalyzer{client: client}
}

func (qa *QuestionAnalyzer) Analyze(ctx context.Context, question, referenceAnswer, explanation string, temperature float32) (types.QuestionAnalysis, error) {
	inputJSON := fmt.Sprintf(`{"question":"%s","reference_answer":"%s","explanation":"%s"}`,
		question, referenceAnswer, explanation)

	resp, err := qa.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: config.AppConfig.OpenAI.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: config.AppConfig.Roles.QuestionAnalyzer.Prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: inputJSON,
				},
			},
			Temperature: temperature,
		},
	)

	if err != nil {
		return types.QuestionAnalysis{}, fmt.Errorf("题目分析请求失败: %v", err)
	}

	var analysis types.QuestionAnalysis
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &analysis); err != nil {
		return types.QuestionAnalysis{}, fmt.Errorf("解析题目分析结果失败: %v", err)
	}

	return analysis, nil
}
