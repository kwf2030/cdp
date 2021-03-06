package cdp

import (
  "fmt"
  "testing"
  "time"
)

type H1 struct{}

func (h1 *H1) OnCdpEvent(msg *Message) {
  fmt.Println("==========OnCdpEvent:", msg.Method)
  fmt.Println(msg.Params)
}

func (h1 *H1) OnCdpResponse(msg *Message) bool {
  fmt.Println("==========OnCdpResponse:", msg.Id, msg.Method)
  fmt.Println(msg.Result)
  return false
}

func TestWindow(t *testing.T) {
  chrome, e := Launch("C:/Program Files (x86)/Google/Chrome/Application/chrome.exe")
  if e != nil {
    panic(e)
  }
  h1 := &H1{}
  tab, e := chrome.NewTab(h1)
  if e != nil {
    panic(e)
  }
  tab.Subscribe(Page.LoadEventFired, Page.WindowOpen)
  tab.Call(Page.Enable, nil)
  tab.Call(Page.Navigate, map[string]interface{}{"url": "https://shanghai.anjuke.com/community/?from=navigation"})
  time.Sleep(time.Second * 5)
  tab.Call(Input.DispatchMouseEvent, map[string]interface{}{"type": "mousePressed", "x": 200, "y": 600, "button": "left", "clickCount": 1})
  tab.Call(Input.DispatchMouseEvent, map[string]interface{}{"type": "mouseReleased", "x": 200, "y": 600, "button": "left", "clickCount": 1})
  time.Sleep(time.Second * 2)
  _, ch := tab.Call(Target.GetTargets, nil)
  msg := <-ch
  fmt.Println(msg)
  chrome.Exit()
}
