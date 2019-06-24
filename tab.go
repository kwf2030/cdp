package cdp

import (
  "net/http"
  "sync"
  "sync/atomic"

  "github.com/gorilla/websocket"
)

type Handler interface {
  OnCdpEvent(*Message)

  OnCdpResponse(*Message) bool
}

// 请求/响应/事件通知
type Message struct {
  // 请求的Id，响应中会带有相同的Id，每次请求Tab.lastMessageId自增后赋值给Message.Id，
  // 事件通知没有该字段
  Id int32 `json:"id,omitempty"`

  // 请求、响应和事件通知都有该字段
  Method string `json:"method,omitempty"`

  // 请求的参数（可选）、事件通知的数据（可选），
  // 响应没有该字段
  Params map[string]interface{} `json:"params,omitempty"`

  // 响应数据（请求和事件通知没有该字段）
  Result map[string]interface{} `json:"result,omitempty"`

  // 同步等待channel，仅在Handler为nil或Handler.OnCdpResponse()返回false的时候会发送一次，
  // 发送的是Message自身
  syncChan chan *Message

  // 一些自定义属性
  What int         `json:"-"`
  Arg  int         `json:"-"`
  Str  string      `json:"-"`
  Obj  interface{} `json:"-"`
}

type Meta struct {
  Id                   string `json:"id"`
  Type                 string `json:"type"`
  Title                string `json:"title"`
  Url                  string `json:"url"`
  FaviconUrl           string `json:"faviconUrl"`
  Description          string `json:"description"`
  DevtoolsFrontendUrl  string `json:"devtoolsFrontendUrl"`
  WebSocketDebuggerUrl string `json:"webSocketDebuggerUrl"`
}

type Tab struct {
  chrome *Chrome

  meta *Meta

  conn *websocket.Conn

  // 每次请求自增
  lastMessageId int32

  // 非零表示Tab已经关闭
  closed int32

  // 广播，用于通知WebSocket关闭读goroutine
  closeChan chan struct{}

  handler Handler

  // 存放两类数据：
  // 1.订阅的事件（string-->bool），key是Message.Method，用于过滤WebSocket读取到的事件，
  // 2.请求的Message（int32-->*Message），key是Message.Id，用于读取到数据时找到对应的请求Message
  data sync.Map
}

func (t *Tab) connect() (*websocket.Conn, error) {
  conn, _, e := websocket.DefaultDialer.Dial(t.meta.WebSocketDebuggerUrl, nil)
  if e != nil {
    return nil, e
  }
  return conn, nil
}

func (t *Tab) read() {
  for {
    select {
    case <-t.closeChan:
      return

    default:
      msg := &Message{}
      e := t.conn.ReadJSON(msg)
      if e != nil {
        t.Close()
        return
      }
      t.dispatch(msg)
    }
  }
}

func (t *Tab) dispatch(msg *Message) {
  // 事件通知
  if msg.Id == 0 {
    // 若注册过该类事件，则进行通知
    if _, ok := t.data.Load(msg.Method); ok && t.handler != nil {
      go t.handler.OnCdpEvent(msg)
    }
    return
  }
  // Message.id非0表示响应
  if v, ok := t.data.Load(msg.Id); ok {
    t.data.Delete(msg.Id)
    req := v.(*Message)
    // 响应没有method字段，
    // 把响应的数据赋值给对应的请求，回调用req作为参数（省去给msg的字段逐个赋值了）
    req.Result = msg.Result
    go func() {
      if t.handler == nil || !t.handler.OnCdpResponse(req) {
        req.syncChan <- req
      }
    }()
  }
}

func (t *Tab) FireEvent(event string, params map[string]interface{}) {
  if t.handler != nil {
    go t.handler.OnCdpEvent(&Message{Method: event, Params: params})
  }
}

func (t *Tab) Call(method string, params map[string]interface{}) (int32, chan *Message) {
  return t.CallAttr(method, params, 0, 0, "", nil)
}

func (t *Tab) CallAttr(method string, params map[string]interface{}, what, arg int, str string, obj interface{}) (int32, chan *Message) {
  if method == "" {
    return 0, nil
  }
  id := atomic.AddInt32(&t.lastMessageId, 1)
  ch := make(chan *Message, 1)
  msg := &Message{
    Id:       id,
    Method:   method,
    Params:   params,
    syncChan: ch,
    What:     what,
    Arg:      arg,
    Str:      str,
    Obj:      obj,
  }
  t.data.Store(id, msg)
  e := t.conn.WriteJSON(msg)
  if e != nil {
    t.Close()
    return 0, nil
  }
  return id, ch
}

func (t *Tab) Subscribe(events ...string) {
  for _, evt := range events {
    if evt != "" {
      t.data.Store(evt, true)
    }
  }
}

func (t *Tab) Unsubscribe(events ...string) {
  for _, evt := range events {
    if evt != "" {
      t.data.Delete(evt)
    }
  }
}

func (t *Tab) Activate() {
  resp, e := http.Get(t.chrome.Endpoint + "/activate/" + t.meta.Id)
  if e == nil {
    drain(resp.Body)
  }
}

func (t *Tab) Close() {
  // 调用一次Close后把Tab.closed标识设为1，防止多次调用
  if !atomic.CompareAndSwapInt32(&t.closed, 0, 1) {
    return
  }
  close(t.closeChan)
  t.conn.Close()
  resp, e := http.Get(t.chrome.Endpoint + "/close/" + t.meta.Id)
  if e == nil {
    drain(resp.Body)
  }
}

func (t *Tab) Closed() bool {
  return atomic.LoadInt32(&t.closed) != 0
}
