// plugin/example/plugin/timetool/tool.go
package main

import (
	"fmt"
	"time"

	"github.com/gophertool/tool/plugin"
)

// TimeTool 时间工具插件
type TimeTool struct {
	// 可以添加一些插件配置或状态
	TimeFormat string
}

// NewTimeTool 创建一个新的时间工具插件实例
func NewTimeTool() *TimeTool {
	return &TimeTool{
		TimeFormat: time.RFC3339,
	}
}

// GetPluginInfo 获取插件信息
func (t *TimeTool) GetPluginInfo() (plugin.PluginInfo, error) {
	return plugin.PluginInfo{
		Name:        "time_tool",
		Version:     "1.0.0",
		Description: "提供时间相关的工具函数",
		Author:      "GopherTool",
	}, nil
}

// GetTools 获取插件提供的工具列表
func (t *TimeTool) GetTools() ([]plugin.Tool, error) {
	// 定义当前时间工具
	currentTimeTool := plugin.NewTool(
		"current_time",
		"获取当前时间",
		plugin.WithString("format",
			plugin.Description("时间格式，例如：2006-01-02 15:04:05"),
			plugin.Default("2006-01-02 15:04:05"),
		),
		plugin.WithString("timezone",
			plugin.Description("时区，例如：Asia/Shanghai, UTC"),
			plugin.Default("Local"),
		),
	)

	// 定义时间转换工具
	timeConvertTool := plugin.NewTool(
		"time_convert",
		"转换时间格式",
		plugin.WithString("time",
			plugin.Description("要转换的时间字符串"),
			plugin.Required(),
		),
		plugin.WithString("source_format",
			plugin.Description("源时间格式"),
			plugin.Required(),
		),
		plugin.WithString("target_format",
			plugin.Description("目标时间格式"),
			plugin.Required(),
		),
	)

	// 定义时间计算工具
	timeCalcTool := plugin.NewTool(
		"time_calc",
		"时间计算（加减天数、小时等）",
		plugin.WithString("time",
			plugin.Description("基准时间，如果为空则使用当前时间"),
		),
		plugin.WithString("format",
			plugin.Description("时间格式"),
			plugin.Default("2006-01-02 15:04:05"),
		),
		plugin.WithInteger("years",
			plugin.Description("要加减的年数"),
			plugin.Default(0),
		),
		plugin.WithInteger("months",
			plugin.Description("要加减的月数"),
			plugin.Default(0),
		),
		plugin.WithInteger("days",
			plugin.Description("要加减的天数"),
			plugin.Default(0),
		),
		plugin.WithInteger("hours",
			plugin.Description("要加减的小时数"),
			plugin.Default(0),
		),
		plugin.WithInteger("minutes",
			plugin.Description("要加减的分钟数"),
			plugin.Default(0),
		),
		plugin.WithInteger("seconds",
			plugin.Description("要加减的秒数"),
			plugin.Default(0),
		),
	)

	return []plugin.Tool{*currentTimeTool, *timeConvertTool, *timeCalcTool}, nil
}

// CallTool 调用工具
func (t *TimeTool) CallTool(toolName string, params map[string]interface{}) (*plugin.CallToolResult, error) {
	switch toolName {
	case "current_time":
		return t.getCurrentTime(params)
	case "time_convert":
		return t.convertTime(params)
	case "time_calc":
		return t.calculateTime(params)
	default:
		return plugin.NewErrorResult(fmt.Sprintf("未知的工具: %s", toolName)), nil
	}
}

// getCurrentTime 获取当前时间
func (t *TimeTool) getCurrentTime(params map[string]interface{}) (*plugin.CallToolResult, error) {
	// 获取格式参数，如果没有则使用默认格式
	format, ok := params["format"].(string)
	if !ok || format == "" {
		format = "2006-01-02 15:04:05"
	}

	// 获取时区参数，如果没有则使用本地时区
	timezone, ok := params["timezone"].(string)
	if !ok || timezone == "" {
		timezone = "Local"
	}

	// 加载时区
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return plugin.NewErrorResult(fmt.Sprintf("无效的时区: %s, 错误: %v", timezone, err)), nil
	}

	// 获取当前时间并格式化
	currentTime := time.Now().In(loc).Format(format)

	// 创建结果
	result := plugin.NewCallToolResult()
	result.AddTextContent(currentTime, "current_time")

	// 添加一些元数据
	result.SetMeta("timezone", timezone)
	result.SetMeta("format", format)
	result.SetMeta("timestamp", time.Now().Unix())

	return result, nil
}

