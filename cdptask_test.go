package cdp

import (
  "fmt"
  "testing"
  "time"
)

type HTask struct {
  name string
}

func (h *HTask) OnCdpEvent(msg *Message) {
  fmt.Println("==========OnCdpEvent:", h.name, msg.Method)
  fmt.Println(msg.Params)
}

func (h *HTask) OnCdpResponse(msg *Message) bool {
  fmt.Println("==========OnCdpResponse:", h.name, msg.Id, msg.Method)
  fmt.Println(msg.Result)
  return true
}

func TestTask(t *testing.T) {
  chrome, e := Launch("C:/Program Files (x86)/Google/Chrome/Application/chrome.exe")
  if e != nil {
    panic(e)
  }
  taskTB(chrome)
  taskJD(chrome)
  taskAmazon(chrome)
  time.Sleep(time.Second * 10)
  chrome.Exit()
}

func taskTB(chrome *Chrome) {
  h := &HTask{name: "TaoBao"}
  NewTask(chrome).
    Action(NewAction(Page.Enable, nil)).
    Action(NewAction(Page.Navigate, map[string]interface{}{"url": "https://item.taobao.com/item.htm?id=549226118434"})).
    Until(Page.LoadEventFired).
    Action(NewEvalAction("document.querySelector('#J_PromoPriceNum').textContent")).
    Run(h)
}

func taskJD(chrome *Chrome) {
  h := &HTask{name: "JingDong"}
  t := NewTask(chrome).
    Action(NewAction(Page.Enable, nil)).
    Action(NewAction(Page.Navigate, map[string]interface{}{"url": "https://item.jd.com/3693867.html"})).
    Until(Page.LoadEventFired).
    Action(NewEvalAction("document.querySelector('.J-p-3693867').textContent")).
    Run(h)
  time.AfterFunc(time.Second*5, func() {
    t.Finish()
  })
}

func taskAmazon(chrome *Chrome) {
  h := &HTask{name: "Amazon"}
  NewTask(chrome).
    Action(NewAction(Page.Enable, nil)).
    Action(NewAction(Page.Navigate, map[string]interface{}{"url": "https://www.amazon.cn/dp/B072RBZ7T1/"})).
    Until(Page.LoadEventFired).
    Action(NewEvalAction("document.querySelector('.a-color-price').textContent")).
    Run(h)
}
