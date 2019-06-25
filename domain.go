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
  Enable              string
  Disable             string
  DescribeNode        string
  Focus               string
  GetAttributes       string
  GetDocument         string
  QuerySelector       string
  QuerySelectorAll    string
  RemoveAttribute     string
  RemoveNode          string
  RequestChildNodes   string
  RequestNode         string
  ResolveNode         string
  SetAttributeValue   string
  SetAttributesAsText string
  SetNodeName         string
  SetNodeValue        string

  ChildNodeCountUpdated string
  ChildNodeInserted     string
  ChildNodeRemoved      string
  DocumentUpdated       string
}{
  "DOM.enable",
  "DOM.disable",
  "DOM.describeNode",
  "DOM.focus",
  "DOM.getAttributes",
  "DOM.getDocument",
  "DOM.querySelector",
  "DOM.querySelectorAll",
  "DOM.removeAttribute",
  "DOM.removeNode",
  "DOM.requestChildNodes",
  "DOM.requestNode",
  "DOM.resolveNode",
  "DOM.setAttributeValue",
  "DOM.setAttributesAsText",
  "DOM.setNodeName",
  "DOM.setNodeValue",

  "DOM.childNodeCountUpdated",
  "DOM.childNodeInserted",
  "DOM.childNodeRemoved",
  "DOM.documentUpdated",
}

// Chrome里有两类Input事件：可信的与非可信的，
// 可信事件是用户和页面交互触发的事件，如鼠标或键盘，
// 非可信事件是指由Web API触发的，如document.createEvent()或element.click()。
// 例如模拟点击事件，可以使用Runtime.Evaluate，传入表达式document.querySelector('xxx').click()，
// 但这属于非可信事件，如果点击会创建新Tab，浏览器会识别为弹出窗口进行拦截（Page.windowOpen事件仍会触发），
// 所以对于这种情况，只能使用Input.dispatchMouseEvent。
// 注：Input.dispatchMouseEvent模拟点击的话需要连续两次调用（type分别为mousePressed和mouseReleased），
// 并且必须带上button和clickCount参数
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
  "Target.activateTarget",
  "Target.attachToTarget",
  "Target.closeTarget",
  "Target.createTarget",
  "Target.detachFromTarget",
  "Target.getTargets",
  "Target.sendMessageToTarget",
  "Target.setDiscoverTargets",

  "Target.attachedToTarget",
  "Target.detachedFromTarget",
  "Target.receivedMessageFromTarget",
  "Target.targetCreated",
  "Target.targetDestroyed",
  "Target.targetCrashed",
  "Target.targetInfoChanged",
}
