package cdp

import (
  "fmt"
  "sync"
  "testing"
)

var wgSync sync.WaitGroup

type HSync struct {
  name string
  expr string
  ch   chan struct{}
  tab  *Tab
}

func (h *HSync) OnCdpEvent(msg *Message) {
  fmt.Println("==========OnCdpEvent:", h.name, msg.Method)
  fmt.Println(msg.Params)
  if msg.Method == Page.LoadEventFired {
    _, ch := h.tab.Call(Runtime.Evaluate, map[string]interface{}{"returnByValue": true, "expression": h.expr})
    resp := <-ch
    fmt.Println("expr result:", h.name, resp.Method, resp.Id, resp.Result)
    h.ch <- struct{}{}
    h.tab.Close()
    wgSync.Done()
  }
}

func (h *HSync) OnCdpResponse(msg *Message) bool {
  fmt.Println("==========OnCdpResponse:", h.name, msg.Id, msg.Method)
  fmt.Println(msg.Result)
  return false
}

func TestTabSync(t *testing.T) {
  chrome, e := Launch("C:/Program Files (x86)/Google/Chrome/Application/chrome.exe")
  if e != nil {
    panic(e)
  }
  wgSync.Add(3)
  // 如果cap为零最后一次会阻塞
  ch := make(chan struct{}, 1)
  go func() {
    fs := []func(*Chrome, chan struct{}){tabSyncTB, tabSyncJD, tabSyncAmazon}
    for _, f := range fs {
      <-ch
      f(chrome, ch)
    }
  }()
  ch <- struct{}{}
  wgSync.Wait()
  chrome.Exit()
}

func tabSyncTB(chrome *Chrome, ch chan struct{}) {
  h := &HSync{name: "TaoBao", expr: "document.querySelector('#J_PromoPriceNum').textContent", ch: ch}
  tab, e := chrome.NewTab(h)
  if e != nil {
    panic(e)
  }
  h.tab = tab
  tab.Subscribe(Page.LoadEventFired)
  tab.Call(Page.Enable, nil)
  tab.Call(Page.Navigate, map[string]interface{}{"url": "https://item.taobao.com/item.htm?id=549226118434"})
}

func tabSyncJD(chrome *Chrome, ch chan struct{}) {
  h := &HSync{name: "JingDong", expr: "document.querySelector('.J-p-3693867').textContent", ch: ch}
  tab, e := chrome.NewTab(h)
  if e != nil {
    panic(e)
  }
  h.tab = tab
  tab.Subscribe(Page.LoadEventFired)
  tab.Call(Page.Enable, nil)
  tab.Call(Page.Navigate, map[string]interface{}{"url": "https://item.jd.com/3693867.html"})
}

func tabSyncAmazon(chrome *Chrome, ch chan struct{}) {
  h := &HSync{name: "Amazon", expr: "document.querySelector('.a-color-price').textContent", ch: ch}
  tab, e := chrome.NewTab(h)
  if e != nil {
    panic(e)
  }
  h.tab = tab
  tab.Subscribe(Page.LoadEventFired)
  tab.Call(Page.Enable, nil)
  tab.Call(Page.Navigate, map[string]interface{}{"url": "https://www.amazon.cn/dp/B072RBZ7T1/"})
}
