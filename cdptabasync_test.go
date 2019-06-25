package cdp

import (
  "fmt"
  "sync"
  "testing"
)

var wg sync.WaitGroup

type H struct {
  name   string
  expr   string
  tab    *Tab
  callId int32
}

func (h *H) OnCdpEvent(msg *Message) {
  fmt.Println("======OnCdpEvent:", h.name, msg.Method)
  fmt.Println(msg.Params)
  if msg.Method == Page.LoadEventFired {
    id, _ := h.tab.Call(Runtime.Evaluate, map[string]interface{}{"returnByValue": true, "expression": h.expr})
    fmt.Println("id:", id)
    h.callId = id
  }
}

func (h *H) OnCdpResponse(msg *Message) bool {
  fmt.Println("======OnCdpResponse:", h.name, msg.Id, msg.Method)
  fmt.Println(msg.Result)
  if msg.Id == h.callId {
    h.tab.Close()
    wg.Done()
  }
  return true
}

func TestTabAsync(t *testing.T) {
  chrome, e := Launch("C:/Program Files (x86)/Google/Chrome/Application/chrome.exe")
  if e != nil {
    panic(e)
  }
  wg.Add(3)
  tabTB(chrome)
  tabJD(chrome)
  tabAmazon(chrome)
  wg.Wait()
  chrome.Exit()
}

func tabTB(chrome *Chrome) {
  h := &H{name: "TaoBao", expr: "document.querySelector('#J_PromoPriceNum').textContent"}
  tab, e := chrome.NewTab(h)
  if e != nil {
    panic(e)
  }
  h.tab = tab
  tab.Subscribe(Page.LoadEventFired)
  tab.Call(Page.Enable, nil)
  tab.Call(Page.Navigate, map[string]interface{}{"url": "https://item.taobao.com/item.htm?id=549226118434"})
}

func tabJD(chrome *Chrome) {
  h := &H{name: "JingDong", expr: "document.querySelector('.J-p-3693867').textContent"}
  tab, e := chrome.NewTab(h)
  if e != nil {
    panic(e)
  }
  h.tab = tab
  tab.Subscribe(Page.LoadEventFired)
  tab.Call(Page.Enable, nil)
  tab.Call(Page.Navigate, map[string]interface{}{"url": "https://item.jd.com/3693867.html"})
}

func tabAmazon(chrome *Chrome) {
  h := &H{name: "Amazon", expr: "document.querySelector('.a-color-price').textContent"}
  tab, e := chrome.NewTab(h)
  if e != nil {
    panic(e)
  }
  h.tab = tab
  tab.Subscribe(Page.LoadEventFired)
  tab.Call(Page.Enable, nil)
  tab.Call(Page.Navigate, map[string]interface{}{"url": "https://www.amazon.cn/dp/B072RBZ7T1/"})
}
