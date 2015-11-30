//
//  ServerController.swift
//  dbworkbench
//
//  Created by Furkan Yilmaz on 28/11/15.
//  Copyright © 2015 Furkan Yilmaz. All rights reserved.
//

import Cocoa

class ServerController {
    
    static var serverTask =  NSTask()
    static var once = false
    
    static func runServer() {        
        let resPath = NSBundle.mainBundle().resourcePath
        let serverGoExePath = resPath! + "/dbworkbench.exe"
        
        serverTask.launchPath = serverGoExePath
        serverTask.currentDirectoryPath = resPath!
//        serverTask.arguments = ["-dev"]

        let pipe = NSPipe()
        serverTask.standardOutput = pipe
        serverTask.standardError = pipe

        
        let outHandle = pipe.fileHandleForReading
        outHandle.waitForDataInBackgroundAndNotify()
        
        let _ = NSNotificationCenter.defaultCenter().addObserverForName(NSFileHandleDataAvailableNotification, object: outHandle, queue: nil, usingBlock: { notification -> Void in
            let output = outHandle.availableData
            let outstr = NSString(data: output, encoding: NSUTF8StringEncoding)
            
            
            let qualityOfServiceClass = QOS_CLASS_BACKGROUND
            let backgroundQueue = dispatch_get_global_queue(qualityOfServiceClass, 0)
            dispatch_async(backgroundQueue, {
                
                if outstr?.length > 0 {
                    NSLog(outstr! as String)
                }
                
                dispatch_async(dispatch_get_main_queue(), { () -> Void in
                    if !once {
                        SingletonObject.viewController.loadURL()
                        once = true
                    }
                })
            })
            
            
            outHandle.waitForDataInBackgroundAndNotify()
            
        })
        
        serverTask.launch()
        NSLog("Server Started")
    }
    
    static func closeServer() {
        if serverTask.running {
            NSLog("Closing Server")
            serverTask.terminate()
        }
    }
}

class SingletonObject {
    static var viewController = ViewController()
}