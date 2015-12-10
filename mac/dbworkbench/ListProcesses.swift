import Foundation

struct ProcInfo {
    var name: String
    var pid: Int
}

func trimEmptyStringsLeft(var a : [String]) -> [String] {
    while !a.isEmpty {
        if a[0] == "" {
            a.removeAtIndex(0)
        } else {
            return a
        }
    }
    return a
}

func removeEmptyStrings(var a : [String]) -> [String] {
    var res = [String]()
    for var i = 0; i < a.count; i++ {
        let el = a[i]
        if el != "" {
            res.append(el)
        }
    }
    return res
}

// parses output of ps -ax line, which is in the format:
// "  454 ttys000    0:00.05 -bash"
func parsePsLine(s : String) -> ProcInfo? {
    var parts = s.componentsSeparatedByString(" ")
    parts = trimEmptyStringsLeft(parts)
    
    //print(s)
    // this should be pid
    if parts.isEmpty {
        return nil
    }
    let pidStr = parts[0]
    guard let pid = Int(pidStr) else {
        print("failed to convert '\(pidStr)' to int")
        return nil
    }
    parts.removeAtIndex(0)
    
    parts = trimEmptyStringsLeft(parts)
    if parts.isEmpty {
        return nil
    }
    
    //let user = parts[0]
    parts.removeAtIndex(0)
    
    parts = trimEmptyStringsLeft(parts)
    if parts.isEmpty {
        return nil
    }
    
    //let time = parts[0]
    parts.removeAtIndex(0)
    
    parts = trimEmptyStringsLeft(parts)
    if parts.isEmpty {
        return nil
    }
    
    // the rest is process name and arguments
    parts = removeEmptyStrings(parts)
    let name = parts.joinWithSeparator(" ")
    return ProcInfo(name: name, pid: pid)
}

func listProcesses() -> [ProcInfo] {
    let task = NSTask()
    task.launchPath = "/bin/ps"
    task.arguments = ["-ax"]
    
    let pipe = NSPipe()
    task.standardOutput = pipe
    task.launch()
    
    var res = [ProcInfo]()
    let data = pipe.fileHandleForReading.readDataToEndOfFile()
    guard let output = NSString(data: data, encoding: NSUTF8StringEncoding) else {
        print("failed to parse the output")
        return res
    }
    var lines = output.componentsSeparatedByString("\n")
    // first line is the header
    if lines.count > 0 {
        lines.removeAtIndex(0)
    }
    for l in lines {
        if let pi = parsePsLine(l) {
            res.append(pi)
        }
    }
    return res
}
