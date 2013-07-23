package qq

import (
	"fmt"
	"strconv"
)
func (c *Client) On_Login() {
//	fmt.Printf("%s: login %d", c.Addr, c.Id)
}

func (c *Client) On_SetName() {
	c.Write(fmt.Sprintf("setok %d", c.Id))
	s := fmt.Sprintf("login %d %s|%s",  c.Id, c.Info.Name, c.Info.Avatar)
	go c.CurGroup.Broadcast(s, c.Id)
}

func (c *Client) On_Logout() {
	c.Write(fmt.Sprintf("goodbye %d", c.Id))
	s := "logout "+strconv.Itoa(c.Id)
	go c.CurGroup.Broadcast(s, c.Id)
}
