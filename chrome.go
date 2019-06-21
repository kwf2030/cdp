package cdp

import (
  "encoding/json"
  "errors"
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "os"
  "os/exec"
  "strings"
  "sync"
  "time"
)

type Chrome struct {
  // ChromeDevToolsProtocol的Endpoint（http://host:port/json）
  Endpoint string
  Process  *os.Process
}

func Launch(bin string, args ...string) (*Chrome, error) {
  if bin == "" {
    return nil, errors.New("param <bin> is empty")
  }
  _, e := exec.LookPath(bin)
  if e != nil {
    return nil, e
  }
  var port string
  for _, arg := range args {
    if strings.Contains(arg, "--remote-debugging-port") {
      arr := strings.Split(arg, "=")
      if len(arr) != 2 {
        return nil, errors.New("param <args> invalid")
      }
      port = strings.TrimSpace(arr[1])
      break
    }
  }
  if port == "" {
    port = "9222"
    args = append(args, "--remote-debugging-port="+port)
  }
  cmd := exec.Command(bin, args...)
  e = cmd.Start()
  if e != nil {
    return nil, e
  }
  c := &Chrome{"http://127.0.0.1:" + port + "/json", cmd.Process}
  if ok := c.waitForStarted(time.Second * 10); !ok {
    return nil, errors.New("failed to launch chrome")
  }
  return c, nil
}

func Connect(host string, port int) (*Chrome, error) {
  if host == "" || port <= 0 {
    return nil, errors.New("param invalid")
  }
  return &Chrome{fmt.Sprintf("http://%s:%d/json", host, port), nil}, nil
}

func (c *Chrome) Exit() error {
  tab, e := c.NewTab(nil)
  if e != nil {
    return e
  }
  tab.Call(Browser.Close, nil)
  return nil
}

func (c *Chrome) NewTab(h Handler) (*Tab, error) {
  meta := &Meta{}
  resp, e := http.Get(c.Endpoint + "/new")
  if e != nil {
    return nil, e
  }
  e = json.NewDecoder(resp.Body).Decode(meta)
  if e != nil {
    return nil, e
  }
  e = resp.Body.Close()
  if e != nil {
    return nil, e
  }
  if meta.Id == "" || meta.WebSocketDebuggerUrl == "" {
    return nil, errors.New("failed to create tab")
  }
  t := &Tab{
    chrome:    c,
    meta:      meta,
    closeChan: make(chan struct{}),
    handler:   h,
    data:      sync.Map{},
  }
  t.conn, e = t.connect()
  if e != nil {
    t.Close()
    return nil, e
  }
  go t.read()
  return t, nil
}

func (c *Chrome) waitForStarted(timeout time.Duration) bool {
  client := &http.Client{Timeout: time.Second}
  t := time.After(timeout)
  for {
    select {
    case <-t:
      return false
    default:
      resp, e := client.Get(c.Endpoint)
      if e != nil {
        break
      }
      drain(resp.Body)
      return true
    }
  }
}

func drain(r io.ReadCloser) {
  _, e := ioutil.ReadAll(r)
  if e == nil {
    _ = r.Close()
  }
}
