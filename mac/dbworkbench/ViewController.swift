import Cocoa
import WebKit

class ViewController: NSViewController {
    
    @IBOutlet weak var webView: WebView!
    
    let urlpath = "http://localhost:5444"
    var awakeFromNibHappened = false
    
    override func awakeFromNib() {
        log("awakeFromNib")
        super.awakeFromNib()

        // can happen more than once
        if !awakeFromNibHappened {
            loadUsageData()
            startBackend(self)
            awakeFromNibHappened = true
        }
    }
    
    override func viewDidLoad() {
        log("viewDidLoad")
        super.viewDidLoad()
    }
    
    override func viewWillAppear() {
        log("viewWillAppear")
        super.viewWillAppear()
        let w = self.view.window
        getAppDelegate().window = w;
        preferredContentSize = view.fittingSize
    }
    
    override func viewWillDisappear() {
        log("viewWillDisappear")
        stopBackend()
    }
    
    func loadURL() {
        log("loadURL")
        let requesturl = NSURL(string: urlpath)
        let request = NSURLRequest(URL: requesturl!)
        
        webView.mainFrame.loadRequest(request)
    }
}

