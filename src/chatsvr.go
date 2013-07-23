package main

import (
	"./etc"
	"./lib/go.net/websocket"
	"./qq"
	"fmt"
	"io"
	"log"
	"net/http"
	//"reflect"
	"strconv"
	. "strings"
)

func MainLoop(ws *websocket.Conn) {
	var cmd []string = make([]string, 3)
	var to_id int

	client := G.AddClient(ws)

	if client == nil {
		log.Printf("Too many connection. Addr=%s", client.Addr)
		return
	}

	//	defer func() {
	//		if x := recover(); x != nil {
	//			log.Printf("run time panic: %v", x)
	//		}
	//	}()

Loop:
	for {
		var reply string
		if err := websocket.Message.Receive(ws, &reply); err != nil {
			log.Printf("Receive Error: %s", err)
			break
		}
		if etc.Debug {
			log.Printf("%s: %s", client.Addr, reply)
		}
		if EqualFold("", TrimSpace(reply)) {
			fmt.Println("empty")
			break
		}
		cmd = SplitN(reply, " ", 3)
		//至少有2个参数
		if len(cmd) < 3 {
			goto unknow
		}
		switch cmd[0] {
		case "set":
			switch cmd[1] {
			//修改名称
			case "name":
				info := SplitN(cmd[2], "|", 2)
				client.SetName(info[0], info[1])
				continue
				break
			default:
				goto unknow
			}
			break
		case "get":
			switch cmd[1] {
			case "name":
				client.Write("names " + G.GetClients())
				break
			default:
				goto unknow
			}
			break

		//向某个人发送信息
		case "sendto":
			to_id, _ = strconv.Atoi(cmd[1])
			G.Sendto(client.Id, to_id, cmd[2])
			break

		//向某个组或所有人发送信息
		case "sendm":
			to_id, _ = strconv.Atoi(cmd[1])
			G.Broadcast(fmt.Sprintf("msgall %d %s", client.Id, cmd[2]), client.Id)
			break
		case "logout":
			G.RemoveClient(client.Id)
			break Loop

		default:
			goto unknow
		}
	unknow:
		if err := websocket.Message.Send(ws, "cmd unknow"); err != nil {
			log.Printf("Can't send.%s", err)
			panic(100)
		}
	}
	G.RemoveClient(client.Id)
}

var G *qq.Group

func main() {
	//http.Handle("/xhr_poll/", )
	G = qq.NewGroup(etc.MaxClient)
	http.Handle("/chatroom", websocket.Handler(MainLoop))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", http.RedirectHandler("/static/html/index.html", 301))
	http.HandleFunc("/config.js", func(w http.ResponseWriter, r *http.Request) {
		var js_debug string = "false"
		if etc.Debug {
			js_debug = "true"
		}
		io.WriteString(w, fmt.Sprintf("var config = {host:'%s', port: %s, debug: %s};\n", etc.ServerHost, etc.ServerPort, js_debug))
	})
	var svr_ = etc.ServerHost + ":" + etc.ServerPort
	log.Printf("Server start on %s", svr_)
	if err := http.ListenAndServe(svr_, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
