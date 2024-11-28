# AI 智能批改系统

这是一个基于 OpenAI API 的智能批改系统，通过多角色协作完成对学生答案的全面评估。

## 系统架构

### 角色设计
系统包含四个专业角色，形成完整的评估链：

1. **题目分析专家**
   - 识别核心知识点
   - 评估题目难度
   - 分析关键步骤
   - 预判常见错误
   - 制定评分标准

2. **答案分析专家**
   - 分析解题思路
   - 识别创新点
   - 评估优缺点
   - 分析错误原因
   - 评估知识掌握程度

3. **数学教师（并行评分）**
   - 根据前序分析结果评分
   - 给出详细评语
   - 提供改进建议

4. **教学总监**
   - 综合所有分析结果
   - 权衡教师评分
   - 给出最终评定
   - 提供全面反馈

### 技术栈
- Go 1.22+
- Gin Web Framework
- OpenAI API
- YAML 配置

## 项目结构

## API 文档

### 批改试题 (POST /api/grade)
```json
{
"question": "求解方程：x² + 2x + 1 = 0",
"reference_answer": "这是一个完全平方式，可以写成(x+1)² = 0，所以x = -1",
"analysis": "这道题考察学生对完全平方式的识别和求解能力",
"student_answer": "x² + 2x = -1，用求根公式得到x = -1"
}
```


#### 请求参数说明
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| question | string | 是 | 题目内容，包含具体的题目描述 |
| reference_answer | string | 是 | 参考答案，包含标准解答过程和结果 |
| analysis | string | 是 | 题目解析，说明考察重点和解题思路 |
| student_answer | string | 是 | 学生作答，学生的实际解答内容 |

#### 响应体 (Response Body)
```json
{
    "question_analysis": {
        "knowledge_points": [
            "方程的解法",
            "完全平方公式",
            "因式分解"
        ],
        "difficulty_level": "简单",
        "key_steps": [
            "识别完全平方式",
            "应用因式分解",
            "得出解"
        ],
        "scoring_criteria": {
            "calculation_process": 40,
            "concept_understanding": 40,
            "result_accuracy": 20
        },
        "common_mistakes": [
            "未识别完全平方式",
            "错误地使用求根公式",
            "忽略解的个数"
        ],
        "evaluation_focus": "关注学生对完全平方形式的识别能力和解题步骤的完整性"
    },
    "answer_analysis": {
        "thinking_process": "学生试图通过移项并使用求根公式解决方程，但未能识别这是一个完全平方形式的方程。",
        "strengths": [],
        "weaknesses": [
            "未识别完全平方形式",
            "错误地移项，导致结果不正确",
            "错误地应用求根公式"
        ],
        "innovation_points": [],
        "knowledge_mastery": "对方程的解法掌握不够，特别是对完全平方公式的理解存在缺陷。",
        "error_analysis": "学生忽略了方程可以因式分解的特性，导致错误使用求根公式，产生错误解。"
    },
    "teacher_results": [
        {
            "teacher_role": "教师1",
            "score": 30,
            "comments": "学生未能识别完全平方形式，错误地移项并应用求根公式，导致解答错误。没有正确理解方程的性质和解法。",
            "suggestions": "建议学生复习完全平方公式的识别与应用，练习因式分解和求解简单二次方程，以提高解题能力。"
        },
        {
            "teacher_role": "教师2",
            "score": 30,
            "comments": "学生未能识别方程为完全平方形式，错误地移项并应用求根公式，因此得到的结果不正确。解题过程中缺乏对方程特性的理解。",
            "suggestions": "建议学生复习完全平方公式的概念，并练习识别完全平方形式的方程。同时，多做因式分解的练习，以增强对不同解法的熟悉度。"
        }
    ],
    "final_result": {
        "final_score": 30,
        "final_comments": "学生在解题过程中未能正确识别完全平方形式，导致解法错误。建议加强对完全平方公式和因式分解的理解与应用，提升解题能力。",
        "explanation": "两位教师均给予最低分数，强调了学生在识别方程性质和解法上的重大缺陷。最终分数反映了学生在该题目上的表现，建议着重关注基础概念的掌握。"
    }
}
```

#### 请求体 (Request Body)

#### 响应字段说明

##### 题目分析 (question_analysis)
| 字段 | 类型 | 说明 |
|------|------|------|
| knowledge_points | string[] | 知识点列表，题目涉及的核心知识点 |
| difficulty_level | string | 难度级别，可能值：简单/中等/困难 |
| key_steps | string[] | 关键步骤，解题的主要步骤列表 |
| scoring_criteria | object | 评分标准，各部分的分值分配 |
| common_mistakes | string[] | 常见错误，可能出现的典型错误 |
| evaluation_focus | string | 评估重点，重点关注的评分要素 |

##### 答案分析 (answer_analysis)
| 字段 | 类型 | 说明 |
|------|------|------|
| thinking_process | string | 思维过程，学生的解题思路分析 |
| strengths | string[] | 优点列表，解答中的亮点 |
| weaknesses | string[] | 缺点列表，解答中的不足 |
| innovation_points | string[] | 创新点，解答中的创新思维 |
| knowledge_mastery | string | 知识掌握程度的评估 |
| error_analysis | string | 错误分析，详细的错误原因 |

##### 教师评分 (teacher_results)
| 字段 | 类型 | 说明 |
|------|------|------|
| teacher_role | string | 教师角色，标识具体的评分教师 |
| score | number | 分数，0-100的数值 |
| comments | string | 评语，对答案的具体点评 |
| suggestions | string | 建议，改进的具体建议 |

##### 最终结果 (final_result)
| 字段 | 类型 | 说明 |
|------|------|------|
| final_score | number | 最终分数，0-100的综合评分 |
| final_comments | string | 最终评语，综合所有分析的总体评价 |
| explanation | string | 评分说明，解释最终分数的评定依据 |

### 请求示例
使用 curl 发送请求：
```shell
curl -X POST http://localhost:8080/api/grade \
-H "Content-Type: application/json" \
-d '{
"question": "求解方程：x² + 2x + 1 = 0",
"reference_answer": "这是一个完全平方式，可以写成(x+1)² = 0，所以x = -1",
"analysis": "这道题考察学生对完全平方式的识别和求解能力",
"student_answer": "x² + 2x = -1，用求根公式得到x = -1"}'
```