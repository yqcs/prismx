package font

import (
	"fmt"
)

// 前景色
const (
	textBlack = iota + 30
	textRed
	textGreen
	textYellow
	textBlue
	textPurple
	textCyan
	textWhite
)

// 背景色
const (
	backBlack = iota + 40
	backRed
	backGreen
	backYellow
	backBlue
	backPurple
	backCyan
	backWhite
)

func Black(str string) string {
	return textColor(textBlack, str)
}

func Red(str string) string {
	return textColor(textRed, str)
}
func Yellow(str string) string {
	return textColor(textYellow, str)
}
func Green(str string) string {
	return textColor(textGreen, str)
}
func Cyan(str string) string {
	return textColor(textCyan, str)
}
func Blue(str string) string {
	return textColor(textBlue, str)
}
func Purple(str string) string {
	return textColor(textPurple, str)
}
func White(str string) string {
	return textColor(textWhite, str)
}

func BackBlack(str string) string {
	return backColor(backBlack, str)
}
func BackRed(str string) string {
	return backColor(backRed, str)
}
func BackYellow(str string) string {
	return backColor(backYellow, str)
}
func BackGreen(str string) string {
	return backColor(backGreen, str)
}
func BackCyan(str string) string {
	return backColor(backCyan, str)
}
func BackBlue(str string) string {
	return backColor(backBlue, str)
}
func BackPurple(str string) string {
	return backColor(backPurple, str)
}
func BackWhite(str string) string {
	return backColor(backWhite, str)
}

// backColor 返回字符串，如果是Windows平台则不作任何处理
func backColor(color int, str string) string {
	//if runtime.GOOS != "windows" {
	//	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
	//}
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
}

// textColor 返回字符串，如果是Windows平台则不作任何处理
func textColor(color int, str string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
	//return str
}
