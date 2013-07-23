package qq

import (
	"../lib/go.net/websocket"
	"encoding/json"
	"fmt"
	"sync"
	"log"
	"strconv"
)

type Group struct {
	sync.Mutex
	Clients map[int]*Client
	MaxClients  int
	CurrentId   int
}

func NewGroup(max int) *Group {
	group := new(Group)
	group.Clients = make(map[int]*Client, max)
	group.MaxClients = max
	return group
}

func (g *Group) Broadcast(content string, fromId int) {
	for id, client := range g.Clients {
		if id == fromId {
			continue
		} else {
			client.Write(content)
		}
	}
}

func (g *Group) AddClient(ws *websocket.Conn) *Client {
	if len(g.Clients) > g.MaxClients {
		return nil
	}
	g.Lock()
	g.CurrentId++
	c := new(Client)
	c.Id = g.CurrentId
	c.Conn = ws
	c.Addr = ws.LocalAddr().String()
	c.CurGroup = g
	g.Clients[c.Id] = c
	g.Unlock()
	c.On_Login()
	return c
}

func (g *Group) RemoveClient(clientId int) bool {
	c := g.Clients[clientId]
	if c == nil {
		return false
	}
	c.On_Logout()
	g.Lock()
	delete(g.Clients, clientId)
	g.Unlock()
	return true
}

func (g *Group) Sendto(formId, toId int, content string) {
	c := g.Clients[toId]
	if c== nil {
		log.Println("User[#%d] is bad.", toId)
		return
	}
	c.Write(fmt.Sprintf("msgfrom %d %s", formId, content))
}

func (g *Group) GetClients() string {
	var jsonStr []byte
	names := make(map[string]UserInfo, len(g.Clients))
	i := 0
	for _, v := range g.Clients {
		id := strconv.Itoa(v.Id)
		names[id] = v.Info
		i++
	}
	jsonStr, err := json.Marshal(names)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(jsonStr)
}
