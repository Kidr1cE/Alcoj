package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有 Origin
	},
}

type Request struct {
	Key string `json:"key"`
}

type Response struct {
	WorkerNum     int      `json:"worker_num"`
	QueueTasks    int      `json:"queue_tasks"`
	FinishedTasks int      `json:"finished_tasks"`
	Workers       []Worker `json:"workers"`
}

type Worker struct {
	WorkerID     string `json:"worker_id"`
	WorkerStatus string `json:"worker_status"`
}

var response = Response{
	WorkerNum:     1,
	QueueTasks:    0,
	FinishedTasks: 0,
	Workers: []Worker{
		{
			WorkerID:     "64f50a0c-f0a4-11ee-90ff-00155d91e788",
			WorkerStatus: "running",
		},
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed: ", err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage failed: ", err)
			return
		}
		key := string(p)
		log.Println("Received key:", key)

		response.QueueTasks++
		response.FinishedTasks++
		content, err := json.Marshal(response)
		if err != nil {
			log.Println("json.Marshal failed: ", err)
			return
		}

		if err := conn.WriteMessage(messageType, content); err != nil {
			log.Println("WriteMessage failed: ", err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/python3", handler)
	log.Println("Listening on :7070")
	err := http.ListenAndServe(":7070", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
