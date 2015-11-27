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

    var urlpath = "http://localhost:5444/"
    
    func loadAddressURL(){
        let requesturl = NSURL(string: urlpath)
        let request = NSURLRequest(URL: requesturl!)
        
//        webView.frame = CGRectMake(0, 0, self.view.frame.size.width, self.view.frame.size.height);
//        self.view = webView;
        webView.mainFrame.loadRequest(request)
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        loadAddressURL() //func above
        
        // Do any additional setup after loading the view, typically from a nib.
    }
    
    override func viewDidLayout() {
//        self.view = webView;
    }
}

