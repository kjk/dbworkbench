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
        let w = self.view.window
        let d = NSApp.delegate as! AppDelegate
        d.window = w;

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

        let dateFmt = NSDateFormatter()
        dateFmt.dateFormat = "'log-'yy-MM-dd'-mac.txt"
        let logPathName = dateFmt.stringFromDate(NSDate())
        if (!NSFileManager.defaultManager().fileExistsAtPath(dataPath)) {
            do {
                try NSFileManager.defaultManager().createDirectoryAtPath(dataPath, withIntermediateDirectories: true, attributes: nil)
            } catch _ {
                // failed
            }

        }
        
        let logDirectory: NSString = dataPath
        let logpath = logDirectory.stringByAppendingPathComponent(logPathName)
        freopen(logpath.cStringUsingEncoding(NSASCIIStringEncoding)!, "a+", stderr)
    }

}

