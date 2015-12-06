//
//  ViewController.swift
//  dbworkbench
//
//  Created by Furkan Yilmaz on 26/11/15.
//  Copyright © 2015 Furkan Yilmaz. All rights reserved.
//

import Cocoa
import WebKit

class ViewController: NSViewController {
    
    @IBOutlet weak var webView: WebView!

    let urlpath = "http://localhost:5444"
    var once = false
    
    override func awakeFromNib() {
        NSLog("awakeFromNib")
        
        if !once {
            redirectLogToDocuments()
            once = true
        }
        
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        NSLog("viewDidLoad")
        
        viewController = self
    }
    
    override func viewWillAppear() {
        super.viewWillAppear()
        NSLog("viewWillAppear")
        
        preferredContentSize = view.fittingSize
    }
    
    override func viewWillDisappear() {
        NSLog("viewWillDisappear")
        closeServer()
    }
    
    func loadURL() {
        NSLog("loadURL")
        let requesturl = NSURL(string: urlpath)
        let request = NSURLRequest(URL: requesturl!)
        
        webView.mainFrame.loadRequest(request)
    }
    
    func redirectLogToDocuments() {
        
//        let homeDirectory = NSHomeDirectory() as NSString
        let dataPath = NSString.pathWithComponents([NSHomeDirectory(), "Library", "Application Support", "Database Workbench", "log"])

        if (!NSFileManager.defaultManager().fileExistsAtPath(dataPath)) {
            do {
                try NSFileManager.defaultManager().createDirectoryAtPath(dataPath, withIntermediateDirectories: true, attributes: nil)
            } catch _ {
                // failed
            }

        }
        
        let logDirectory: NSString = dataPath
        let logpath = logDirectory.stringByAppendingPathComponent("maclog.txt")
        freopen(logpath.cStringUsingEncoding(NSASCIIStringEncoding)!, "a+", stderr)
    }

}

