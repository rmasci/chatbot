package main

import (
	"fmt"
	"github.com/aichaos/rivescript-go"
	"github.com/bitfield/script"
	"log"
	"strings"
)

func isChatbot(chatbot, message string) string {
	mesgSlice := strings.Fields(message)
	firstWord := mesgSlice[0]
	msg := strings.Join(mesgSlice[1:], " ")
	if firstWord == chatbot || firstWord == "@+chatbot" {
		return msg
	}
	return ""
}

func getReply(bot *rivescript.RiveScript, username, message string) (string, error) {
	reply, err := bot.Reply(username, message)
	if err != nil {
		return "", err
	}
	replySlice := strings.Fields(reply)
	switch replySlice[0] {
	case "LocalCommand":
		cmd := strings.Join(replySlice[1:], " ")
		log.Println("LocalCommand:", cmd)
		reply, err = script.Exec(cmd).String()
		if err != nil {
			return "", err
		}
		retStr := fmt.Sprintf("</div><div class=\"col-10 mb-1 small\"><pre>%s</pre></div></div>", reply)
		return retStr, nil
	}
	return reply, nil
}
