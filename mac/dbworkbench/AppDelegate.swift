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
        // Insert code here to initialize your application
        print("App Opened")
        ServerController.runServer()
    }

    func applicationWillTerminate(aNotification: NSNotification) {
        // CMD + Q
        print("CMD + Q")
    }


}

