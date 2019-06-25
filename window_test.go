package cdp

import (
  "fmt"
  "sync"
  "testing"
  "time"
)

var wg1 sync.WaitGroup

type H1 struct{}

func (h1 *H1) OnCdpEvent(msg *Message) {
  fmt.Println("======OnCdpEvent:", msg.Method)
  fmt.Println(msg.Params)
}

func (h1 *H1) OnCdpResponse(msg *Message) bool {
  fmt.Println("======OnCdpResponse:", msg.Method, msg.Result)
  return false
}

func TestWindow(t *testing.T) {
  // chrome, e := Launch("C:/Program Files (x86)/Google/Chrome/Application/chrome.exe")
  chrome, e := Launch("C:/App/Chromium/chrome.exe")
  if e != nil {
    panic(e)
  }
  wg1.Add(1)
  h1 := &H1{}
  tab, e := chrome.NewTab(h1)
  if e != nil {
    panic(e)
  }
  tab.Subscribe(Page.LoadEventFired, Page.WindowOpen, Target.AttachedToTarget,
    Target.DetachedFromTarget, Target.ReceivedMessageFromTarget, Target.TargetCreated,
    Target.TargetDestroyed, Target.TargetCrashed, Target.TargetInfoChanged)
  tab.Call(Page.Enable, nil)
  tab.Call(Page.Navigate, map[string]interface{}{"url": "https://shanghai.anjuke.com/community/?from=navigation"})
  time.Sleep(time.Second * 5)
  tab.Call(Input.DispatchMouseEvent, map[string]interface{}{"type": "mousePressed", "x": 200, "y": 600, "button": "left", "clickCount": 1})
  tab.Call(Input.DispatchMouseEvent, map[string]interface{}{"type": "mouseReleased", "x": 200, "y": 600, "button": "left", "clickCount": 1})
  _, ch := tab.Call(Target.GetTargets, nil)
  msg := <-ch
  fmt.Println(msg)
  wg1.Wait()
  chrome.Exit()
}
