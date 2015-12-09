import Cocoa
import WebKit

class ViewController: NSViewController {
    
    @IBOutlet weak var webView: WebView!

    let urlpath = "http://localhost:5444"
    var awakeFromNibHappened = false
    
    override func awakeFromNib() {
        super.awakeFromNib()

        log("awakeFromNib")

        // can happen more than once
        if !awakeFromNibHappened {
            loadUsageData()
            runServer(self)
            awakeFromNibHappened = true
        }
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        log("viewDidLoad")
    }
    
    override func viewWillAppear() {
        super.viewWillAppear()
        let w = self.view.window
        let d = NSApp.delegate as! AppDelegate
        d.window = w;

        log("viewWillAppear")
        
        preferredContentSize = view.fittingSize
    }
    
    override func viewWillDisappear() {
        log("viewWillDisappear")
        closeServer()
    }
    
    func loadURL() {
        log("loadURL")
        let requesturl = NSURL(string: urlpath)
        let request = NSURLRequest(URL: requesturl!)
        
        webView.mainFrame.loadRequest(request)
    }
    }

