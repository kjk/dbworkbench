using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Diagnostics;
using System.IO;
using System.Windows.Forms;
using System.Net.Http;

namespace DatabaseWorkbench
{
    public partial class Form1 : Form
    {
        WebBrowser _webBrowser;
        Process _backendProcess;
        string _websiteURL = "http://localhost:5555";
        //string _websiteURL = "http://databaseworkbench.com";
        bool _cleanFinish = false;

        private void InitializeComponent2()
        {
            Layout += Form1_Layout;
            Load += Form1_Load;
            FormClosing += Form1_FormClosing;
            SuspendLayout();
            _webBrowser = new WebBrowser()
            {
                AllowNavigation = true,
            };
            Controls.Add(_webBrowser);
            ResumeLayout(true);
        }

        private void Form1_FormClosing(object sender, FormClosingEventArgs e)
        {
            // TODO: weird. It looks like _backendProcess gets killed after we set _cleanFinish
            // but before Console.WriteLine().
            _cleanFinish = true;
            Console.WriteLine($"Form1_FormClosing: backend exited: {_backendProcess.HasExited}");
            // TODO: if we have multiple forms, only when last form is being closed
            if (_backendProcess != null && !_backendProcess.HasExited)
            {
                Console.WriteLine($"Form1_FormClosing: killing backend");
                _backendProcess.Kill();
            }
        }

        // Find directory where dbworkbench.exe is
        // Returns "" if not found
        private string FindGoBackendDirectory()
        {
            var path = Application.ExecutablePath;
            var dir = Path.GetDirectoryName(path);
            while (dir != null)
            {
                path = Path.Combine(dir, "dbworkbench.exe");
                if (File.Exists(path))
                {
                    return dir;
                }
                var newDir = Path.GetDirectoryName(dir);
                if (dir == newDir)
                {
                    return "";
                }
                dir = newDir;
            }
            return "";
        }

        private bool StartGoBackend()
        {
            var dir = FindGoBackendDirectory();
            if (dir == "")
            {
                // TODO: log
                return false;
            }
            // explanation of StartInfo flags:
            // http://blogs.msdn.com/b/jmstall/archive/2006/09/28/createnowindow.aspx
            var p = new Process();
            _backendProcess = p;
            p.StartInfo.WorkingDirectory = dir;
            p.StartInfo.FileName = "dbworkbench.exe";
            p.StartInfo.UseShellExecute = true;
            p.StartInfo.CreateNoWindow = true;
            p.StartInfo.WindowStyle = ProcessWindowStyle.Hidden;
            p.Exited += Backend_Exited;
            var ok = p.Start();
            return ok;
        }

        // happens when go backend process has finished e.g. because of an error
        private void Backend_Exited(object sender, EventArgs e)
        {
            Console.WriteLine($"Backend_Exited. _clienFinish:{_cleanFinish}");
            if (_cleanFinish)
            {
                // we killed the process ourselves
                return;
            }
            // TODO: show better error message
            MessageBox.Show("backend exited unexpectedly!");
            Close();
        }

        // returns null if failed to download
        private async Task<string> UrlDownloadAsStringAsync(string uri)
        {
            try
            {
                using (HttpClient client = new HttpClient())
                using (HttpResponseMessage response = await client.GetAsync(uri))
                using (HttpContent content = response.Content)
                {
                    if (response.StatusCode != System.Net.HttpStatusCode.OK)
                    {
                        // TODO: log to a file
                        Console.WriteLine($"UrlDownloadAsStringAsync(): response.StatusCode: {response.StatusCode}");
                        return null;
                    }
                    string result = await content.ReadAsStringAsync();
                    return result;
                }
            }
            catch
            {
                Console.WriteLine($"UrlDownloadAsStringAsync: exception happened");
                return null;
            }
        }

        // returns ver, downloadUrl
        private Tuple<string, string> ParseAutoUpdateCheckResponse(string s)
        {
            var downloadUrl = "";
            var ver = "";
            var lines = s.Split(new string[] { "\n" }, StringSplitOptions.RemoveEmptyEntries);
            foreach (var line in lines)
            {
                var parts = line.Split(new char[] { ':' }, 2);
                if (parts.Length != 2)
                {
                    Console.WriteLine($"AutoUpdateCheck: unexpected line: {line}");
                    continue;
                }
                var name = parts[0];
                var val = parts[1];
                if (name == "ver")
                {
                    ver = val;
                }
                else if (name == "url")
                {
                    downloadUrl = val;
                }
            }
            if (ver == "" || downloadUrl == "")
            {
                Console.WriteLine($"AutoUpdateCheck: unexpected response, missing 'ver' or 'url' lines");
            }
            return new Tuple<string, string>(ver, downloadUrl);
        }

        private async Task AutoUpdateCheck()
        {
            var myVer = "0.1"; // TODO: get from AssemblyInfo.cs
            var uri = _websiteURL + "/api/winupdatecheck?ver=" + myVer;
            var result = await UrlDownloadAsStringAsync(uri);
            if (result == null)
            {
                // TODO: log to a file
                Console.WriteLine("AutoUpdateCheck(): result is null");
                return;
            }

            Console.WriteLine($"result: {result}");
            var verUrl = ParseAutoUpdateCheckResponse(result);
            var ver = verUrl.Item1;
            // TODO: only trigger auto-update if ver > myVer
            if (ver == "" || ver == myVer)
            {
                Console.WriteLine($"AutoUpdateCheck: latest version {ver} is same as mine {myVer}");
                return;
            }
        }

        private async void Form1_Load(object sender, EventArgs e)
        {
            if (!StartGoBackend())
            {
                // TODO: better way to show error message
                MessageBox.Show("Backend didn't start. Quitting.");
                Close();
                return;
            }
            _webBrowser.Navigate("http://127.0.0.1:5444");
            await AutoUpdateCheck();
        }

        private void Form1_Layout(object sender, LayoutEventArgs e)
        {
            var area = ClientSize;
            _webBrowser.SetBounds(0, 0, area.Width, area.Height);
        }

        public Form1()
        {
            InitializeComponent();
            InitializeComponent2();
        }
    }
}
