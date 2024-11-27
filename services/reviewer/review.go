package reviewer

import (
	"context"
	"encoding/json"
	"fmt"
	"grading-api/config"
	"grading-api/types"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type Reviewer struct {
	client *openai.Client
}

func NewReviewer(client *openai.Client) *Reviewer {
	return &Reviewer{client: client}
}

func (r *Reviewer) Review(ctx context.Context, req types.GradingRequest, questionAnalysis types.QuestionAnalysis, answerAnalysis types.AnswerAnalysis, teacherResults []types.TeacherResult, temperature float32) (types.FinalResult, error) {
	// 构建包含所有前序结果的输入
	input := map[string]interface{}{
		"question": map[string]interface{}{
			"content":          req.Question,
			"reference_answer": req.ReferenceAnswer,
			"analysis":         req.Analysis,
			"student_answer":   req.StudentAnswer,
		},
		"question_analysis": questionAnalysis,
		"answer_analysis":   answerAnalysis,
		"teachers":          teacherResults,
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return types.FinalResult{}, fmt.Errorf("构建输入JSON失败: %v", err)
	}

	resp, err := r.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: config.AppConfig.OpenAI.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: config.AppConfig.Roles.Reviewer.Prompt,
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
		log.Printf("最终审核请求失败: %v", err)
		return types.FinalResult{}, fmt.Errorf("最终审核请求失败: %v", err)
	}

	log.Printf("OpenAI Response: %s", resp.Choices[0].Message.Content)

	// 清理响应内容中的markdown格式
	content := cleanMarkdownJSON(resp.Choices[0].Message.Content)

	var result types.FinalResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		log.Printf("解析失败的响应内容: %s", content)
		return types.FinalResult{}, fmt.Errorf("解析最终审核结果失败: %v, 响应内容: %s", err, content)
	}

	return result, nil
}

// cleanMarkdownJSON 清理markdown格式，提取纯JSON内容
func cleanMarkdownJSON(content string) string {
	// 移除可能的markdown代码块标记
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")

	// 清理首尾的空白字符
	content = strings.TrimSpace(content)

	return content
}
