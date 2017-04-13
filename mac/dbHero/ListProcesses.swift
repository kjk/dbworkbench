import Foundation

struct ProcInfo {
    var pathAndArgs: String
    var pid: Int
}

// parses output of ps -ax line, which is in the format:
// "  454 ttys000    0:00.05 -bash"
func parsePsLine(_ s : String) -> ProcInfo? {
    var parts = s.characters.split(separator: " ", maxSplits: 3, omittingEmptySubsequences: true).map(String.init)
    if parts.count != 4 {
        log("parsePsLine: failed to parse: '\(s)'")
        return nil
    }
    
    let pidStr = parts[0]
    guard let pid = Int(pidStr) else {
        log("parsePsLine: failed to convert pid '\(pidStr)' to int")
        return nil
    }
    return ProcInfo(pathAndArgs: parts[3], pid: pid)
}

func listProcesses() -> [ProcInfo] {
    let task = Process()
    task.launchPath = "/bin/ps"
    task.arguments = ["-ax"]
    
    let pipe = Pipe()
    task.standardOutput = pipe
    task.launch()
    
    var res = [ProcInfo]()
    let data = pipe.fileHandleForReading.readDataToEndOfFile()
    guard let output = NSString(data: data, encoding: String.Encoding.utf8.rawValue) as? String else {
        log("listProcesses: failed to parse the output")
        return res
    }
    var lines = output.characters.split(separator: "\n", omittingEmptySubsequences: true).map(String.init)
    // first line is the header
    if lines.count > 0 {
        lines.remove(at: 0)
    }
    for l in lines {
        if let pi = parsePsLine(l) {
            res.append(pi)
        }
    }
    return res
}
