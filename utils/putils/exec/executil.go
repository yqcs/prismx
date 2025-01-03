package executil

import (
	"encoding/base64"
	"io"
	"os/exec"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/unicode"
)

const (
	// token of prefix for space, escape char '\'
	tkSpPrefix = '\\'
	// the space
	tkSp = ' '
)

// splitCmdArgs splits cmds with spaces.
// It recognizes "\ " as a " " (space) in arguments of the command.
// '\' is escape char only effect with space ' '.
// so "\\a" is also the "\\a" in argment, not "\a".
// e.g. "b\a" -> "b\a".
func splitCmdArgs(cmds string) []string {
	raw := []byte(cmds)
	length := len(raw)

	argments := make([]string, 0, 4)

	escape := false
	// argment recognizing
	argRec := false
	var rawArg []byte

	fillArgChar := func(c byte) []byte {
		if rawArg == nil {
			rawArg = make([]byte, 0, 16)
		}
		rawArg = append(rawArg, c)
		return rawArg
	}

	finishArg := func() {
		// finish one argment recognized
		if len(rawArg) > 0 {
			argment := string(rawArg)
			argments = append(argments, argment)
		}
		rawArg = nil
	}

	for i := 0; i < length; i++ {
		c := raw[i]
		switch c {
		case tkSp:
			if escape {
				fillArgChar(tkSp)
				escape = false // exit excape
			} else if argRec {
				argRec = false
				finishArg()
			}
		case tkSpPrefix:
			if escape {
				// like "\\"
				fillArgChar(tkSpPrefix) // first '\'
			} else {
				escape = true
			}
			if !argRec {
				argRec = true
			}
		default:
			if !argRec {
				argRec = true
			}
			if escape {
				escape = false
				fillArgChar(tkSpPrefix)
			}
			fillArgChar(c)
		}
	}

	if escape {
		fillArgChar(tkSpPrefix)
	}

	if argRec {
		finishArg()
	}

	return argments
}

// convert spaces in arg to "\ "
// so it can split with func [splitCmdAgrs]
// func safeArg(arg string) string {
// 	return strings.ReplaceAll(arg, " ", "\\ ")
// }

// func splitt(s string) (tokens []string) {
// 	for _, ss := range strings.Split(s, " ") {
// 		tokens = append(tokens, strings.Split(ss, "\n")...)
// 	}
// 	return
// }

// Run the specified command in os shell (sh or powershell.exe) and return the output
func Run(cmd string) (string, error) {
	if runtime.GOOS == "windows" {
		return RunPS(cmd)
	}
	return RunSh(splitCmdArgs(cmd)...)
}

/*
RunSafe necessarily does not prevent command injection but prevents/limits damage to an extent
"the os/exec package intentionally does not invoke the system shell and does not expand
any glob patterns or handle other expansions, pipelines, or redirections typically done by shells"
This also mitigates windows security risk in go<1.19. refer https://pkg.go.dev/os/exec
*/
func RunSafe(cmd ...string) (string, error) {
	execpath, err := exec.LookPath(cmd[0])
	if err != nil {
		if runtime.GOOS == "windows" {
			patherror := errors.New("RunSafe does not allow relative exection of binaries (ex ./main) due to security reasons")
			return "", errors.Wrap(err, patherror.Error())
		}
		return "", err
	}

	var cmdArgs []string

	if runtime.GOOS == "windows" {
		/* When command is run in Windows using exec.Command, args are passed as quoted strings
		and are parsed/converted to args before command is run by Go Internally
		*/
		cmdArgs = cmd[1:]
	} else {
		cmdArgs = splitCmdArgs(strings.Join(cmd[1:], " "))
	}

	cmdExec := exec.Command(execpath, cmdArgs...)

	in, _ := cmdExec.StdinPipe()
	errorOut, _ := cmdExec.StderrPipe()
	out, _ := cmdExec.StdoutPipe()
	defer in.Close()
	defer errorOut.Close()
	defer out.Close()

	if err := cmdExec.Start(); err != nil {
		return "", errors.Wrapf(err, "failed to start command:\n%v", strings.Join(cmd, " "))
	}

	outData, _ := io.ReadAll(out)
	errorData, _ := io.ReadAll(errorOut)

	var adbError error = nil

	if err := cmdExec.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			adbError = errors.New("return error")
			outData = errorData
		} else {
			return "", errors.Wrap(err, "process i/o error")
		}
	}

	return string(outData), adbError
}

// RunSh the specified command through sh
func RunSh(cmd ...string) (string, error) {
	cmdExec := exec.Command("sh", "-c", strings.Join(cmd, " "))
	in, _ := cmdExec.StdinPipe()
	errorOut, _ := cmdExec.StderrPipe()
	out, _ := cmdExec.StdoutPipe()
	defer in.Close()
	defer errorOut.Close()
	defer out.Close()

	if err := cmdExec.Start(); err != nil {
		errorData, _ := io.ReadAll(errorOut)
		return "", errors.Wrapf(err, "failed to start process %v", string(errorData))
	}

	outData, _ := io.ReadAll(out)
	errorData, _ := io.ReadAll(errorOut)

	var adbError error = nil

	if err := cmdExec.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			adbError = errors.New("sh return error")
			outData = errorData
		} else {
			return "", errors.New("start sh process error")
		}
	}

	return string(outData), adbError
}

// RunPS runs the specified command through powershell.exe
func RunPS(cmd string) (string, error) {

	/*
		Run command in powershell using -EncodedCommand flag and base64 of actual command
		this makes it possible to run complex quotation marks or curly braces without manual escaping
		more details can be found at:
		https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_powershell_exe?view=powershell-5.1#-encodedcommand-base64encodedcommand
	*/

	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	encodedCmd, err := utf16.NewEncoder().String(cmd)
	if err != nil {
		return "", err
	}
	b64cmd := base64.StdEncoding.EncodeToString([]byte(encodedCmd))

	cmdExec := exec.Command("powershell.exe", "-EncodedCommand", b64cmd)

	in, _ := cmdExec.StdinPipe()
	errorOut, _ := cmdExec.StderrPipe()
	out, _ := cmdExec.StdoutPipe()
	defer in.Close()
	defer errorOut.Close()
	defer out.Close()

	if err := cmdExec.Start(); err != nil {
		return "", errors.New("start powershell.exe process error")
	}

	outData, _ := io.ReadAll(out)
	errorData, _ := io.ReadAll(errorOut)

	var adbError error = nil

	if err := cmdExec.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			adbError = errors.New("powershell.exe return error")
			outData = errorData
		} else {
			return "", errors.New("start powershell.exe process error")
		}
	}

	return string(outData), adbError
}
