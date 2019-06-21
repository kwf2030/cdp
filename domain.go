// Domain的方法/事件都是字符串，
// 为方便使用，这里预置了一些较为常用的，
// 其它的请参考官方文档（https://chromedevtools.github.io/devtools-protocol/tot）

package cdp

var Browser = struct {
  Close      string
  GetVersion string
}{
  "Browser.close",
  "Browser.getVersion",
}

var DOM = struct {
  Enable            string
  Disable           string
  DescribeNode      string
  GetDocument       string
  QuerySelector     string
  QuerySelectorAll  string
  RequestChildNodes string
  RequestNode       string
  ResolveNode       string

  ChildNodeCountUpdated string
  ChildNodeInserted     string
  ChildNodeRemoved      string
  DocumentUpdated       string
}{
  "DOM.enable",
  "DOM.disable",
  "DOM.describeNode",
  "DOM.getDocument",
  "DOM.querySelector",
  "DOM.querySelectorAll",
  "DOM.requestChildNodes",
  "DOM.requestNode",
  "DOM.resolveNode",

  "DOM.childNodeCountUpdated",
  "DOM.childNodeInserted",
  "DOM.childNodeRemoved",
  "DOM.documentUpdated",
}

var Input = struct {
  DispatchKeyEvent   string
  DispatchMouseEvent string
  DispatchTouchEvent string
}{
  "Input.dispatchKeyEvent",
  "Input.dispatchMouseEvent",
  "Input.dispatchTouchEvent",
}

var Network = struct {
  Enable              string
  Disable             string
  ClearBrowserCache   string
  ClearBrowserCookies string
  DeleteCookies       string
  GetAllCookies       string
  GetCookies          string
  GetResponseBody     string
  GetRequestPostData  string

  DataReceived               string
  EventSourceMessageReceived string
  LoadingFailed              string
  LoadingFinished            string
  RequestWillBeSent          string
  ResponseReceived           string
  WebSocketClosed            string
  WebSocketCreated           string
  WebSocketFrameReceived     string
  WebSocketFrameSent         string
}{
  "Network.enable",
  "Network.disable",
  "Network.clearBrowserCache",
  "Network.clearBrowserCookies",
  "Network.deleteCookies",
  "Network.getAllCookies",
  "Network.getCookies",
  "Network.getResponseBody",
  "Network.getRequestPostData",

  "Network.dataReceived",
  "Network.eventSourceMessageReceived",
  "Network.loadingFailed",
  "Network.loadingFinished",
  "Network.requestWillBeSent",
  "Network.responseReceived",
  "Network.webSocketClosed",
  "Network.webSocketCreated",
  "Network.webSocketFrameReceived",
  "Network.webSocketFrameSent",
}

var Page = struct {
  Enable            string
  Disable           string
  BringToFront      string
  CaptureScreenshot string
  Navigate          string
  Reload            string
  StopLoading       string

  DomContentEventFired string
  FrameAttached        string
  FrameDetached        string
  FrameNavigated       string
  LifecycleEvent       string
  LoadEventFired       string
  WindowOpen           string
}{
  "Page.enable",
  "Page.disable",
  "Page.bringToFront",
  "Page.captureScreenshot",
  "Page.navigate",
  "Page.reload",
  "Page.stopLoading",

  "Page.domContentEventFired",
  "Page.frameAttached",
  "Page.frameDetached",
  "Page.frameNavigated",
  "Page.lifecycleEvent",
  "Page.loadEventFired",
  "Page.windowOpen",
}

var Runtime = struct {
  Enable                  string
  Disable                 string
  AwaitPromise            string
  CallFunctionOn          string
  CompileScript           string
  Evaluate                string
  GetProperties           string
  GlobalLexicalScopeNames string
  QueryObjects            string
  ReleaseObject           string
  ReleaseObjectGroup      string
  RunScript               string

  ConsoleAPICalled string
  ExceptionRevoked string
  ExceptionThrown  string
}{
  "Runtime.enable",
  "Runtime.disable",
  "Runtime.awaitPromise",
  "Runtime.callFunctionOn",
  "Runtime.compileScript",
  "Runtime.evaluate",
  "Runtime.getProperties",
  "Runtime.globalLexicalScopeNames",
  "Runtime.queryObjects",
  "Runtime.releaseObject",
  "Runtime.releaseObjectGroup",
  "Runtime.runScript",

  "Runtime.consoleAPICalled",
  "Runtime.exceptionRevoked",
  "Runtime.exceptionThrown",
}

var Target = struct {
  ActivateTarget      string
  AttachToTarget      string
  CloseTarget         string
  CreateTarget        string
  DetachFromTarget    string
  GetTargets          string
  SendMessageToTarget string
  SetDiscoverTargets  string

  AttachedToTarget          string
  DetachedFromTarget        string
  ReceivedMessageFromTarget string
  TargetCreated             string
  TargetDestroyed           string
  TargetCrashed             string
  TargetInfoChanged         string
}{
  "Runtime.activateTarget",
  "Runtime.attachToTarget",
  "Runtime.closeTarget",
  "Runtime.createTarget",
  "Runtime.detachFromTarget",
  "Runtime.getTargets",
  "Runtime.sendMessageToTarget",
  "Runtime.setDiscoverTargets",

  "Runtime.attachedToTarget",
  "Runtime.detachedFromTarget",
  "Runtime.receivedMessageFromTarget",
  "Runtime.targetCreated",
  "Runtime.targetDestroyed",
  "Runtime.targetCrashed",
  "Runtime.targetInfoChanged",
}
