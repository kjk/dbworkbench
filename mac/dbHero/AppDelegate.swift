import Cocoa

extension NSBundle {
    
    var shortVersion: String {
        if let ver = self.infoDictionary?["CFBundleShortVersionString"] as? String {
            return ver
        }
        return "unknown"
    }
    
    var version: String {
        if let ver = self.infoDictionary?["CFBundleVersion"] as? String {
            return ver
        }
        return "unknown"
    }
    
}

func parseAutoUpdateCheck(s : String) -> (ver: String?, url: String?) {
    var ver : String?
    var url : String?
    let parts = s.componentsSeparatedByString("\n")
    for p in parts {
        let parts = p.componentsSeparatedByString(": ")
        if parts.count != 2 {
            continue
        }
        let name = parts[0]
        let val = parts[1]
        if name == "ver" {
            ver = val
        } else if name == "url" {
            url = val
        }
    }
    return (ver, url)
}

func versionArrayElToInt(a : [String], pos : Int) -> Int {
    if pos >= a.count {
        return 0
    }

    let s = a[pos]
    if let n = Int(s) {
        return n
    }
    log("versionArrayElToInt: '\(s)' is not a valid number")
    return 0
}

// return true if ver1 > ver2
// version is in the format "0.1.3", "1.2" etc.
func programVersionGreater(ver1 : String, ver2 : String) -> Bool {
    let parts1 = ver1.componentsSeparatedByString(".")
    let parts2 = ver2.componentsSeparatedByString(".")
    var n = parts1.count
    if parts2.count > n {
        n = parts2.count
    }
    for (var i = 0; i < n; i++) {
        let n1 = versionArrayElToInt(parts1, pos: i)
        let n2 = versionArrayElToInt(parts2, pos: i)
        if n1 > n2 {
            return true
        }
        if n1 < n2 {
            return false
        }
        // are equal so check the next element
    }
    return false // are equal so not greater
}

// http://stackoverflow.com/questions/5868567/unique-identifier-of-a-mac
func getMacSerialNumber() -> String {
    let platformExpert: io_service_t = IOServiceGetMatchingService(kIOMasterPortDefault, IOServiceMatching("IOPlatformExpertDevice"));
    let serialNumberAsCFString = IORegistryEntryCreateCFProperty(platformExpert, kIOPlatformSerialNumberKey, kCFAllocatorDefault, 0);
    IOObjectRelease(platformExpert);
    // Take the unretained value of the unmanaged-any-object
    // (so we're not responsible for releasing it)
    // and pass it back as a String or, if it fails, an empty string
    return (serialNumberAsCFString.takeUnretainedValue() as? String) ?? ""
}

/*
// if we decide to support 10.9, we'll need this version
// return os version in the "10.11.1" form
func getOsVersion2() -> String {
    if #available(OSX 10.10, *) {
        let os = NSProcessInfo().operatingSystemVersion
        return "\(os.majorVersion).\(os.minorVersion).\(os.patchVersion)"
    } else {
        let ver = rint(NSAppKitVersionNumber)
        if ver >= Double(NSAppKitVersionNumber10_10_Max) {
            return "10.10.5+"
        }
        if ver >= Double(NSAppKitVersionNumber10_10_5) {
            return "10.10.5"
        }
        if ver >= Double(NSAppKitVersionNumber10_10_4) {
            return "10.10.4"
        }
        if ver >= Double(NSAppKitVersionNumber10_10_3) {
            return "10.10.3"
        }
        if ver >= Double(NSAppKitVersionNumber10_10_2) {
            return "10.10.2"
        }
        if ver >= Double(NSAppKitVersionNumber10_10) {
            return "10.10"
        }
        if ver >= Double(NSAppKitVersionNumber10_9) {
            return "10.9"
        }
        if ver >= Double(NSAppKitVersionNumber10_8) {
            return "10.8"
        }
        return "unknown: \(ver)"
    }
}
*/

func getOsVersion() -> String {
    let os = NSProcessInfo().operatingSystemVersion
    return "\(os.majorVersion).\(os.minorVersion).\(os.patchVersion)"
}

func getHostName() -> String {
    return NSProcessInfo.processInfo().hostName
}

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate {

    weak var window : NSWindow?

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
        let computerId = getMacSerialNumber()
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

    func applicationDidFinishLaunching(aNotification: NSNotification) {
        log("applicationDidFinishLaunching")
        autoUpdateCheck()
    }

    func applicationWillTerminate(aNotification: NSNotification) {
        closeLogFile()
    }
    
    func applicationShouldTerminateAfterLastWindowClosed(sender: NSApplication) -> Bool {
        return true
    }

    func goToWebsite() {
        NSWorkspace.sharedWorkspace().openURL(NSURL(string: "https://dbheroapp.com/for-mac")!)
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

func getAppDelegate() -> AppDelegate {
    return NSApp.delegate as! AppDelegate
}
