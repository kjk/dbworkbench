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
            loadUsageData()
            runServer(self)
            redirectNSLogToFile()
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
    }

