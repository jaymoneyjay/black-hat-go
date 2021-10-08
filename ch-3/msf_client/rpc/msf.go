package rpc

import (
	"bytes"
	"fmt"
	"gopkg.in/vmihailenco/msgpack.v2"
	"net/http"
)

type sessionListReq struct {
	_msgpack	struct{}	`msgpack:",asArray"`
	Method 		string
	Token 		string
}
type SessionListRes struct {
	ID				uint32	`msgpack:",omitempty"`
	Type			string	`msgpack:"type"`
	TunnelLocal		string	`msgpack:"tunnel_local"`
	TunnelPeer		string	`msgpack:"tunnel_peer"`
	ViaExploit		string	`msgpack:"via_exploit"`
	ViaPayload		string	`msgpack:"via_payload"`
	Description		string	`msgpack:"desc"`
	Info			string	`msgpack:"info"`
	Workspace		string	`msgpack:"workspace"`
	TargetHost		string	`msgpack:"target_host"`
	Username		string	`msgpack:"username"`
	UUID 			string	`msgpack:"uuid"`
	ExploitUUID		string	`msgpack:"exploit_uuid"`
}
type loginReq struct {
	_msgpack	struct{}	`msgpack:",asArray"`
	Method 		string
	Username	string
	Password 	string

}
type loginRes struct {
	Result 			string	`msgpack:"result"`
	Token			string	`msgpack:"token"`
	Error 			bool	`msgpack:"error"`
	ErrorClass		string	`msgpack:"error_class"`
	ErrorMessage	string	`msgpack:"error_message"`
}
type logoutReq struct {
	_msgpack	struct{}	`msgpack:",asArray"`
	Method		string
	Token		string
	LogoutToken	string
}
type logoutRes struct {
	Result	string	`msgpack:"result"`
}

// Stores authentication information for API call
type Metasploit struct {
	host		string
	user		string
	password 	string
	token 		string
}

// New authenticates a user and returns a new Metasploit session
func New(host, user, password string) (*Metasploit, error) {
	m := &Metasploit{
		host: host,
		user: user,
		password: password,
		token: "",
	}
	if err := m.Login(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Metasploit) send(req, res interface{}) error {
	buf := new(bytes.Buffer)
	msgpack.NewEncoder(buf).Encode(req)

	dest := fmt.Sprintf("http://%s/api", m.host)
	r, err := http.Post(dest, "binary/message-pack", buf)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	if err := msgpack.NewDecoder(r.Body).Decode(&res); err != nil {
		return err
	}
	return nil
}

// Login calls the auth.login method
func (m *Metasploit) Login() error {
	req := loginReq{
		Method:   	"auth.login",
		Username: 	m.user,
		Password: 	m.password,
	}
	var res loginRes
	if err := m.send(&req, &res); err != nil {
		return err
	}
	// TODO: Add error handling for response
	m.token = res.Token
	return nil
}

// Logout calls the auth.logout method
func (m *Metasploit) Logout() error {
	req := logoutReq{
		Method: 		"auth.logout",
		Token: 			m.token,
		LogoutToken: 	m.token,
	}
	var res logoutRes
	if err := m.send(&req, &res); err != nil {
		return nil
	}
	// TODO: Add error handling for response
	m.token = ""
	return nil
}

// SessionList calls the session.list method
func (m *Metasploit) SessionList() (map[uint32]SessionListRes, error) {
	req := sessionListReq{
		Method: "session.list",
		Token: 	m.token,
	}
	res := make(map[uint32]SessionListRes)
	if err := m.send(&req, &res); err != nil {
		return nil, err
	}
	// TODO: Add error handling for response

	for id, session := range res {
		session.ID = id
		res[id] = session
	}
	return res, nil
}


