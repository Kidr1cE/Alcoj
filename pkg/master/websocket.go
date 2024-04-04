package master

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var websocketPort = os.Getenv("WEBSOCKET_PORT")

type BackendRequest struct {
	Key string `json:"key"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type BackendResponse struct {
	WorkerNum     int64            `json:"worker_num"`
	QueueTasks    int64            `json:"queue_tasks"`
	FinishedTasks int64            `json:"finished_tasks"`
	Workers       []*SandboxServer `json:"workers"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed: ", err)
		return
	}
	defer conn.Close()

	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage failed: ", err)
			return
		}
		response := BackendResponse{
			WorkerNum:     master.WorkerNum.Load(),
			QueueTasks:    master.QueueTasks.Load(),
			FinishedTasks: master.FinishedTasks.Load(),
			Workers:       master.Workers,
		}

		// Send master status to client
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

func startWebsocketServer(stopCh chan struct{}) {
	defer func() {
		stopCh <- struct{}{}
	}()

	http.HandleFunc("/"+language, handler)
	log.Printf("Listening on :%s", websocketPort)
	err := http.ListenAndServe(":"+websocketPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
