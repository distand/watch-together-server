package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0:6506", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

type Store struct {
	Url       string `json:"u,omitempty"`
	Time      string `json:"t,omitempty"`
	Timestamp string `json:"ts,omitempty"`
}

type In struct {
	SetUrl       string `json:"su"`
	SetTime      string `json:"st"`
	SetTimestamp string `json:"ts"`
	GetUrl       string `json:"gu"`
	GetTime      string `json:"gt"`
}

var store *Store

func server(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("in: %s", message)
		var in In
		if err = json.Unmarshal([]byte(message), &in); err != nil {
			log.Println("err:", err)
			break
		}
		out := dealMsg(in)
		log.Printf("out: %s", string(out))
		err = c.WriteMessage(mt, out)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func dealMsg(p In) []byte {
	res := Store{}
	if p.SetUrl != "" {
		store.Url = p.SetUrl
	}
	if p.SetTime != "" {
		store.Time = p.SetTime
	}
	if p.SetTimestamp != "" {
		store.Timestamp = p.SetTimestamp
	}
	if p.GetUrl != "" {
		res.Url = store.Url
	}
	if p.GetTime != "" {
		res.Time = store.Time
		res.Timestamp = store.Timestamp
	}
	s, _ := json.Marshal(res)
	return s
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", server)
	store = &Store{}
	log.Fatal(http.ListenAndServe(*addr, nil))
}
