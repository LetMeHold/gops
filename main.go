package main

import (
	"gops/api"
	"log"
	"strings"
)

func Print(term *api.Term, args []string) error {
	term.WriteString(strings.Join(args, " ") + "\n")
	return nil
}

func main() {
	term, err := api.NewTerm("->> ")
	if err != nil {
		log.Fatalf("api.NewTerm faild : %v", err)
	}
	defer term.Close()
	term.AddCmd("print", Print)
	if err = term.Start(); err != nil {
		log.Fatalf("term.Start faild : %v", err)
	}
	/*
		rmt, err := api.NewRemoteByDefaultKey("172.18.180.110:22", "root")
		if err != nil {
			log.Fatalf("NewRemote failed : %v", err)
		}
		defer rmt.Close()
		err = rmt.Run("whoami")
		if err != nil {
			log.Fatalf("Run failed : %v", err)
		}
		err = rmt.Start("ping baidu.com -c 3")
		if err != nil {
			log.Fatalf("Start failed : %v", err)
	*/
}
