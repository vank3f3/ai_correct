package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"grading-api/services/analyzer"
	"grading-api/services/reviewer"
	"grading-api/services/teacher"
	"grading-api/types"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

var gradingHistory []types.GradingResponse

// 处理批改请求
func handleGrading(_ *openai.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.GradingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client := createOpenAIClient(req.OpenAIConfig)
		ctx := c.Request.Context()
		temperature := randomTemperature()

		// 1. 题目分析
		questionAnalyzer := analyzer.NewQuestionAnalyzer(client)
		questionAnalysis, err := questionAnalyzer.Analyze(ctx, req.Question, req.ReferenceAnswer, req.Analysis, temperature)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("题目分析失败: %v", err)})
			return
		}

		// 2. 答案分析
		answerAnalyzer := analyzer.NewAnswerAnalyzer(client)
		answerAnalysis, err := answerAnalyzer.Analyze(ctx, req.Question, req.StudentAnswer, req.ReferenceAnswer, questionAnalysis, temperature)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("答案分析失败: %v", err)})
			return
		}

		// 3. 教师评分（并发）
		teacherGrader := teacher.NewTeacherGrader(client)
		teacherResultsChan := make(chan types.TeacherResult, 2)
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			result, err := teacherGrader.Grade(ctx, req, questionAnalysis, answerAnalysis, "教师1", temperature)
			if err == nil {
				teacherResultsChan <- result
			}
		}()

		go func() {
			defer wg.Done()
			result, err := teacherGrader.Grade(ctx, req, questionAnalysis, answerAnalysis, "教师2", temperature)
			if err == nil {
				teacherResultsChan <- result
			}
		}()

		go func() {
			wg.Wait()
			close(teacherResultsChan)
		}()

		var teacherResults []types.TeacherResult
		for result := range teacherResultsChan {
			teacherResults = append(teacherResults, result)
		}

		// 4. 最终审核
		reviewerClient := reviewer.NewReviewer(client)
		finalResult, err := reviewerClient.Review(ctx, req, questionAnalysis, answerAnalysis, teacherResults, temperature)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("最终审核失败: %v", err)})
			return
		}

		// 构建完整响应
		response := types.GradingResponse{
			QuestionAnalysis: questionAnalysis,
			AnswerAnalysis:   answerAnalysis,
			TeacherResults:   teacherResults,
			FinalResult:      finalResult,
		}

		gradingHistory = append(gradingHistory, response)
		c.JSON(http.StatusOK, response)
	}
}

// 生成随机Temperature，范围在0.6-1.1之间
func randomTemperature() float32 {
	// 确保随机数生成器已初始化
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// 生成0.6-1.1之间的随机数
	return float32(0.6 + rand.Float64()*0.5)
}

// 获取历史记录
func getGradingHistory(c *gin.Context) {
	c.JSON(http.StatusOK, gradingHistory)
}

// 清除历史记录
func clearGradingHistory(c *gin.Context) {
	gradingHistory = []types.GradingResponse{}
	c.JSON(http.StatusOK, gin.H{"message": "Grading history cleared"})
}
