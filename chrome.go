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
  "strconv"
  "strings"
  "sync"
  "time"

  "github.com/kwf2030/commons/base"
)

var (
  ErrLaunchChrome = errors.New("failed to launch chrome")
  ErrCreateTab    = errors.New("failed to create tab")
)

var (
  ArgDebugPort  = "--remote-debugging-port"
  ArgHeadless   = "--headless"
  ArgIgnoreCert = "--ignore-certificate-errors"
)

type Chrome struct {
  // ChromeDevToolsProtocol的Endpoint（http://host:port/json），
  // 请求该地址返回的是Meta数组
  Endpoint string
  Host     string
  Port     int
  Process  *os.Process
}

func Launch(bin string, args ...string) (*Chrome, error) {
  if bin == "" {
    return nil, base.ErrInvalidArgument
  }
  _, e := exec.LookPath(bin)
  if e != nil {
    return nil, e
  }
  var port string
  for _, arg := range args {
    if strings.Contains(arg, ArgDebugPort) {
      arr := strings.Split(arg, "=")
      if len(arr) != 2 {
        return nil, base.ErrInvalidArgument
      }
      port = strings.TrimSpace(arr[1])
      break
    }
  }
  if port == "" {
    port = "9222"
    args = append(args, ArgDebugPort+"="+port)
  }
  cmd := exec.Command(bin, args...)
  e = cmd.Start()
  if e != nil {
    return nil, e
  }
  h := "127.0.0.1"
  p, _ := strconv.Atoi(port)
  c := &Chrome{Endpoint: fmt.Sprintf("http://%s:%d/json", h, p), Host: h, Port: p, Process: cmd.Process}
  if ok := c.waitForStarted(time.Second * 10); !ok {
    return nil, ErrLaunchChrome
  }
  return c, nil
}

func Connect(host string, port int) (*Chrome, error) {
  if host == "" || port <= 0 {
    return nil, errors.New("param invalid")
  }
  endpoint := fmt.Sprintf("http://%s:%d/json", host, port)
  resp, e := http.Get(endpoint)
  if e != nil {
    return nil, e
  }
  drain(resp.Body)
  return &Chrome{Endpoint: endpoint, Host: host, Port: port, Process: nil}, nil
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
    return nil, ErrCreateTab
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
    r.Close()
  }
}
