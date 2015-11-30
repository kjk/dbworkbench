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
        
        let resPath = NSBundle.mainBundle().resourcePath
        let serverGoExePath = resPath! + "/dbworkbench.exe"
        print(resPath)
        
        
        serverTask.launchPath = serverGoExePath
        serverTask.currentDirectoryPath = resPath!
//        serverTask.arguments = ["-dev"]

        let pipe = NSPipe()
        serverTask.standardOutput = pipe
        serverTask.standardError = pipe
        
        serverTask.launch()
        
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