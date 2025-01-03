package logger

import (
	"bytes"
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"io"
	"os"
	"prismx_cli/utils/logger/color"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	Global.EnableColor()
}

type (
	logger struct {
		prefix     string
		level      uint32
		output     io.Writer
		levels     []string
		color      *color.Color
		bufferPool sync.Pool
		mutex      sync.Mutex
	}

	lvl uint8
)

const (
	debug lvl = iota + 1
	info
	warn
	errorType
	OFF
	panicLevel
	fatalLevel
	scanInfo
	scanSuccess
	scanFailure
)

var (
	Global = New("$")
)

func New(prefix string) (l *logger) {
	l = &logger{
		level:  uint32(info),
		prefix: prefix,
		color:  color.New(),
		bufferPool: sync.Pool{
			New: func() any {
				return bytes.NewBuffer(make([]byte, 256))
			},
		},
	}
	l.initLevels()
	l.SetOutput(colorable.NewColorableStdout())
	return
}

func (l *logger) initLevels() {
	l.levels = []string{
		"-",
		l.color.Grey("DEBUG"),
		l.color.Blue("INFO"),
		l.color.Yellow("WARN"),
		l.color.Red("ERROR", color.U),
		"",
		l.color.Magenta("PANIC", color.U),
		l.color.Red("FATAL"),
		l.color.Cyan("SCAN"),
		l.color.GreenBg("VULN"),
		l.color.Red("FAILURE"),
	}
}

func (l *logger) DisableColor() {
	l.color.Disable()
	l.initLevels()
}

func (l *logger) EnableColor() {
	l.color.Enable()
	l.initLevels()
}

func (l *logger) SetPrefix(p string) {
	l.prefix = p
}

func (l *logger) Level() lvl {
	return lvl(atomic.LoadUint32(&l.level))
}

func (l *logger) SetLevel(level lvl) {
	atomic.StoreUint32(&l.level, uint32(level))
}

func (l *logger) Output() io.Writer {
	return l.output
}

func (l *logger) SetOutput(w io.Writer) {
	l.output = w
	if w, ok := w.(*os.File); !ok || !isatty.IsTerminal(w.Fd()) {
		l.EnableColor()
	}
}

func (l *logger) Color() *color.Color {
	return l.color
}

func (l *logger) print(i ...any) {
	l.log(0, "", i...)
}

func (l *logger) printf(format string, args ...any) {
	l.log(0, format, args...)
}

func (l *logger) debug(i ...any) {
	l.log(debug, "", i...)
}

func (l *logger) debugf(format string, args ...any) {
	l.log(debug, format, args...)
}

func (l *logger) info(i ...any) {
	l.log(info, "", i...)
}

func (l *logger) infof(format string, args ...any) {
	l.log(info, format, args...)
}

func (l *logger) warn(i ...any) {
	l.log(warn, "", i...)
}

func (l *logger) warnf(format string, args ...any) {
	l.log(warn, format, args...)
}

func (l *logger) error(i ...any) {
	l.log(errorType, "", i...)
}

func (l *logger) errorf(format string, args ...any) {
	l.log(errorType, format, args...)
}

func (l *logger) fatal(i ...any) {
	l.log(fatalLevel, "", i...)
	os.Exit(-1)
}

func (l *logger) fatalf(format string, args ...any) {
	l.log(fatalLevel, format, args...)
	os.Exit(-1)
}

func (l *logger) panic(i ...any) {
	l.log(panicLevel, "", i...)
	os.Exit(-1)
}

func (l *logger) panicf(format string, args ...any) {
	l.log(panicLevel, format, args...)
	os.Exit(-1)
}

func (l *logger) message(i ...any) {
	l.log(scanInfo, "", i...)
}
func (l *logger) messagef(format string, i ...any) {
	l.log(scanInfo, format, i...)
}
func (l *logger) success(i ...any) {
	l.log(scanSuccess, "", i...)
}
func (l *logger) successf(format string, i ...any) {
	l.log(scanSuccess, format, i...)
}
func (l *logger) failure(i ...any) {
	l.log(scanFailure, "", i...)
}
func (l *logger) failuref(format string, i ...any) {
	l.log(scanFailure, format, i...)
}

func Print(i ...any) {
	Global.print(i...)
}

func Printf(format string, args ...any) {
	Global.printf(format, args...)
}

func Debug(i ...any) {
	Global.debug(i...)
}

func Debugf(format string, args ...any) {
	Global.debugf(format, args...)
}

func Info(i ...any) {
	Global.info(i...)
}

func Infof(format string, args ...any) {
	Global.infof(format, args...)
}

func Warn(i ...any) {
	Global.warn(i...)
}

func Warnf(format string, args ...any) {
	Global.warnf(format, args...)
}

func Error(i ...any) {
	Global.error(i...)
}

func Errorf(format string, args ...any) {
	Global.errorf(format, args...)
}

func Fatal(i ...any) {
	Global.fatal(i...)
}

func Fatalf(format string, args ...any) {
	Global.fatalf(format, args...)
}

func Panic(i ...any) {
	Global.panic(i...)
}

func Panicf(format string, args ...any) {
	Global.panicf(format, args...)
}

//----------------扫描消息池--------------------

// ScanSuccess 发现问题，上报消息池
func ScanSuccess(i ...any) {
	Global.success(i)
}

// ScanFailuref 扫描过程出现错误，上报消息池
func ScanFailuref(format string, args ...any) {
	Global.failuref(format, args...)
}

// ScanMessage 扫描信息，上报消息池
func ScanMessage(i ...any) {
	Global.message(i...)
}

// ScanMessagef 扫描信息，上报消息池
func ScanMessagef(format string, args ...any) {
	Global.messagef(format, args...)
}

func (l *logger) log(level lvl, format string, args ...any) {
	if level >= l.Level() || level == 0 {
		buf := l.bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer l.bufferPool.Put(buf)
		message := ""
		buf = bytes.NewBuffer([]byte(fmt.Sprintf("[%s - %s]%s ", l.levels[level], time.Now().Format("15:04:05"), l.prefix)))
		if format == "" {
			message = fmt.Sprint(args...)
		} else {
			message = fmt.Sprintf(format, args...)
		}
		buf.WriteString(message + "\n")
		l.mutex.Lock()
		l.output.Write(buf.Bytes())
		l.mutex.Unlock()
	}
}
