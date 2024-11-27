package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	OpenAI OpenAIConfig `yaml:"openai"`
	Roles  RolesConfig  `yaml:"roles"`
}

type OpenAIConfig struct {
	DefaultAPIKey  string `yaml:"default_api_key"`
	DefaultBaseURL string `yaml:"default_base_url"`
	Model          string `yaml:"model"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
}

type RoleConfig struct {
	Name   string `yaml:"name"`
	Prompt string `yaml:"prompt"`
}

type RolesConfig struct {
	QuestionAnalyzer RoleConfig `yaml:"question_analyzer"`
	AnswerAnalyzer   RoleConfig `yaml:"answer_analyzer"`
	Teacher          RoleConfig `yaml:"teacher"`
	Reviewer         RoleConfig `yaml:"reviewer"`
}

var AppConfig Config

// LoadConfig 从文件加载配置
func LoadConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		return err
	}

	return nil
}

// GetOpenAIConfig 获取OpenAI配置
func GetOpenAIConfig(customAPIKey, customBaseURL string) OpenAIConfig {
	config := AppConfig.OpenAI

	// 如果提供了自定义配置，则使用自定义配置
	if customAPIKey != "" {
		config.DefaultAPIKey = customAPIKey
	}
	if customBaseURL != "" {
		config.DefaultBaseURL = customBaseURL
	}

	return config
}
