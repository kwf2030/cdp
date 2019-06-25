package main

import (
  "bufio"
  "encoding/json"
  "fmt"
  "os"
  "time"

  "github.com/kwf2030/cdp"
)

type REPL struct{}

func (r *REPL) OnCdpEvent(msg *cdp.Message) {
  fmt.Println("======OnCdpEvent:", msg.Method)
  fmt.Println(msg.Params)
}

func (r *REPL) OnCdpResponse(msg *cdp.Message) bool {
  fmt.Println("======OnCdpResponse:", msg.Id, msg.Method)
  fmt.Println(msg.Result)
  return true
}

func main() {
  chrome, e := cdp.Launch("C:/App/Chromium/chrome.exe")
  // chrome, e := cdp.Connect("127.0.0.1", 9222)
  if e != nil {
    panic(e)
  }

  tab, e := chrome.NewTab(&REPL{})
  if e != nil {
    panic(e)
  }
  tab.Subscribe(cdp.Page.LoadEventFired, cdp.Page.WindowOpen,
    cdp.Target.AttachedToTarget, cdp.Target.DetachedFromTarget, cdp.Target.ReceivedMessageFromTarget,
    cdp.Target.TargetCreated, cdp.Target.TargetDestroyed, cdp.Target.TargetCrashed,
    cdp.Target.TargetInfoChanged)

  scanner := bufio.NewScanner(os.Stdin)
  for {
    fmt.Print("Method: ")
    scanner.Scan()
    method := scanner.Text()
    if method == "exit" {
      chrome.Exit()
      return
    }

    fmt.Print("Params: ")
    scanner.Scan()
    str := scanner.Text()
    var params map[string]interface{}
    if str != "" {
      e := json.Unmarshal([]byte(str), &params)
      if e != nil {
        fmt.Println(e)
        continue
      }
    }

    id, _ := tab.Call(method, params)
    fmt.Println("id: ", id)
    time.Sleep(time.Second)
  }
}
