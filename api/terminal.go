package api

import (
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"strings"
)

var autoGroup = []string{
	"hello, world",
	"hello, china",
	"golang",
	"goto",
}

type Completion []string

type CmdCall func(term *Term, args []string) error

type Term struct {
	completions map[string]*Completion
	terminal        *terminal.Terminal
	fd          int
	state       *terminal.State
	cmds map[string]CmdCall
}

func NewTerm(prompt string) (*Term, error) {
	term := new(Term)
	term.fd = int(os.Stdin.Fd())
	var err error
	term.state, err = terminal.MakeRaw(term.fd)
	if err != nil {
		return nil, err
	}
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	term.terminal = terminal.NewTerminal(screen, prompt)
    term.cmds = make(map[string]CmdCall)
	return term, nil
}

func (term *Term) Close() error {
	return terminal.Restore(term.fd, term.state)
}

func (term *Term) WriteString(data string) {
    term.terminal.Write([]byte(data))
}

func (term *Term) Start() error {
	for {
		line, err := term.terminal.ReadLine()
		if err != nil {
			return err
		}
		if line == "" {
			continue
		}
		if line == "exit" {
			break
		}
        args := strings.Fields(line)
        if call, ok := term.cmds[args[0]]; ok {
            call(term, args[1:])
            continue
        }
		// term.terminal.Write([]byte(line + "\n"))
	}
	return nil
}

func (term *Term) AddCmd(cmd string, call CmdCall) {
    term.cmds[cmd] = call
}

/*
	term.AutoCompleteCallback =
		func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
			if key != 9 { // 非tab键不处理
				return "", 0, false
			}
			if pos == 0 {
				new := strings.Join(autoGroup, "\t")
				term.Write([]byte(new + "\n"))
				return "", 0, false
			}
			var tmp []string
			for _, v := range autoGroup {
				if strings.HasPrefix(v, line) {
					tmp = append(tmp, v)
				}
			}
			if len(tmp) == 0 {
				return line, pos, false
			} else if len(tmp) == 1 {
				return tmp[0], len(tmp[0]), true
			}
			new := strings.Join(tmp, "\t")
			term.Write([]byte(new + "\n"))
			return "", 0, false
		}
*/
