#!/usr/bin/env xcrun swift

import Foundation

let task = NSTask()
task.launchPath = "/bin/ps"
task.arguments = ["-ax"]

let pipe = NSPipe()
task.standardOutput = pipe
task.launch()

let data = pipe.fileHandleForReading.readDataToEndOfFile()
if let output = NSString(data: data, encoding: NSUTF8StringEncoding) {
	print(output)
} else {
	print("failed to parse the output")
}

