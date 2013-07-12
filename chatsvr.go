package main

import (
	"./lib/go.net/websocket"
	"fmt"
	"io"
	"log"
	"net/http"
	//"reflect"
	"runtime"
	"strconv"
	. "strings"
	"sync"
	"net"
	"encoding/json"
)

type Client struct {
	conn *websocket.Conn
}

type Config struct {
	host string
	port string
	max_conn int
}

var svr_cnf Config = Config{port: "80", max_conn: 1000}

var cur_id int = 0
var mutex sync.Mutex
var err error

interface ChatUser{
	Write(string)
}

func (Client *) Write(content string) {

}

var clients map[int]Client = make(map[int]Client, svr_cnf.max_conn)
var userlist map[string]string =  make(map[string]string, svr_cnf.max_conn)

func broadcast(content string, client_id int) {
	for id, client := range(clients) {
		if id == client_id {
			continue
		} else {
			websocket.Message.Send(client.conn, content)
		}
	}
}

func cmd_set_name(client_id int, name string) {
	var ws *websocket.Conn = clients[client_id].conn
	mutex.Lock()
	userlist[strconv.Itoa(client_id)] = name
	mutex.Unlock()
	go broadcast(fmt.Sprintf("login %d %s", client_id, name), client_id)
	websocket.Message.Send(ws, fmt.Sprintf("setok %d", client_id))
}

func cmd_unknow(client_id int) {
	var ws *websocket.Conn = clients[client_id].conn
	if err = websocket.Message.Send(ws, "cmd unknow"); err != nil {
		log.Printf("Can't send.%s", err)
		panic(100)
	}
}

func cmd_proxy() {

}

func cmd_get_name(client_id int) {	
	var jsonStr []byte
	
	jsonStr, err := json.Marshal(userlist)
	if err != nil {
		fmt.Println("error:", err)
	}
	var ws *websocket.Conn = clients[client_id].conn
	websocket.Message.Send(ws, fmt.Sprintf("names %s", jsonStr))
}

func cmd_sendto(client_id int, to_id int, content string) {
	var ws *websocket.Conn = clients[to_id].conn
	websocket.Message.Send(ws, fmt.Sprintf("msgto %d %s", client_id, content))
}

func cmd_sendall(client_id int, group_id int, content string) {
	var send string = fmt.Sprintf("msgall %d %s", client_id, content);
	go broadcast(send, client_id);
}

func comet_loop() {

}

func CharRoom(ws *websocket.Conn) {
	var cmd []string = make([]string, 3)
	var addr = ws.LocalAddr().String()
	var client_id int
	var to_id int

	mutex.Lock()
	cur_id++
	client_id = cur_id
	clients[client_id] = Client{conn: ws}
	mutex.Unlock()

	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic: %v", x)
		}
	}()

	log.Printf("New connection. Addr=%s", addr)

	Loop:
	for {
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			log.Printf("Can't send.%s", err)
			break
		}

		if EqualFold("", TrimSpace(reply)) {
			fmt.Println("empty")
			break
		}

		cmd = SplitN(reply, " ", 3)

		//至少有2个参数
		if len(cmd) < 3 {
			cmd_unknow(client_id)
		}

		log.Printf("%s: %s-%s: %s", addr, cmd[0], cmd[1], cmd[2])
		switch cmd[0] {
		case "set":
			switch cmd[1] {
			//修改名称
			case "name":
				cmd_set_name(client_id, cmd[2])
				continue
				break
			default:
				cmd_unknow(client_id)
				break
			}
			break
		case "get":
			switch cmd[1] {
			case "name":
				cmd_get_name(client_id)
				break
			default:
				cmd_unknow(client_id)
				break
			}
			break
		
		//向某个人发送信息
		case "sendto":
			to_id, _ = strconv.Atoi(cmd[1])
			cmd_sendto(client_id, to_id, cmd[2])
			break

		//向某个组或所有人发送信息
		case "sendm":
			to_id, _ = strconv.Atoi(cmd[1])
			cmd_sendall(client_id, to_id, cmd[2])
			break
		case "logout":
			websocket.Message.Send(ws, fmt.Sprintf("goodbye %d", client_id))
			break Loop;	
		
		default:
			cmd_unknow(client_id)
			break
		}
	}
	
	mutex.Lock()
	delete(clients, client_id)
	delete(userlist, strconv.Itoa(client_id))
	broadcast("logout "+strconv.Itoa(client_id), client_id)
	mutex.Unlock()
}

func main() {
	ifs, _ := net.InterfaceAddrs() 
	svr_cnf.host = ifs[0].String()
	
	http.Handle("/xhr_poll/", )
	http.Handle("/chatroom", websocket.Handler(CharRoom))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", http.RedirectHandler("/static/html/index.html", 301))
	http.HandleFunc("/config.js", func(w http.ResponseWriter, r *http.Request){
		var conf map[string]string = make(map[string]string, 2)
		conf["host"] = svr_cnf.host
		conf["port"] = svr_cnf.port
		jsonStr, err := json.Marshal(conf); if err!=nil {
			fmt.Printf("json.Marshal fail.error=%s\n", err)
		}
		io.WriteString(w, fmt.Sprintf("var config = %s;\n", jsonStr))
	})

	runtime.GOMAXPROCS(4)

	var svr_ = svr_cnf.host + ":" + svr_cnf.port

	log.Printf("Server start on %s", svr_)
	if err := http.ListenAndServe(svr_, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
