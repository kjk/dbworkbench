import Cocoa
import AppKit

var backendTask = NSTask()
var waitsForMoreServerOutput = true

// kill other instances of mac app. this is needed for the case when a user
// downloads an update and runs it without quitting the previous version
func killOtherAppInstances() {
    let wsk = NSWorkspace.sharedWorkspace()
    let processes = wsk.runningApplications
    guard let myPath = NSRunningApplication.currentApplication().executableURL?.path else {
        return
    }
    for proc in processes {
        guard let appUrl = proc.executableURL else {
            continue
        }
        guard let path = appUrl.path else {
            continue
        }
        guard path.containsString("dbHero.app") else {
            continue
        }
        //log("path: \(path)")
        guard path != myPath else {
            continue
        }
        proc.forceTerminate()
        let pid = proc.processIdentifier
        log("killing process \(pid) '\(path)'")
        proc.forceTerminate()
        // wait up to 3 secs for process to terminate
        // TODO: proc.terminated seems to be false even if process is gone
        var i = 3
        while i > 0 {
            if proc.terminated {
                return
            }
            sleep(1)
            i -= 1
        }
    }
}

// kill all dbherohelper.exe processes because if they are running, we'll fail
// to start. They might exist in case the mac app didn't exit cleanly and didn't
// kill its backend. Or maybe the user launched another instance from a different
// path (e.g. downloaded an upgrade and is still running the previous version)
func killBackendIfRunning() {
    let procs = listProcesses()
    for p in procs {
        if p.pathAndArgs.containsString("/dbherohelper.exe") {
            log("killing process \(p.pid) '\(p.pathAndArgs)'")
            kill(pid_t(p.pid), SIGKILL) // SIGINT ?
        }
    }
}

func getDataDir() -> String {
    return NSString.pathWithComponents([NSHomeDirectory(), "Library", "Application Support", "dbHero"])
}

var backendUsage = ""

// must be executed before starting backend in order to read usage.json
func loadUsageData() {
    let path = NSString.pathWithComponents([getDataDir(), "usage.json"])
    do {
        let s = try NSString(contentsOfFile: path, encoding: NSUTF8StringEncoding)
        backendUsage = s as String;
        // delete so that we don't send duplicate data
        try NSFileManager.defaultManager().removeItemAtPath(path)
    }
    catch let error as NSError {
        log("loadUsageData: error: \(error)")
    }
}

// Maybe: instead of passing appDelegate, use notifcations or getDelegate()
func startBackend(appDelegate : AppDelegate) {
    let resPath = NSBundle.mainBundle().resourcePath
    let backendGoExePath = resPath! + "/dbherohelper.exe"

    killOtherAppInstances()
    killBackendIfRunning()

    // NSTask.launch() will throw objc exception if file doesn't exist
    // it'll be ultimately caught by Cocoa if we're running on ui thread
    // but it'll not execute any of our code after that. Ideally we would
    // handle the exception but that's not possible in swift and I'm too lazy
    // to write this part in Objective-C
    // So we guard against the most common reason for this error with up-front check
    // I also tried running on background thread, but that would crash the whole app
    // due to uncought exception
    let exists = NSFileManager.defaultManager().fileExistsAtPath(backendGoExePath)
    if !exists {
        appDelegate.showBackendFailedError()
        return
    }
    
    backendTask.terminationHandler = { task -> Void in
        log("backendTask terminated");
        // this runs on non-main thread so marshal on main thread
        dispatch_async(dispatch_get_main_queue(),{
            appDelegate.showBackendFailedError()
        })
    }

    backendTask.launchPath = backendGoExePath
    backendTask.currentDirectoryPath = resPath!
    //        serverTask.arguments = ["-dev"]

    let pipe = NSPipe()
    backendTask.standardOutput = pipe
    backendTask.standardError = pipe

    let outHandle = pipe.fileHandleForReading
    outHandle.waitForDataInBackgroundAndNotify()

    let _ = NSNotificationCenter.defaultCenter().addObserverForName(NSFileHandleDataAvailableNotification, object: outHandle, queue: nil, usingBlock: { notification -> Void in

        if !waitsForMoreServerOutput {
            return
        }

        let output = outHandle.availableData
        let outStr = NSString(data: output, encoding: NSUTF8StringEncoding)
        // wait until backend prints "Started running on..."
        if outStr?.length > 0 {
            let s = outStr! as String
            if (s.containsString("failed with")) {
                // TODO: notify about the error in the UI
                // this could be "http.ListendAndServer() failed with listen tcp 127.0.0.1:5444: bind: address already in use"
                log("startBackend: failed because output is: \(s)")
                waitsForMoreServerOutput = false
                return
            }
            if (s.containsString("Started running on")) {
                log("startBackend: backend started, loading url")
                waitsForMoreServerOutput = false
                appDelegate.loadURL()
                return
            }
        }
        outHandle.waitForDataInBackgroundAndNotify()
    })

    backendTask.launch()
}

func stopBackend() {
    log("stopping backend")
    if backendTask.running {
        backendTask.terminationHandler = nil
        backendTask.terminate()
    }
}

