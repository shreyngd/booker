package slotserver

import (
	"encoding/json"
	"log"

	"github.com/shreyngd/booker/models"
)

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const AllRoomsList = "all-rooms"

type Message struct {
	Action  string  `json:"action"`
	Message string  `json:"message"`
	Target  string  `json:"target"`
	Sender  *Client `json:"sender"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}

type MessageAddRoom struct {
	Data   models.AddRoomReply `json:"data"`
	Action string              `json:"action"`
}

func (message *MessageAddRoom) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}
