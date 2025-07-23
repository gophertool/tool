package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Level int

const (
	DEBUG Level = 0x0000a1
	INFO  Level = 0x0000b2
	WARN  Level = 0x0000c3
	ERROR Level = 0x0000d4
	DATA  Level = 0x0000e5
	NONE  Level = 0x0000f6
)

type Logger struct {
	log      *log.Logger
	modifier func(string) string
	filter   func(string) bool
}

// callerLevel 全局调用者层级设置，默认为3
var callerLevel = 3

func (l *Logger) Printf(format string, s ...any) {
	expr := fmt.Sprintf(format, s...)
	l.Println(expr)
}

func (l *Logger) Println(s ...any) {
	expr := fmt.Sprint(s...)
	if l.modifier != nil {
		expr = l.modifier(expr)
	}
	if l.filter != nil {
		if l.filter(expr) {
			return
		}
	}
	_, _, depth := findCallerWithLevel(callerLevel)
	_ = l.log.Output(depth, expr)
}

var info = &Logger{
	log.New(os.Stdout, "\r[I]", log.Ldate|log.Ltime|log.Lshortfile),
	Green,
	nil,
}

var warn = &Logger{
	log.New(os.Stdout, "\r[W]", log.Ldate|log.Ltime|log.Llongfile),
	Yellow,
	nil,
}

var err = &Logger{
	log.New(os.Stderr, "\r[E]", log.Ldate|log.Ltime|log.Llongfile),
	Red,
	nil,
}

var dbg = &Logger{
	log.New(os.Stdout, "\r[D]", log.Ldate|log.Ltime|log.Llongfile),
	debugModifier,
	debugFilter,
}

// findCaller 寻找真正的调用者位置
// 跳过本包内的函数调用，找到第一个非log包的调用位置
func findCaller() (file string, line int, depth int) {
	return findCallerWithLevel(3)
}

// findCallerWithLevel 寻找真正的调用者位置，允许用户指定起始层级
// startLevel: 开始搜索的调用栈层级，通常为3（跳过runtime.Caller、findCallerWithLevel、调用者）
// 返回值: file-文件路径, line-行号, depth-实际找到的调用栈深度
func findCallerWithLevel(startLevel int) (file string, line int, depth int) {
	const pkgPath = "github.com/gophertool/tool/log"
	const maxDepth = 15

	// 参数校验
	if startLevel < 1 {
		startLevel = 3
	}

	// 从指定层级开始查找
	for depth = startLevel; depth < maxDepth; depth++ {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			break
		}

		// 如果调用来自log包之外，则认为找到了真正的调用位置
		if !strings.Contains(file, pkgPath) {
			return file, line, depth
		}
	}

	// 如果找不到非log包的调用，则使用默认值
	// 如果startLevel过大，则使用较小的fallback值
	fallbackLevel := startLevel
	if fallbackLevel > 10 {
		fallbackLevel = 3
	}
	_, file, line, ok := runtime.Caller(fallbackLevel)
	if !ok {
		// 如果fallback也失败，使用最小值
		_, file, line, _ = runtime.Caller(1)
		return file, line, 1
	}
	return file, line, fallbackLevel
}

func debugModifier(s string) string {
	file, line, _ := findCallerWithLevel(callerLevel)
	file = file[strings.LastIndex(file, "/")+1:]
	logStr := fmt.Sprintf("%s%s(%d) %s", "> ", file, line, s)
	logStr = Yellow(logStr)
	return logStr
}

// SetCallerLevel 设置日志调用者查找的起始层级
// level: 起始层级，通常为3，如果需要跳过更多层级可以增加此值
// 此函数影响所有日志输出中显示的调用者位置
func SetCallerLevel(level int) {
	if level < 1 {
		level = 3
	}
	callerLevel = level
}

// GetCallerLevel 获取当前的调用者层级设置
// 返回当前设置的调用栈起始层级
func GetCallerLevel() int {
	return callerLevel
}

