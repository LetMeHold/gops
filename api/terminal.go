package api

import (
	"golang.org/x/crypto/ssh/terminal"
	"io"
	//"log"
	"os"
	//"strings"
)

var autoGroup = []string{
	"hello, world",
	"hello, china",
	"golang",
	"goto",
}

type Completion []string

type Term struct {
	completions map[string]*Completion
	terminal        *terminal.Terminal
	fd          int
	state       *terminal.State
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
	return term, nil
}

func (term *Term) Close() error {
	return terminal.Restore(term.fd, term.state)
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
		term.terminal.Write([]byte(line + "\n"))
	}
	return nil
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
