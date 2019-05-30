// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"github.com/greymd/ojichat/generator"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
	"strings"
)

func getSenderID(event *linebot.Event) string {
	switch event.Source.Type {
	case linebot.EventSourceTypeGroup:
		return event.Source.GroupID
	case linebot.EventSourceTypeRoom:
		return event.Source.RoomID
	case linebot.EventSourceTypeUser:
		return event.Source.UserID
	}
	return ""
}

func getSenderName(bot *linebot.Client, from string) string {
	if len(from) == 0 {
		return ""
	}
	if from[0:1] == "U" {
		senderProfile, err := bot.GetProfile(from).Do()
		if err != nil {
			return ""
		}
		return senderProfile.DisplayName
	}
	return ""
}

func getSenderName2(bot *linebot.Client, event *linebot.Event) string {
	switch event.Source.Type {
	case linebot.EventSourceTypeGroup:
		senderProfile, err2 := bot.GetGroupMemberProfile(event.Source.GroupID, event.Source.UserID).Do()
		if err2 != nil {
			return ""
		} else {
			return senderProfile.DisplayName
		}
	case linebot.EventSourceTypeRoom:
		senderProfile, err2 := bot.GetRoomMemberProfile(event.Source.RoomID, event.Source.UserID).Do()
		if err2 != nil {
			return ""
		} else {
			return senderProfile.DisplayName
		}
	}
	return getSenderName(bot, getSenderID(event))
}

func main() {
	bot, err := linebot.New(
		os.Getenv("OJILINEBOT_CHANNEL_SECRET"),
		os.Getenv("OJILINEBOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		//fmt.Printf("Hello world %s\n",req.Body)
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					fmt.Printf("%s\n", message.Text)
					if strings.Contains(message.Text, "おじさん") {
						config := generator.Config{}
						config.TargetName = getSenderName2(bot, event)
						config.EmojiNum = 4
						config.PunctiuationLebel = 0
						result, err := generator.Start(config)
						fmt.Printf("%s\n", result)

						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(result)).Do(); err != nil {
							log.Print(err)
						}
					}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+os.Getenv("OJILINEBOT_PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
