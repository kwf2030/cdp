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

  "github.com/kwf2030/commons/base"
)

var (
  ErrLaunchChrome = errors.New("failed to launch chrome")
  ErrCreateTab    = errors.New("failed to create tab")
)

type Chrome struct {
  // ChromeDevToolsProtocol的Endpoint（http://host:port/json）
  // [ {
  //    "description": "",
  //    "devtoolsFrontendUrl": "/devtools/inspector.html?ws=127.0.0.1:9222/devtools/page/5D5FE2210AF9A5DAFAA2D69159C6CD52",
  //    "id": "5D5FE2210AF9A5DAFAA2D69159C6CD52",
  //    "title": "新标签页",
  //    "type": "page",
  //    "url": "chrome://newtab/",
  //    "webSocketDebuggerUrl": "ws://127.0.0.1:9222/devtools/page/5D5FE2210AF9A5DAFAA2D69159C6CD52"
  // }, {
  //    "description": "",
  //    "devtoolsFrontendUrl": "/devtools/inspector.html?ws=127.0.0.1:9222/devtools/page/853C0E933FD62DAD9ABBDFC9C3C47084",
  //    "faviconUrl": "https://www.baidu.com/favicon.ico",
  //    "id": "853C0E933FD62DAD9ABBDFC9C3C47084",
  //    "title": "百度一下，你就知道",
  //    "type": "page",
  //    "url": "https://www.baidu.com/",
  //    "webSocketDebuggerUrl": "ws://127.0.0.1:9222/devtools/page/853C0E933FD62DAD9ABBDFC9C3C47084"
  // }, {
  //    "description": "",
  //    "devtoolsFrontendUrl": "/devtools/inspector.html?ws=127.0.0.1:9222/devtools/page/CF2C261EEFA71ACB7803D25CFE93386C",
  //    "faviconUrl": "https://www.jd.com/favicon.ico",
  //    "id": "CF2C261EEFA71ACB7803D25CFE93386C",
  //    "title": "京东(JD.COM)-正品低价、品质保障、配送及时、轻松购物！",
  //    "type": "page",
  //    "url": "https://www.jd.com/",
  //    "webSocketDebuggerUrl": "ws://127.0.0.1:9222/devtools/page/CF2C261EEFA71ACB7803D25CFE93386C"
  // } ]
  Endpoint string
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
    if strings.Contains(arg, "--remote-debugging-port") {
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
    args = append(args, "--remote-debugging-port="+port)
  }
  cmd := exec.Command(bin, args...)
  e = cmd.Start()
  if e != nil {
    return nil, e
  }
  c := &Chrome{"http://127.0.0.1:" + port + "/json", cmd.Process}
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
  return &Chrome{endpoint, nil}, nil
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
