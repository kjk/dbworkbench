import Cocoa
import WebKit

// Maybe: rememver zoom level

// same as in Chrome
let zoomLevels : [Double] = [
    0.25, 0.33, 0.5, 0.67, 0.75, 0.9, 1, 1.1, 1.25, 1.5, 1.75, 2, 2.5, 3, 4, 5
]

// returns an index at which zoomLevel[i] has the smallest difference from z
func findClosestZoomLevelIdx(z : Double) -> Int {
    let min = zoomLevels[0]
    if z <= min {
        return 0
    }
    let n = zoomLevels.count
    let max = zoomLevels[n-1]
    if z >= max {
        return n-1
    }
    var prevDiff = fabs(min - z)
    for i in 1...n {
        let diff = fabs(zoomLevels[i] - z)
        if diff > prevDiff {
            return i - 1
        }
        prevDiff = diff
    }
    return n-1
}

func findNextZoomLevel(z : Double) -> Double {
    var idx = findClosestZoomLevelIdx(z) + 1
    if idx >= zoomLevels.count {
        idx = zoomLevels.count - 1
    }
    return zoomLevels[idx]
}

func findPrevZoomLevel(z : Double) -> Double {
    var idx = findClosestZoomLevelIdx(z) - 1
    if idx < 0 {
        idx = 0
    }
    return zoomLevels[idx]
}

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate {

    @IBOutlet weak var window: NSWindow!
    var webView: WKWebView!
    let backendURL = "http://localhost:5444"

    func autoUpdateCheck() {
        let myVer = NSBundle.mainBundle().shortVersion
        let url = NSURL(string: "https://dbheroapp.com/api/macupdatecheck?ver=" + myVer)!;
        //let url = NSURL(string: "http://localhost:5555/api/macupdatecheck?ver=" + ver); // for testing
        log("url: \(url)")
        let req = NSMutableURLRequest(URL: url)
        let session = NSURLSession.sharedSession()
        req.HTTPMethod = "POST"
        // should roughly match BuildAutoUpdatePostData() in Form1.cs for win
        var s = "ver: \(myVer)\n"
        s += "ostype: mac\n"
        let osVer = getOsVersion()
        let hostName = getHostName()
        let userName = NSUserName()
        let computerId = getUniqueMachineId()
        s += "osversion: \(osVer)\n"
        s += "user: \(userName)\n"
        s += "machine: \(hostName)\n";
        s += "serial: \(computerId)\n"
        s += "---------------\n"; // separator
        if backendUsage != "" {
            s += backendUsage
        }
        req.HTTPBody = s.dataUsingEncoding((NSUTF8StringEncoding))
        let task = session.dataTaskWithRequest(req, completionHandler: {data, response, error -> Void in
            // error is not nil e.g. when the server is not running
            if error != nil {
                log("autoUpdateCheck(): url download failed with error: \(error)")
                return
            }
            guard let httpRsp = response as? NSHTTPURLResponse else {
                log("autoUpdateCheck(): '\(response)' is not NSHTTPURLResponse")
                return
            }
            if httpRsp.statusCode != 200 {
                log("autoUpdateCheck(): response returned status code \(httpRsp.statusCode) which is not 200")
                return
            }
            let dataStr = NSString(data: data!, encoding: NSUTF8StringEncoding)
            let urlVer = parseAutoUpdateCheck(dataStr as! String)
            if programVersionGreater(urlVer.ver!, ver2: myVer) {
                self.notifyAboutUpdate(urlVer.ver!)
            }
        })
        task.resume()
    }
    
    func notifyAboutUpdate(ver : String) {
        let alert = NSAlert()
        alert.messageText = "Update available"
        alert.informativeText = "A new version \(ver) is available. Do you want to update?"
        alert.addButtonWithTitle("Update")
        alert.addButtonWithTitle("No")
        alert.beginSheetModalForWindow(self.window!, completionHandler: {res -> Void in
            if res == NSAlertFirstButtonReturn {
                // TODO: a more specific page with just a download button to
                // download new version and instructions on how to update
                self.goToWebsite()
            }
        })
    }
    
    func showBackendFailedError() {
        let alert = NSAlert()
        alert.messageText = "Backend failed"
        alert.informativeText = "Failed to start"
        alert.addButtonWithTitle("Exit the app")
        // TODO: maybe a way to contact support
        alert.beginSheetModalForWindow(self.window!, completionHandler: {res -> Void in
            log("shutting down the app")
            NSApp.terminate(nil)
        })
    }

    func urlReq(url: String) -> NSURLRequest {
        let u = NSURL(string: url)
        return NSURLRequest(URL: u!)
    }

    
    func viewMidPoint(view : NSView) -> CGPoint {
        let b = view.bounds
        return CGPoint(x: b.width / 2, y: b.height / 2)
    }

    @IBAction func zoomIn(sender: AnyObject) {
        let currZoom = Double(webView.magnification)
        let newZoom = findNextZoomLevel(currZoom)
        webView.setMagnification(CGFloat(newZoom), centeredAtPoint: viewMidPoint(webView))
        //log("zoomIn: from \(currZoom) to \(newZoom)")
    }

    @IBAction func zoomOut(sender: AnyObject) {
        let currZoom = Double(webView.magnification)
        let newZoom = findPrevZoomLevel(currZoom)
        webView.setMagnification(CGFloat(newZoom), centeredAtPoint: viewMidPoint(webView))
        //log("zoomIn: from \(currZoom) to \(newZoom)")
    }

    @IBAction func actualSize(sender: AnyObject) {
        webView.setMagnification(1.0, centeredAtPoint: viewMidPoint(webView))
    }

    func loadURL() {
        log("loadURL")
        webView.loadRequest(urlReq(backendURL))
    }

    func applicationDidFinishLaunching(aNotification: NSNotification) {
        log("applicationDidFinishLaunching")
        webView = WKWebView(frame: NSRect(x: 0, y: 0, width: 400, height: 400))
        webView.allowsMagnification = true
        webView.translatesAutoresizingMaskIntoConstraints = false
        let cv = window.contentView!
        cv.subviews.append(webView)

        let v = window.contentView!

        let cWidth = NSLayoutConstraint(item: webView,
            attribute: NSLayoutAttribute.Width,
            relatedBy: .Equal,
            toItem: v,
            attribute: .Width,
            multiplier: 1,
            constant: 0)

        let cHeight = NSLayoutConstraint(item: webView,
            attribute: NSLayoutAttribute.Height,
            relatedBy: .Equal,
            toItem: v,
            attribute: .Height,
            multiplier: 1,
            constant: 0)

        let constraints = [
            cWidth,
            cHeight,
        ]

        // Sepculation: contentView has constraints that bind its size to size
        // of the window. We need to lower priority of our constraints below
        // that of contentView. Otherwise ultimately we bind the size of the
        // window to current size of webView, making it un-resizeable
        for c in constraints {
            c.priority = NSLayoutPriorityDragThatCannotResizeWindow
        }

        v.addConstraints(constraints)

        //webView.mainFrame.loadRequest(urlReq("https://blog.kowalczyk.info"))

        loadUsageData()
        startBackend(self)
        autoUpdateCheck()
    }

    func applicationShouldTerminateAfterLastWindowClosed(sender: NSApplication) -> Bool {
        return true
    }

    func applicationWillTerminate(aNotification: NSNotification) {
        stopBackend()
        closeLogFile()
    }

    func goToWebsite() {
        NSWorkspace.sharedWorkspace().openURL(NSURL(string: "https://dbheroapp.com")!)
    }
    
    @IBAction func goToWebsite(sender: NSMenuItem) {
        goToWebsite()
    }
    
    @IBAction func goToSupport(sender: NSMenuItem) {
        NSWorkspace.sharedWorkspace().openURL(NSURL(string: "https://dbheroapp.com/support")!)
    }
    
    @IBAction func goToFeedback(sender: NSMenuItem) {
        NSWorkspace.sharedWorkspace().openURL(NSURL(string: "https://dbheroapp.com/feedback")!)
    }
}

