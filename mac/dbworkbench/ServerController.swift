//
//  ServerController.swift
//  dbworkbench
//
//  Created by Furkan Yilmaz on 28/11/15.
//  Copyright © 2015 Furkan Yilmaz. All rights reserved.
//

import Cocoa

var viewController = ViewController()

var serverTask =  NSTask()
var didLoadURL = false

func runServer() {
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
                if !didLoadURL {
                    didLoadURL = true
                    viewController.loadURL()
                }
            })
        })
        
        
        outHandle.waitForDataInBackgroundAndNotify()
        
    })
    
    serverTask.launch()
    NSLog("server started")
}

func closeServer() {
    if serverTask.running {
        NSLog("closing server")
        serverTask.terminate()
    }
}

