//
//  ServerController.swift
//  dbworkbench
//
//  Created by Furkan Yilmaz on 28/11/15.
//  Copyright © 2015 Furkan Yilmaz. All rights reserved.
//

import Foundation

class ServerController {
    
    static var serverTask =  NSTask()
    
    static func runServer() {
        // TODO: launch server executable
        
//        serverTask.launchPath = "/usr/sbin/sysctl"
//        serverTask.arguments = ["-w", "net.inet.tcp.delayed_ack=0"]
//        
//        let pipe = NSPipe()
//        serverTask.standardOutput = pipe
//        serverTask.standardError = pipe
        
//        serverTask.launch()
                
//        let data = pipe.fileHandleForReading.readDataToEndOfFile()
//        let output: String = NSString(data: data, encoding: NSUTF8StringEncoding) as! String
        
//        print(output)
    }
    
    static func closeServer() {
        serverTask.interrupt()
    }
    
    func findGoServerDirectory() -> String? {
        // TODO: find the directory
        
        return ""
    }
    
}