// convertTime 转换时间格式
func (t *TimeTool) convertTime(params map[string]interface{}) (*plugin.CallToolResult, error) {
	// 获取必要参数
	timeStr, ok := params["time"].(string)
	if !ok || timeStr == "" {
		return plugin.NewErrorResult("缺少必要参数: time"), nil
	}

	sourceFormat, ok := params["source_format"].(string)
	if !ok || sourceFormat == "" {
		return plugin.NewErrorResult("缺少必要参数: source_format"), nil
	}

	targetFormat, ok := params["target_format"].(string)
	if !ok || targetFormat == "" {
		return plugin.NewErrorResult("缺少必要参数: target_format"), nil
	}

	// 解析时间
	parsedTime, err := time.Parse(sourceFormat, timeStr)
	if err != nil {
		return plugin.NewErrorResult(fmt.Sprintf("时间解析失败: %v", err)), nil
	}

	// 转换格式
	convertedTime := parsedTime.Format(targetFormat)

	// 创建结果
	result := plugin.NewCallToolResult()
	result.AddTextContent(convertedTime, "converted_time")

	// 添加一些元数据
	result.SetMeta("original_time", timeStr)
	result.SetMeta("source_format", sourceFormat)
	result.SetMeta("target_format", targetFormat)

	return result, nil
}

// calculateTime 计算时间
func (t *TimeTool) calculateTime(params map[string]interface{}) (*plugin.CallToolResult, error) {
	// 获取基准时间，如果没有则使用当前时间
	var baseTime time.Time
	timeStr, ok := params["time"].(string)
	if ok && timeStr != "" {
		// 获取时间格式
		format, ok := params["format"].(string)
		if !ok || format == "" {
			format = "2006-01-02 15:04:05"
		}

		// 解析时间
		var err error
		baseTime, err = time.Parse(format, timeStr)
		if err != nil {
			return plugin.NewErrorResult(fmt.Sprintf("时间解析失败: %v", err)), nil
		}
	} else {
		// 使用当前时间
		baseTime = time.Now()
	}

	// 获取各个时间单位的增量
	years := getIntParam(params, "years", 0)
	months := getIntParam(params, "months", 0)
	days := getIntParam(params, "days", 0)
	hours := getIntParam(params, "hours", 0)
	minutes := getIntParam(params, "minutes", 0)
	seconds := getIntParam(params, "seconds", 0)

	// 计算新时间
	newTime := baseTime.AddDate(years, months, days)
	newTime = newTime.Add(time.Duration(hours) * time.Hour)
	newTime = newTime.Add(time.Duration(minutes) * time.Minute)
	newTime = newTime.Add(time.Duration(seconds) * time.Second)

	// 获取输出格式
	format, ok := params["format"].(string)
	if !ok || format == "" {
		format = "2006-01-02 15:04:05"
	}

	// 格式化结果
	resultTime := newTime.Format(format)

	// 创建结果
	result := plugin.NewCallToolResult()
	result.AddTextContent(resultTime, "calculated_time")

	// 添加一些元数据
	result.SetMeta("base_time", baseTime.Format(format))
	result.SetMeta("years", years)
	result.SetMeta("months", months)
	result.SetMeta("days", days)
	result.SetMeta("hours", hours)
	result.SetMeta("minutes", minutes)
	result.SetMeta("seconds", seconds)

	return result, nil
}

// getIntParam 从参数中获取整数值，如果获取失败则返回默认值
func getIntParam(params map[string]interface{}, name string, defaultValue int) int {
	value, ok := params[name]
	if !ok {
		return defaultValue
	}

	// 尝试转换为整数
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		// 尝试将字符串转换为整数
		var intValue int
		_, err := fmt.Sscanf(v, "%d", &intValue)
		if err == nil {
			return intValue
		}
	}

	return defaultValue
}

func main() {
	// 创建一个时间工具插件实例
	timeToolPlugin := NewTimeTool()

	// 启动插件服务
	plugin.ServePlugin(timeToolPlugin)
}
