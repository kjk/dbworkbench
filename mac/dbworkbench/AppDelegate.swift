import Cocoa

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate {

    func applicationDidFinishLaunching(aNotification: NSNotification) {
        NSLog("applicationDidFinishLaunching")
    }

    func applicationWillTerminate(aNotification: NSNotification) {
        NSLog("applicationWillTerminate")
    }
    
    func applicationShouldTerminateAfterLastWindowClosed(sender: NSApplication) -> Bool {
        return true
    }
}

