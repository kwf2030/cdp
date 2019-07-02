package cdp

import (
  "time"
)

const defaultEvent = "__default__"

type Action interface {
  Method() string
  Params() map[string]interface{}
}

type action struct {
  method string
  params map[string]interface{}
}

func NewAction(method string, params map[string]interface{}) Action {
  if method == "" {
    return nil
  }
  return &action{method: method, params: params}
}

func (a *action) Method() string {
  return a.method
}

func (a *action) Params() map[string]interface{} {
  return a.params
}

type waitAction time.Duration

func (wa waitAction) Method() string {
  return "<wait>"
}

func (wa waitAction) Params() map[string]interface{} {
  return nil
}

type evalAction struct {
  *action
  expressions []string
}

func NewEvalAction(expressions ...string) Action {
  if len(expressions) == 0 {
    return nil
  }
  return &evalAction{
    action: &action{
      method: Runtime.Evaluate,
      params: map[string]interface{}{"returnByValue": true},
    },
    expressions: expressions,
  }
}

type Task struct {
  chrome *Chrome
  tab    *Tab

  // 一个Domain事件对应多个Action（DomainEvent-->[]Action），
  // 没有事件的Action的key为DEFAULT
  actions map[string][]Action

  // 当前事件（用于链式调用）
  evt string

  handler Handler
}

func NewTask(c *Chrome) *Task {
  if c == nil {
    return nil
  }
  t := &Task{
    chrome:  c,
    actions: make(map[string][]Action, 2),
    evt:     defaultEvent,
  }
  t.actions[defaultEvent] = make([]Action, 0, 2)
  return t
}

func (t *Task) OnCdpEvent(msg *Message) {
  if actions, ok := t.actions[msg.Method]; ok {
    for _, action := range actions {
      t.runAction(action)
    }
  } else {
    if t.handler != nil {
      t.handler.OnCdpEvent(msg)
    }
  }
}

func (t *Task) OnCdpResponse(msg *Message) bool {
  if t.handler != nil {
    return t.handler.OnCdpResponse(msg)
  }
  return true
}

func (t *Task) Finish() {
  t.tab.Close()
}

func (t *Task) Action(action Action) *Task {
  if action != nil {
    t.actions[t.evt] = append(t.actions[t.evt], action)
  }
  return t
}

func (t *Task) Until(event string) *Task {
  if event != "" {
    if _, ok := t.actions[event]; !ok {
      t.evt = event
      t.actions[event] = make([]Action, 0, 16)
    }
  }
  return t
}

func (t *Task) Wait(duration time.Duration) *Task {
  if duration > 0 {
    t.Action(waitAction(duration))
  }
  return t
}

func (t *Task) Run(h Handler) *Task {
  tab, e := t.chrome.NewTab(t)
  if e != nil {
    return t
  }
  t.tab = tab
  t.handler = h
  for event := range t.actions {
    if event != defaultEvent {
      tab.Subscribe(event)
    }
  }
  for _, action := range t.actions[defaultEvent] {
    t.runAction(action)
  }
  return t
}

func (t *Task) runAction(action Action) {
  switch a := action.(type) {
  case waitAction:
    time.Sleep(time.Duration(a))
  case *evalAction:
    for _, expr := range a.expressions {
      if expr != "" {
        a.params["expression"] = expr
        t.tab.Call(a.Method(), a.params)
      }
    }
  default:
    t.tab.Call(a.Method(), a.Params())
  }
}
