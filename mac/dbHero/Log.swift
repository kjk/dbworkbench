import Foundation

var shouldOpenLog = true
var logLock = NSLock()
var logFile : OutputStream?

func openLogFileIfNeeded() {
    if logFile != nil {
        return
    }
    if !shouldOpenLog {
        return
    }

    let dir =  NSString.path(withComponents: [getDataDir(), "log"])

    let dateFmt = DateFormatter()
    dateFmt.dateFormat = "'log-'yy-MM-dd'-mac.txt"
    let logName = dateFmt.string(from: Date())
    if (!FileManager.default.fileExists(atPath: dir)) {
        do {
            try FileManager.default.createDirectory(atPath: dir, withIntermediateDirectories: true, attributes: nil)
        } catch _ {
            // failed
            shouldOpenLog = false
            return
        }

    }

    let logPath = NSString.path(withComponents: [dir, logName])
    logFile = OutputStream(toFileAtPath: logPath, append: true)
    logFile?.open()
}

func closeLogFile() {
    logLock.lock()
    logFile?.close()
    logFile = nil
    logLock.unlock()
}

func log(_ s : String) {
    logLock.lock()
    openLogFileIfNeeded()
    print(s)
    if let lf = logFile {
        var s2 = s
        if !s.hasSuffix("\n") {
            s2 = s + "\n"
        }
        let encodedDataArray = [UInt8](s2.utf8)
        let n = lf.write(encodedDataArray, maxLength: encodedDataArray.count)
        if n == -1 && false {
            print("write failed with error: '\(String(describing: logFile?.streamError))'")
        }
    }
    logLock.unlock()
}
