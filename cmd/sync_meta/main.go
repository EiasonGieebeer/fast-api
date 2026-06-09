package main

import (
	"fmt"
	"os"

	"github.com/QuantumNous/fast-api/common"
	"github.com/QuantumNous/fast-api/model"
	"github.com/QuantumNous/fast-api/service"
	"github.com/QuantumNous/fast-api/setting/ratio_setting"
	"github.com/joho/godotenv"
)

type modelMeta struct {
	ContextLength    string
	MaxOutputTokens  string
	ParameterCount   string
	KnowledgeCutoff  string
	ReleaseDate      string
	InputModalities  string
	OutputModalities string
	Capabilities     string
}

var specs = map[string]modelMeta{
	// === Anthropic (parameter counts undisclosed) ===
	"claude-haiku-4-5": {
		ContextLength: "200K", MaxOutputTokens: "64K", ParameterCount: "",
		KnowledgeCutoff: "2025-02", ReleaseDate: "2025-10-01",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","tools","json_mode","structured_output","system_prompt","caching"]`,
	},
	"claude-opus-4-6": {
		ContextLength: "1M", MaxOutputTokens: "128K", ParameterCount: "",
		KnowledgeCutoff: "2025-08", ReleaseDate: "2026-02-05",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},
	"claude-opus-4-7": {
		ContextLength: "1M", MaxOutputTokens: "128K", ParameterCount: "",
		KnowledgeCutoff: "2026-01", ReleaseDate: "2026-04-16",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter","web_search"]`,
	},
	"claude-opus-4-8": {
		ContextLength: "1M", MaxOutputTokens: "128K", ParameterCount: "",
		KnowledgeCutoff: "2026-01", ReleaseDate: "2026-05-28",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter","web_search"]`,
	},
	"claude-sonnet-4-6": {
		ContextLength: "1M", MaxOutputTokens: "64K", ParameterCount: "",
		KnowledgeCutoff: "2025-08", ReleaseDate: "2026-02-17",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},

	// === OpenAI ===
	"gpt-4o": {
		ContextLength: "128K", MaxOutputTokens: "16K", ParameterCount: "1.8T",
		KnowledgeCutoff: "2023-10", ReleaseDate: "2024-05-13",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","tools","json_mode","structured_output","system_prompt","caching"]`,
	},
	"gpt-4.1": {
		ContextLength: "1M", MaxOutputTokens: "32K", ParameterCount: "1.8T",
		KnowledgeCutoff: "2024-06", ReleaseDate: "2025-04-14",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},
	"gpt-4o-mini-tts": {
		ContextLength: "TTS 模型", MaxOutputTokens: "音频", ParameterCount: "8B",
		KnowledgeCutoff: "", ReleaseDate: "2024",
		InputModalities: `["text"]`, OutputModalities: `["audio"]`,
		Capabilities: `["streaming","system_prompt"]`,
	},
	"gpt-image-2": {
		ContextLength: "图像生成", MaxOutputTokens: "图像", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "",
		InputModalities: `["text"]`, OutputModalities: `["image"]`,
		Capabilities: `["streaming"]`,
	},
	"gpt-5.2": {
		ContextLength: "400K", MaxOutputTokens: "128K", ParameterCount: "",
		KnowledgeCutoff: "2025-08", ReleaseDate: "2025-12-11",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter","web_search"]`,
	},
	"gpt-5.2-codex": {
		ContextLength: "400K", MaxOutputTokens: "128K", ParameterCount: "",
		KnowledgeCutoff: "2025-08", ReleaseDate: "2026-01",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},
	"gpt-5.1-codex": {
		ContextLength: "400K", MaxOutputTokens: "128K", ParameterCount: "",
		KnowledgeCutoff: "2024-09", ReleaseDate: "2025-11-19",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},
	"gpt-5.3-codex": {
		ContextLength: "400K", MaxOutputTokens: "128K", ParameterCount: "",
		KnowledgeCutoff: "2025-08", ReleaseDate: "2026-02-05",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},
	"gpt-5.4": {
		ContextLength: "1.05M", MaxOutputTokens: "128K", ParameterCount: "",
		KnowledgeCutoff: "2025-08", ReleaseDate: "2026-03-05",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter","web_search"]`,
	},
	"gpt-5.4-mini": {
		ContextLength: "400K", MaxOutputTokens: "128K", ParameterCount: "",
		KnowledgeCutoff: "2025-08", ReleaseDate: "2026-03-17",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","tools","json_mode","structured_output","system_prompt","caching"]`,
	},
	"gpt-5.4-nano": {
		ContextLength: "400K", MaxOutputTokens: "128K", ParameterCount: "",
		KnowledgeCutoff: "2025-08", ReleaseDate: "2026-03-17",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","tools","json_mode","structured_output","system_prompt","caching"]`,
	},
	"gpt-5.5": {
		ContextLength: "1.05M", MaxOutputTokens: "128K", ParameterCount: "",
		KnowledgeCutoff: "2025-12", ReleaseDate: "2026-04-23",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter","web_search"]`,
	},
	"gpt-chat-latest": {
		ContextLength: "自定义别名", MaxOutputTokens: "", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter","web_search"]`,
	},

	// === Google / Gemini ===
	"gemini-2.5-flash": {
		ContextLength: "1M", MaxOutputTokens: "65K", ParameterCount: "",
		KnowledgeCutoff: "2025-01", ReleaseDate: "2025-05-20",
		InputModalities: `["text","image","audio","video"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","tools","json_mode","structured_output","system_prompt","caching"]`,
	},
	"gemini-2.5-flash-lite": {
		ContextLength: "1M", MaxOutputTokens: "65K", ParameterCount: "",
		KnowledgeCutoff: "2025-01", ReleaseDate: "2025-06-17",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","tools","json_mode","structured_output","system_prompt","caching"]`,
	},
	"gemini-2.5-flash-image": {
		ContextLength: "图像生成", MaxOutputTokens: "图像", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2025",
		InputModalities: `["text"]`, OutputModalities: `["image"]`,
		Capabilities: `["streaming","system_prompt"]`,
	},
	"gemini-2.5-pro": {
		ContextLength: "1M", MaxOutputTokens: "65K", ParameterCount: "",
		KnowledgeCutoff: "2025-01", ReleaseDate: "2025-03-24",
		InputModalities: `["text","image","audio","video"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},
	"gemini-2.5-flash-preview-tts": {
		ContextLength: "TTS 模型", MaxOutputTokens: "音频", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2025",
		InputModalities: `["text"]`, OutputModalities: `["audio"]`,
		Capabilities: `["streaming","system_prompt"]`,
	},
	"gemini-2.5-pro-preview-tts": {
		ContextLength: "TTS 模型", MaxOutputTokens: "音频", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2025",
		InputModalities: `["text"]`, OutputModalities: `["audio"]`,
		Capabilities: `["streaming","system_prompt"]`,
	},
	"gemini-3-flash-preview": {
		ContextLength: "1M", MaxOutputTokens: "65K", ParameterCount: "",
		KnowledgeCutoff: "2025-01", ReleaseDate: "2025-12-17",
		InputModalities: `["text","image","audio","video"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching"]`,
	},
	"gemini-3.1-flash-lite-preview": {
		ContextLength: "1M [已下线]", MaxOutputTokens: "65K", ParameterCount: "",
		KnowledgeCutoff: "2025-01", ReleaseDate: "2026-02",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","tools","json_mode","structured_output","system_prompt","caching"]`,
	},
	"gemini-3.1-flash-image-preview": {
		ContextLength: "图像生成", MaxOutputTokens: "图像", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2026-04",
		InputModalities: `["text"]`, OutputModalities: `["image"]`,
		Capabilities: `["streaming","system_prompt"]`,
	},
	"gemini-3.1-flash-tts-preview": {
		ContextLength: "TTS 模型", MaxOutputTokens: "音频", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2026-05",
		InputModalities: `["text"]`, OutputModalities: `["audio"]`,
		Capabilities: `["streaming","system_prompt"]`,
	},
	"gemini-3-pro-image-preview": {
		ContextLength: "图像生成", MaxOutputTokens: "图像", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2026-03",
		InputModalities: `["text","image"]`, OutputModalities: `["image"]`,
		Capabilities: `["streaming","vision","system_prompt"]`,
	},
	"gemini-3.5-flash": {
		ContextLength: "1M", MaxOutputTokens: "65K", ParameterCount: "",
		KnowledgeCutoff: "2025-01", ReleaseDate: "2026-05-19",
		InputModalities: `["text","image","audio","video"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching","web_search"]`,
	},
	"gemini-embedding-001": {
		ContextLength: "8K", MaxOutputTokens: "向量嵌入", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2024-11",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["embeddings"]`,
	},
	"text-embedding-004": {
		ContextLength: "2K", MaxOutputTokens: "向量嵌入", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2024-11",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["embeddings"]`,
	},

	// === DeepSeek ===
	"deepseek-v4-flash": {
		ContextLength: "1M", MaxOutputTokens: "384K", ParameterCount: "284B (13B active)",
		KnowledgeCutoff: "2026-04", ReleaseDate: "2026-04-24",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","reasoning","tools","json_mode","structured_output","system_prompt","caching"]`,
	},
	"deepseek-v4-pro": {
		ContextLength: "1M", MaxOutputTokens: "384K", ParameterCount: "1.6T (49B active)",
		KnowledgeCutoff: "2026-04", ReleaseDate: "2026-04-24",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},

	// === xAI / Grok ===
	"grok-4.3": {
		ContextLength: "1M", MaxOutputTokens: "无上限", ParameterCount: "",
		KnowledgeCutoff: "2025-12", ReleaseDate: "2026-04-30",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","web_search"]`,
	},

	// === Zhipu / GLM ===
	"glm-5": {
		ContextLength: "200K", MaxOutputTokens: "128K", ParameterCount: "744B (40B active)",
		KnowledgeCutoff: "", ReleaseDate: "2026-02-11",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},
	"glm-5.1": {
		ContextLength: "200K", MaxOutputTokens: "131K", ParameterCount: "744B (40B active)",
		KnowledgeCutoff: "", ReleaseDate: "2026-04-07",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},
	"glm-5-turbo": {
		ContextLength: "200K", MaxOutputTokens: "128K", ParameterCount: "744B (40B active)",
		KnowledgeCutoff: "", ReleaseDate: "2026-03-16",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","reasoning","tools","json_mode","structured_output","system_prompt","caching"]`,
	},
	"glm-5v-turbo": {
		ContextLength: "200K", MaxOutputTokens: "128K", ParameterCount: "744B (40B active)",
		KnowledgeCutoff: "", ReleaseDate: "2026-04-01",
		InputModalities: `["text","image","video","file"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching"]`,
	},

	// === Alibaba / Qwen ===
	"qwen3.6-plus": {
		ContextLength: "1M", MaxOutputTokens: "65K", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2026-03-31",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","reasoning","tools","json_mode","structured_output","system_prompt","caching"]`,
	},
	"qwen3.7-plus": {
		ContextLength: "1M", MaxOutputTokens: "65K", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2026-06-01",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","reasoning","tools","json_mode","structured_output","system_prompt","caching"]`,
	},
	"qwen3.7-max": {
		ContextLength: "1M", MaxOutputTokens: "65K", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2026-05-19",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},
	"wan2.7-image-pro": {
		ContextLength: "图像生成", MaxOutputTokens: "图像", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2026-04-01",
		InputModalities: `["text"]`, OutputModalities: `["image"]`,
		Capabilities: `["streaming"]`,
	},

	// === Moonshot / Kimi ===
	"kimi-k2.5": {
		ContextLength: "256K", MaxOutputTokens: "256K", ParameterCount: "1T (32B active)",
		KnowledgeCutoff: "", ReleaseDate: "2026-01-27",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},
	"kimi-k2.6": {
		ContextLength: "256K", MaxOutputTokens: "256K", ParameterCount: "1T (32B active)",
		KnowledgeCutoff: "", ReleaseDate: "2026-04-20",
		InputModalities: `["text","image"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","vision","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter","web_search"]`,
	},

	// === MiniMax ===
	"MiniMax-M2.1": {
		ContextLength: "200K", MaxOutputTokens: "196K", ParameterCount: "230B (10B active)",
		KnowledgeCutoff: "", ReleaseDate: "2025-12-23",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},
	"MiniMax-M2.5": {
		ContextLength: "200K", MaxOutputTokens: "196K", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2026-02-12",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},
	"MiniMax-M2.7": {
		ContextLength: "200K", MaxOutputTokens: "131K", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2026-03-18",
		InputModalities: `["text"]`, OutputModalities: `["text"]`,
		Capabilities: `["function_calling","streaming","reasoning","tools","json_mode","structured_output","system_prompt","caching","code_interpreter"]`,
	},

	// === ByteDance / Doubao ===
	"seedream-4-0": {
		ContextLength: "图像生成", MaxOutputTokens: "图像", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2025-09-09",
		InputModalities: `["text"]`, OutputModalities: `["image"]`,
		Capabilities: `["streaming"]`,
	},
	"seedream-4-5": {
		ContextLength: "图像生成", MaxOutputTokens: "图像", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2025-12-04",
		InputModalities: `["text"]`, OutputModalities: `["image"]`,
		Capabilities: `["streaming"]`,
	},
	"seedream-5-0-lite": {
		ContextLength: "图像生成", MaxOutputTokens: "图像", ParameterCount: "",
		KnowledgeCutoff: "", ReleaseDate: "2026-02-13",
		InputModalities: `["text"]`, OutputModalities: `["image"]`,
		Capabilities: `["streaming","web_search"]`,
	},
}

func main() {
	// Replicate initialization chain from main.go
	// .env is in project root; try a few common paths
	_ = godotenv.Load("../.env") // when run from cmd/sync_meta/
	_ = godotenv.Load(".env")    // when run from project root
	common.InitEnv()
	ratio_setting.InitRatioSettings()
	service.InitHttpClient()
	service.InitTokenEncoders()

	err := model.InitDB()
	if err != nil {
		fmt.Println("FAIL: failed to initialize database:", err)
		os.Exit(1)
	}

	count := 0
	skipped := 0
	for name, meta := range specs {
		var m model.Model
		if err := model.DB.Where("model_name = ?", name).First(&m).Error; err != nil {
			fmt.Printf("SKIP %s: not found in DB\n", name)
			skipped++
			continue
		}
		m.ContextLength = meta.ContextLength
		m.MaxOutputTokens = meta.MaxOutputTokens
		m.ParameterCount = meta.ParameterCount
		m.KnowledgeCutoff = meta.KnowledgeCutoff
		m.ReleaseDate = meta.ReleaseDate
		m.InputModalities = meta.InputModalities
		m.OutputModalities = meta.OutputModalities
		m.Capabilities = meta.Capabilities
		m.ShowSignals = 1
		if err := m.Update(); err != nil {
			fmt.Printf("FAIL %s: %v\n", name, err)
			continue
		}
		fmt.Printf("  OK  %s\n", name)
		count++
	}
	fmt.Printf("\nUpdated: %d  Skipped: %d  Total specs: %d\n", count, skipped, len(specs))
}