func debugFilter(_ string) bool {
	//Debug 过滤器
	//if strings.Contains(s, "STEP1:CONNECT") {
	//	return true
	//}
	return false
}

var data = &Logger{
	log.New(os.Stdout, "\r", 0),
	nil,
	nil,
}

func Printf(level Level, format string, s ...any) {
	Println(level, fmt.Sprintf(format, s...))
}

func Println(level Level, s ...any) {
	logStr := fmt.Sprint(s...)
	switch level {
	case DEBUG:
		dbg.Println(logStr)
	case INFO:
		info.Println(logStr)
	case WARN:
		warn.Println(logStr)
	case ERROR:
		err.Println(logStr)
	case DATA:
		data.Println(logStr)
	default:
		return
	}
}

func Debug(s ...any) {
	dbg.Println(fmt.Sprint(s...))
}

func Info(s ...any) {
	info.Println(fmt.Sprint(s...))
}

func Warn(s ...any) {
	warn.Println(fmt.Sprint(s...))
}
func Error(s ...any) {
	err.Println(fmt.Sprint(s...))
}

func Data(s ...any) {
	logStr := fmt.Sprint(s...)
	data.Println(logStr)
}
func Debugf(format string, s ...any) {
	Debug(fmt.Sprintf(format, s...))
}

func Infof(format string, s ...any) {
	Info(fmt.Sprintf(format, s...))
}
func Warnf(format string, s ...any) {
	Warn(fmt.Sprintf(format, s...))
}
func Errorf(format string, s ...any) {
	Error(fmt.Sprintf(format, s...))
}
func Dataf(format string, s ...any) {
	Data(fmt.Sprintf(format, s...))
}

var empty = &Logger{log.New(io.Discard, "", 0), nil, nil}

func SetLevel(level Level) {
	if level > ERROR {
		err = empty
	}
	if level > WARN {
		warn = empty
	}
	if level > INFO {
		info = empty
	}
	if level > DEBUG {
		dbg = empty
	}
	if level > NONE {
		//nothing
	}
}

func SetOutput(writer io.Writer) {
	data.modifier = func(s string) string {
		_, _ = writer.Write([]byte(Clear(s)))
		_, _ = writer.Write([]byte("\r\n"))
		return s
	}
}

func SetOutputFile(level Level, fileName string) {
	// 获取文件路径中的目录部分
	dirName := filepath.Dir(fileName)

	// 创建目录（如果目录不存在）
	if e := os.MkdirAll(dirName, 0755); e != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create directory: %s\n", e)

		return
	}

	// 尝试打开或创建文件
	file, e := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if e != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to open/create log file: %s\n", e)
		return
	}

	// 定义一个函数来设置日志输出到控制台和文件
	setLogOutput := func(logger *Logger, output io.Writer) {
		logger.log.SetOutput(io.MultiWriter(file, os.Stdout)) // 同时输出到文件和控制台
	}

	// 根据日志级别设置日志输出
	switch level {
	case DEBUG:
		setLogOutput(dbg, file)
	case INFO:
		setLogOutput(info, file)
	case WARN:
		setLogOutput(warn, file)
	case ERROR:
		setLogOutput(err, file)
	case NONE:
		return
	default:
		return
	}
}

func LogString(level Level, s string) string {
	var buffer bytes.Buffer
	// 将 logger 的输出改为 buffer
	l := log.New(&buffer, "", log.Ldate|log.Ltime|log.Llongfile)
	switch level {
	case DEBUG:
		l.SetPrefix("#DEBUG")
	case INFO:
		l.SetPrefix("#INFO")
	case WARN:
		l.SetPrefix("#WARN")
	case ERROR:
		l.SetPrefix("#ERROR")
	case DATA:
		l.SetPrefix("#DATA")
	}
	l.Println(s)
	return Clear(buffer.String())
}
