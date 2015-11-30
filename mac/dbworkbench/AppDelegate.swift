//
//  AppDelegate.swift
//  dbworkbench
//
//  Created by Furkan Yilmaz on 26/11/15.
//  Copyright © 2015 Furkan Yilmaz. All rights reserved.
//

import Cocoa

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate {

    func applicationDidFinishLaunching(aNotification: NSNotification) {
        NSLog("App Opened")
        ServerController.runServer()
    }

    func applicationWillTerminate(aNotification: NSNotification) {
        NSLog("applicationWillTerminate")
    }
    
    func applicationShouldTerminateAfterLastWindowClosed(sender: NSApplication) -> Bool {
        return true
    }
}

