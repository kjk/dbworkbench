import Foundation

var shouldOpenLog = true
var logLock = NSLock()
var logFile : NSOutputStream?

func openLogFileIfNeeded() {
    if logFile != nil {
        return
    }
    if !shouldOpenLog {
        return
    }
    
    let dir =  NSString.pathWithComponents([getDataDir(), "log"])
    
    let dateFmt = NSDateFormatter()
    dateFmt.dateFormat = "'log-'yy-MM-dd'-mac.txt"
    let logName = dateFmt.stringFromDate(NSDate())
    if (!NSFileManager.defaultManager().fileExistsAtPath(dir)) {
        do {
            try NSFileManager.defaultManager().createDirectoryAtPath(dir, withIntermediateDirectories: true, attributes: nil)
        } catch _ {
            // failed
            shouldOpenLog = false
            return
        }
        
    }
    
    let logPath = NSString.pathWithComponents([dir, logName])
    logFile = NSOutputStream(toFileAtPath: logPath, append: true)
    logFile?.open()
}

func closeLogFile() {
    logLock.lock()
    logFile?.close()
    logFile = nil
    logLock.unlock()
}

func log(s : String) {
    logLock.lock()
    openLogFileIfNeeded()
    print(s)
    if let lf = logFile {
        // TODO: only if doesn't end with '\n' already
        let s2 = s + "\n"
        let encodedDataArray = [UInt8](s2.utf8)
        let n = lf.write(encodedDataArray, maxLength: encodedDataArray.count)
        if n == -1 && false {
            print("write failed with error: '\(logFile?.streamError)'")
        }
    }
    logLock.unlock()
}
