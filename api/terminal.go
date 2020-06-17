package api

import (
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"sort"
	"strings"
)

type Completion map[string][]string

type CmdCall func(term *Term, args []string) error

type Term struct {
	completion Completion
	current    string
	terminal   *terminal.Terminal
	fd         int
	state      *terminal.State
	cmds       map[string]CmdCall
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
	term.terminal.AutoCompleteCallback = term.autoComplete
	term.cmds = make(map[string]CmdCall)
	term.completion = make(Completion)
	term.completion["default"] = []string{"exit"}
	term.current = "default"
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
		input := strings.TrimSpace(line)
		if input == "exit" {
			break
		}
		args := strings.Fields(input)
		if call, ok := term.cmds[args[0]]; ok {
			call(term, args[1:])
			continue
		}
	}
	return nil
}

func (term *Term) AddCmd(cmd string, call CmdCall) {
	term.cmds[cmd] = call
	term.completion["default"] = append(term.completion["default"], cmd)
}

func (term *Term) autoComplete(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	if key != 9 { // 非tab键不处理
		return "", 0, false
	}
	if pos == 0 { // 输入为空
		term.WriteString(strings.Join(term.completion[term.current], "\t") + "\n")
		return "", 0, false
	}
	var tmp []string
	for _, word := range term.completion[term.current] {
		if strings.HasPrefix(word, line) {
			tmp = append(tmp, word)
		}
	}
	if len(tmp) == 0 { // 没有匹配的字符串
		return "", 0, false
	} else if len(tmp) == 1 { // 只有一个匹配
		return tmp[0] + " ", len(tmp[0]) + 1, true
	}
	sort.Strings(tmp)
	prefix, err := GetMaxPrefix(tmp) // 获取最长的相同前缀
	if err != nil {
		term.WriteString(err.Error() + "\n")
		return "", 0, false
	}
	if prefix == line { // 需要补全的内容与输入相同
		term.WriteString(strings.Join(tmp, "\t") + "\n")
		return "", 0, false
	}
	return prefix, len(prefix), true
}
