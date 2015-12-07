import Cocoa
import WebKit

class ViewController: NSViewController {
    
    @IBOutlet weak var webView: WebView!

    let urlpath = "http://localhost:5444"
    var awakeFromNibHappened = false
    
    override func awakeFromNib() {
        super.awakeFromNib()

        NSLog("awakeFromNib")

        // can happen more than once
        if !awakeFromNibHappened {
            runServer(self)
            redirectLogToDocuments()
            awakeFromNibHappened = true
        }
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        NSLog("viewDidLoad")
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

