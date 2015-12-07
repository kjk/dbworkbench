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

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate {

    
    func autoUpdateCheck() {
        let ver = NSBundle.mainBundle().shortVersion

        //Then just cast the object as a String, but be careful, you may want to double check for nil
        //let url = NSURL(string: "http://databaseworkbench.com/api/macupdatecheck?ver=" + ver);
        let url = NSURL(string: "http://localhost:5555/api/macupdatecheck?ver=" + ver);
        NSLog("url: \(url)")
        let req = NSMutableURLRequest(URL: url!)
        let session = NSURLSession.sharedSession()
        req.HTTPMethod = "POST"
        let s = "ver: \(ver)\n"
        req.HTTPBody = s.dataUsingEncoding((NSUTF8StringEncoding))
        let task = session.dataTaskWithRequest(req, completionHandler: {data, response, error -> Void in
            // error is not nil e.g. when the server is not running
            if error != nil {
                NSLog("autoUpdateCheck(): url download failed with error: \(error)")
                return
            }
            guard let httpRsp = response as? NSHTTPURLResponse else {
                NSLog("autoUpdateCheck(): '\(response)' is not NSHTTPURLResponse")
                return
            }
            if httpRsp.statusCode != 200 {
                NSLog("autoUpdateCheck(): response returned status code \(httpRsp.statusCode) which is not 200")
                return
            }
            let dataStr = NSString(data: data!, encoding: NSUTF8StringEncoding)
            NSLog("got post response \(dataStr)")
            let urlVer = parseAutoUpdateCheck(dataStr as! String)
            if urlVer.ver != ver {
                // TODO: notify the user that update is available
                NSLog("should autoUpdate because versions different! my ver: \(ver), latest ver: \(urlVer.ver)")
            }
            
            
        })
        task.resume()
    }

    func applicationDidFinishLaunching(aNotification: NSNotification) {
        NSLog("applicationDidFinishLaunching")
        autoUpdateCheck()
    }

    func applicationWillTerminate(aNotification: NSNotification) {
        NSLog("applicationWillTerminate")
    }
    
    func applicationShouldTerminateAfterLastWindowClosed(sender: NSApplication) -> Bool {
        return true
    }
}

