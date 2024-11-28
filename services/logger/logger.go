package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RequestLog struct {
	Timestamp time.Time   `json:"timestamp"`
	Request   interface{} `json:"request"`
	Response  interface{} `json:"response"`
}

type Logger struct {
	logDir string
	mu     sync.Mutex
}

func NewLogger(logDir string) (*Logger, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	return &Logger{
		logDir: logDir,
	}, nil
}

// compressJSON 压缩JSON数据（去除空格和换行符）
func compressJSON(data []byte) ([]byte, error) {
	// 创建一个buffer来存储压缩后的JSON
	buf := new(bytes.Buffer)

	// 使用json.Compact压缩JSON
	if err := json.Compact(buf, data); err != nil {
		return nil, fmt.Errorf("压缩JSON失败: %v", err)
	}

	return buf.Bytes(), nil
}

func (l *Logger) Log(request, response interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 创建日志记录
	log := RequestLog{
		Timestamp: time.Now(),
		Request:   request,
		Response:  response,
	}

	// 将日志记录转换为JSON
	logJSON, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 压缩JSON数据
	compressedJSON, err := compressJSON(logJSON)
	if err != nil {
		return fmt.Errorf("压缩数据失败: %v", err)
	}

	// 生成日志文件名
	filename := fmt.Sprintf("%s.log", time.Now().Format("2006010215"))
	filepath := filepath.Join(l.logDir, filename)

	// 打开日志文件（追加模式）
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}
	defer file.Close()

	// 写入压缩后的数据，每条记录一行
	compressedJSON = append(compressedJSON, '\n')
	if _, err := file.Write(compressedJSON); err != nil {
		return fmt.Errorf("写入日志失败: %v", err)
	}

	return nil
}

// ReadLog 读取日志内容
func (l *Logger) ReadLog(filename string) ([]RequestLog, error) {
	filepath := filepath.Join(l.logDir, filename)

	// 读取文件内容
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("读取日志文件失败: %v", err)
	}

	var logs []RequestLog

	// 按行处理
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		// 解析JSON
		var log RequestLog
		if err := json.Unmarshal(line, &log); err != nil {
			return nil, fmt.Errorf("JSON解析失败: %v", err)
		}

		logs = append(logs, log)
	}

	return logs, nil
}
