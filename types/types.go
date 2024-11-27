package types

// GradingRequest 批改请求结构
type GradingRequest struct {
	Question        string       `json:"question" binding:"required"`         // 题目
	ReferenceAnswer string       `json:"reference_answer" binding:"required"` // 参考答案
	Analysis        string       `json:"analysis" binding:"required"`         // 解析
	StudentAnswer   string       `json:"student_answer" binding:"required"`   // 学生作答
	OpenAIConfig    OpenAIConfig `json:"openai_config" binding:"required"`    // OpenAI配置
	Stream          bool         `json:"stream,omitempty"`                    // 是否流式响应
}

// OpenAIConfig OpenAI配置结构
type OpenAIConfig struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url,omitempty"`
}

// QuestionAnalysis 题目分析结果
type QuestionAnalysis struct {
	KnowledgePoints []string       `json:"knowledge_points"`
	DifficultyLevel string         `json:"difficulty_level"`
	KeySteps        []string       `json:"key_steps"`
	ScoringCriteria map[string]int `json:"scoring_criteria"`
	CommonMistakes  []string       `json:"common_mistakes"`
	EvaluationFocus string         `json:"evaluation_focus"`
}

// AnswerAnalysis 答案分析结果
type AnswerAnalysis struct {
	ThinkingProcess  string   `json:"thinking_process"`
	Strengths        []string `json:"strengths"`
	Weaknesses       []string `json:"weaknesses"`
	InnovationPoints []string `json:"innovation_points"`
	KnowledgeMastery string   `json:"knowledge_mastery"`
	ErrorAnalysis    string   `json:"error_analysis"`
}

// TeacherResult 教师评分结果
type TeacherResult struct {
	TeacherRole string  `json:"teacher_role"` // 教师角色
	Score       float64 `json:"score"`        // 分数
	Comments    string  `json:"comments"`     // 评语
	Suggestions string  `json:"suggestions"`  // 建议
}

// FinalResult 最终评分结果
type FinalResult struct {
	FinalScore    float64 `json:"final_score"`    // 最终分数
	FinalComments string  `json:"final_comments"` // 最终评语
	Explanation   string  `json:"explanation"`    // 分数说明
}

// GradingResponse 批改结果结构
type GradingResponse struct {
	QuestionAnalysis QuestionAnalysis `json:"question_analysis"`
	AnswerAnalysis   AnswerAnalysis   `json:"answer_analysis"`
	TeacherResults   []TeacherResult  `json:"teacher_results"`
	FinalResult      FinalResult      `json:"final_result"`
}
