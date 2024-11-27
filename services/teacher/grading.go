package teacher

import (
	"context"
	"encoding/json"
	"fmt"
	"grading-api/config"
	"grading-api/types"

	"github.com/sashabaranov/go-openai"
)

type TeacherGrader struct {
	client *openai.Client
}

func NewTeacherGrader(client *openai.Client) *TeacherGrader {
	return &TeacherGrader{client: client}
}

func (tg *TeacherGrader) Grade(ctx context.Context, req types.GradingRequest, questionAnalysis types.QuestionAnalysis, answerAnalysis types.AnswerAnalysis, teacherRole string, temperature float32) (types.TeacherResult, error) {
	// 构建包含前序分析结果的输入
	input := map[string]interface{}{
		"question":          req.Question,
		"reference_answer":  req.ReferenceAnswer,
		"explanation":       req.Analysis,
		"student_answer":    req.StudentAnswer,
		"question_analysis": questionAnalysis,
		"answer_analysis":   answerAnalysis,
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return types.TeacherResult{}, fmt.Errorf("构建输入JSON失败: %v", err)
	}

	resp, err := tg.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: config.AppConfig.OpenAI.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: config.AppConfig.Roles.Teacher.Prompt,
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
		return types.TeacherResult{}, fmt.Errorf("教师评分请求失败: %v", err)
	}

	var result types.TeacherResult
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return types.TeacherResult{}, fmt.Errorf("解析教师评分结果失败: %v", err)
	}

	result.TeacherRole = teacherRole
	return result, nil
}
