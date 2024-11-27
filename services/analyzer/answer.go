package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"grading-api/config"
	"grading-api/types"

	"github.com/sashabaranov/go-openai"
)

type AnswerAnalyzer struct {
	client *openai.Client
}

func NewAnswerAnalyzer(client *openai.Client) *AnswerAnalyzer {
	return &AnswerAnalyzer{client: client}
}

func (aa *AnswerAnalyzer) Analyze(ctx context.Context, question, studentAnswer, referenceAnswer string, questionAnalysis types.QuestionAnalysis, temperature float32) (types.AnswerAnalysis, error) {
	// 构建包含题目分析结果的输入
	input := map[string]interface{}{
		"question":          question,
		"student_answer":    studentAnswer,
		"reference_answer":  referenceAnswer,
		"question_analysis": questionAnalysis,
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return types.AnswerAnalysis{}, fmt.Errorf("构建输入JSON失败: %v", err)
	}

	resp, err := aa.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: config.AppConfig.OpenAI.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: config.AppConfig.Roles.AnswerAnalyzer.Prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: string(inputJSON),
				},
			},
			Temperature: temperature,
		},
	)

	if err != nil {
		return types.AnswerAnalysis{}, fmt.Errorf("答案分析请求失败: %v", err)
	}

	var analysis types.AnswerAnalysis
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &analysis); err != nil {
		return types.AnswerAnalysis{}, fmt.Errorf("解析答案分析结果失败: %v", err)
	}

	return analysis, nil
}
