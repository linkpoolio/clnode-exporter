package client

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/linkpoolio/clnode-exporter/model"
	"bytes"
	log "github.com/sirupsen/logrus"
)

// NodeClient stores the credentials/hostname and an established session cookie
type NodeClient struct {
	Config *model.NodeConfig
	Cookie *http.Cookie
}

// NewNodeClient returns a new NodeConfig with an established session cookie
func NewNodeClient(config *model.NodeConfig) (*NodeClient, error) {
	nc := &NodeClient{Config: config}
	err := nc.SetSessionCookie()
	if err != nil {
		return nil, err
	}
	return nc, nil
}

// HttpGet calls an endpoint on the node, attaching the session cookie
func (nc *NodeClient) HttpGet(endpoint string, excModel interface{}) error {
	client := &http.Client{}
	var resp *http.Response

	i := 0
	for {
		nc.debugHttp(endpoint)
		req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", nc.Config.URL, endpoint), nil)
		if err != nil {
			return err
		}
		req.AddCookie(nc.Cookie)
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode == 401 && i < 5 {
			i++
			err = nc.SetSessionCookie()
			if err != nil {
				return err
			}
		} else {
			if i >= 5 {
				return fmt.Errorf("session returned from the node was invalid, attempted 5 times")
			}
			break
		}
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, excModel)
	if err != nil {
		return err
	}
	return nil
}

// SetSessionCookie establishes a new session on the node and retains the cookie to the NodeClient
func (nc *NodeClient) SetSessionCookie() error {
	session := &model.Session{Email: nc.Config.Username, Password: nc.Config.Password}
	b, err := json.Marshal(session)
	if err != nil {
		return err
	}
	nc.debugHttp("/sessions")
	resp, err := http.Post(
		fmt.Sprintf("%s/sessions", nc.Config.URL),
		"application/json",
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code of %d was returned when trying to get a session", resp.StatusCode)
	}
	if len(resp.Cookies()) == 0 {
		return fmt.Errorf("no cookie was returned after getting a session")
	}
	nc.Cookie = resp.Cookies()[0]
	return nil
}

func (nc *NodeClient) debugHttp(endpoint string) {
	log.WithFields(log.Fields{
		"url": nc.Config.URL,
		"endpoint": endpoint,
	}).Debug("node api call")
}