package main

import (
  "bufio"
  "encoding/json"
  "fmt"
  "os"
  "strings"
  "sync"
  "time"

  "github.com/kwf2030/cdp"
)

var (
  wg     sync.WaitGroup
  chrome *cdp.Chrome
  tab    *cdp.Tab
)

type REPL struct{}

func (r *REPL) OnCdpEvent(msg *cdp.Message) {
  fmt.Printf("\n==========OnCdpEvent: {method: %s}\n", msg.Method)
  fmt.Println(msg.Params)
}

func (r *REPL) OnCdpResponse(msg *cdp.Message) bool {
  fmt.Printf("==========OnCdpResponse: {id: %d, method: %s}\n", msg.Id, msg.Method)
  fmt.Println(msg.Result)
  go readStdin()
  return true
}

func readStdin() {
  scanner := bufio.NewScanner(os.Stdin)
  fmt.Printf("Method: ")
  scanner.Scan()
  str := scanner.Text()
  if str == "exit" {
    chrome.Exit()
    time.Sleep(time.Millisecond * 500)
    wg.Done()
    return
  }
  methods := strings.Split(str, ",")
  params := make([]map[string]interface{}, len(methods))
  for i := range methods {
    methods[i] = strings.TrimSpace(methods[i])
    if methods[i] == "" {
      continue
    }

    fmt.Printf("Params(%d): ", i)
    scanner.Scan()
    str = scanner.Text()
    if str != "" {
      var m map[string]interface{}
      e := json.Unmarshal([]byte(str), &m)
      if e != nil {
        fmt.Println(e)
        continue
      }
      params[i] = m
    }
  }

  for i, m := range methods {
    if m != "" {
      tab.Call(m, params[i])
    }
  }
}

func main() {
  wg.Add(1)

  var e error
  chrome, e = cdp.Launch("C:/Program Files (x86)/Google/Chrome/Application/chrome.exe")
  if e != nil {
    panic(e)
  }

  tab, e = chrome.NewTab(&REPL{})
  if e != nil {
    panic(e)
  }
  tab.Subscribe(cdp.Page.LoadEventFired, cdp.Page.WindowOpen,
    cdp.Target.AttachedToTarget, cdp.Target.DetachedFromTarget, cdp.Target.ReceivedMessageFromTarget,
    cdp.Target.TargetCreated, cdp.Target.TargetDestroyed, cdp.Target.TargetCrashed,
    cdp.Target.TargetInfoChanged)

  go readStdin()

  wg.Wait()
}
