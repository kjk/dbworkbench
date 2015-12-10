#!/usr/bin/env xcrun swift

import Foundation

struct ProcInfo {
	var name: String
	var pid: Int
}

// parses output of ps -ax line, which is in the format:
// "  454 ttys000    0:00.05 -bash"
func parsePsLine(s : String) -> ProcInfo? {
    var parts = s.characters.split(" ", maxSplit: 3, allowEmptySlices: false).map(String.init)
	if parts.count != 4 {
		print("failed to parse: '\(s)'")
		return nil
	}

	let pidStr = parts[0]
	guard let pid = Int(pidStr) else {
		print("failed to convert '\(pidStr)' to int")
		return nil
	}
	let name = parts[3]
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

let procs = listProcesses()
for p in procs {
	print("pid: \(p.pid) name: '\(p.name)'")
}