package ami

import "net"

func Open(host, login, password string) (*Client, error) {

}

type Client struct {
	conn net.Conn
}

func (self *Client) connect() {
	if self.conn != nil {
		return
	}
	s.conn = 
}
