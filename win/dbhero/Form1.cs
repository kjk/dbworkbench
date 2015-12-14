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
using System.Reflection;

using Yepi;
using System.Threading;

namespace DbHero
{
    public partial class Form1 : Form
    {
        WebBrowser _webBrowser;
        Process _backendProcess;
        MenuStrip _mainMenu;
        //string _websiteURL = "http://localhost:5555";
        string _websiteURL = "https://dbheroapp.com";
        bool _cleanFinish = false;
        string _updateInstallerPath;

        protected override void WndProc(ref Message m)
        {
            if (m.Msg == NativeMethods.WM_SHOWME)
            {
                BringMeToFront();
            }
            base.WndProc(ref m);
        }

        private void BringMeToFront()
        {
            if (WindowState == FormWindowState.Minimized)
            {
                WindowState = FormWindowState.Normal;
            }
            // get our current "TopMost" value (ours will always be false though)
            bool top = TopMost;
            // make our form jump to the top of everything
            TopMost = true;
            // set it back to whatever it was
            TopMost = top;
        }

        // could also use MainMenu http://stackoverflow.com/questions/2778109/standard-windows-menu-bars-in-windows-forms
        // in which case probably don't need to layout _mainMenu (it'll be part of non-client area)
        private void CreateMenu()
        {
            _mainMenu = new MenuStrip();

            var menuFile = new ToolStripMenuItem("&File");
            _mainMenu.Items.Add(menuFile);
            menuFile.DropDownItems.Add("&Exit", null, FileExit_Click);

            var menuView = new ToolStripMenuItem("&View");
            _mainMenu.Items.Add(menuView);
            menuView.DropDownItems.Add("Zoom In", null, ViewZoomIn_Click);
            menuView.DropDownItems.Add("Zoom Out", null, ViewZoomOut_Click);
            menuView.DropDownItems.Add("Zoom 100%", null, ViewZoom100_Click);

            var menuHelp = new ToolStripMenuItem("&Help");
            _mainMenu.Items.Add(menuHelp);
            menuHelp.DropDownItems.Add("&Website", null, HelpWebsite_Click);
            menuHelp.DropDownItems.Add("&Support", null, HelpSupport_Click);
            menuHelp.DropDownItems.Add("&Feedback", null, HelpFeedback_Click);

#if DEBUG
            menuHelp.DropDownItems.Add(new ToolStripSeparator());
            menuHelp.DropDownItems.Add("Diagnostic page", null, HelpDiagnosticPage_Click);
            menuHelp.DropDownItems.Add("Main page", null, HelpMainPage_Click);
            menuHelp.DropDownItems.Add("Crash main thread", null, HelpCrashMainThread_Click);
            menuHelp.DropDownItems.Add("Crash background thread", null, HelpCrashBackgroundThread_Click);
#endif

            Controls.Add(_mainMenu);
            MainMenuStrip = _mainMenu;
            _mainMenu.PerformLayout();
        }

        // Those work for me but wonder if they are reliable. There are other ways:
        // http://stackoverflow.com/questions/738232/zoom-in-on-a-web-page-using-webbrowser-net-control
        private void ViewZoomIn_Click(object sender, EventArgs e)
        {
            _webBrowser.Focus();
            SendKeys.Send("^{+}"); // [CTRL]+[+]
        }

        private void ViewZoomOut_Click(object sender, EventArgs e)
        { 
            _webBrowser.Focus();
            SendKeys.Send("^{-}"); // [CTRL]+[-]
        }

        private void ViewZoom100_Click(object sender, EventArgs e)
        {
            _webBrowser.Focus();
            SendKeys.Send("^0"); // [CTRL]+[0]
        }

        private void HelpDiagnosticPage_Click(object sender, EventArgs e)
        {
            _webBrowser.Navigate("http://127.0.0.1:5444/diagnostic.html");
        }

        private void HelpMainPage_Click(object sender, EventArgs e)
        {
            _webBrowser.Navigate("http://127.0.0.1:5444");
        }

        private void HelpCrashMainThread_Click(object sender, EventArgs e)
        {
            throw new NotImplementedException();
        }

        private void HelpCrashBackgroundThread_Click(object sender, EventArgs e)
        {
            Task.Run(() =>
            {
                throw new NotImplementedException();
            });
            Thread.Sleep(2000);
            // unobserved exceptions are only reported when finalized, so force
            // finalization
            GC.Collect();
        }

        private void HelpFeedback_Click(object sender, EventArgs e)
        {
            FileUtil.TryLaunchUrl("https://dbheroapp.com/feedback");
        }

        private void HelpSupport_Click(object sender, EventArgs e)
        {
            FileUtil.TryLaunchUrl("https://dbheroapp.com/support");
        }

        private void HelpWebsite_Click(object sender, EventArgs e)
        {
            FileUtil.TryLaunchUrl("https://dbheroapp.com");
        }

        private void FileExit_Click(object sender, EventArgs e)
        {
            this.Close();
        }

        // Find directory where dbherohelper.exe is
        // Returns "" if not found
        private string FindBackendDirectory()
        {
            var path = Application.ExecutablePath;
            var dir = Path.GetDirectoryName(path);
            while (dir != null)
            {
                path = Path.Combine(dir, "dbherohelper.exe");
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

        private bool StartBackendServer()
        {
            var dir = FindBackendDirectory();
            if (dir == "")
            {
                Log.S("StartBackendServer: FindBackendDirectory() failed\n");
                return false;
            }
            // explanation of StartInfo flags:
            // http://blogs.msdn.com/b/jmstall/archive/2006/09/28/createnowindow.aspx
            var p = new Process();
            _backendProcess = p;
            p.StartInfo.WorkingDirectory = dir;
            p.StartInfo.FileName = "dbherohelper.exe";
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
            Log.Line($"Backend_Exited. _clienFinish:{_cleanFinish}");
            if (_cleanFinish)
            {
                // we killed the process ourselves
                return;
            }
            // TODO: show better error message
            MessageBox.Show("backend exited unexpectedly!");
            Close();
        }


        // Note: can't be in utils, must be in this assembly
        public static string AppVer()
        {
            Assembly assembly = Assembly.GetExecutingAssembly();
            FileVersionInfo fvi = FileVersionInfo.GetVersionInfo(assembly.Location);
            return Utils.CleanAppVer(fvi.ProductVersion);
        }

        // must happen before StartBackendServer()
        string _backendUsage = "";
        public void ReadUsage()
        {
            // can happen because UsageFilePath() might not exist on first run
            // TODO: make it File.TryReadAllText()
            try
            {
                _backendUsage = File.ReadAllText(Util.UsageFilePath());
            }
            catch (Exception e)
            {
                Log.E(e);
            }
            // delete so that we don't the same data multiple times
            FileUtil.TryFileDelete(Util.UsageFilePath());
        }

        /* the data we POST as part of auto-update check is in format:
        ver: ${ver}
        os: 6.1
        other: ${other_val}
        -----------
        ${usage data from backend}
        */
        private string BuildAutoUpdatePostData()
        {
            var computerInfo = Util.GetComputerInfo();

            var s = "";
            s += $"ver: {AppVer()}\n";
            s += "ostype: windows\n";
            s += $"user: {computerInfo.UserName}\n";
            s += $"os: {computerInfo.OsVersion}\n";
            if (computerInfo.NetworkCardId != "")
            {
                s += $"networkCardId: {computerInfo.NetworkCardId}\n";
            }
            s += $"machine: {computerInfo.MachineName}\n";
            s += $"net: {computerInfo.InstalledNetVersions}\n";

            s += "---------------\n"; // separator
            if (_backendUsage != "")
            {
                s += _backendUsage + "\n";
            }
            return s;
        }

        private async Task AutoUpdateCheck()
        {
            // we might have downloaded installer previously, in which case
            // don't re-download
            var tmpInstallerPath = Util.UpdateInstallerTmpPath();
            if (File.Exists(tmpInstallerPath))
            {
                _updateInstallerPath = tmpInstallerPath;
                NotifyUpdateAvailable("");
                return;
            }

            var myVer = AppVer();
            var postData = BuildAutoUpdatePostData();
            Log.Line(postData);

            var url = _websiteURL + "/api/winupdatecheck?ver=" + myVer;
            var result = await Http.PostStringAsync(url, postData);
            if (result == null)
            {
                Log.Line("AutoUpdateCheck(): result is null");
                return;
            }

            string ver, dlUrl;
            var ok = Utils.ParseUpdateResponse(result, out ver, out dlUrl);
            if (!ok)
            {
                Log.Line($"AutoUpdateCheck: failed to parse ver '{ver}' in '{result}'");
                return;
            }

            if (!Utils.ProgramVersionGreater(ver, myVer))
            {
                Log.Line($"AutoUpdateCheck: not updating because my version{myVer}  is >= than latest available {ver}");
                return;
            }
            var d = await Http.UrlDownloadAsync(dlUrl);
            if (d == null)
            {
                Log.Line($"AutoUpdateCheck: failed to download {dlUrl}");
                return;
            }
            _updateInstallerPath = Util.UpdateInstallerTmpPath();
            try
            {
                File.WriteAllBytes(_updateInstallerPath, d);
            }
            catch
            {
                File.Delete(_updateInstallerPath);
                _updateInstallerPath = null;
                return;
            }
            NotifyUpdateAvailable(ver);
        }

        public void NotifyUpdateAvailable(string ver)
        {
            // TODO: show a nicer dialog
            if (ver != "")
            {
                ver = " " + ver;
            }
            var res = MessageBox.Show($"A new version{ver} is available. Update to new version?", "Update available", MessageBoxButtons.YesNo);
            if (res != DialogResult.Yes)
            {
                return;
            }
            Log.Line($"About to run an updater: {_updateInstallerPath}");
            // move the installer to another, temporary path, so that when the installation is finished
            // and we restart the app, we won't think an update is available
            var tmpInstallerPath = Path.GetTempFileName();
            File.Delete(tmpInstallerPath);
            tmpInstallerPath += ".exe";
            File.Move(_updateInstallerPath, tmpInstallerPath);
            _updateInstallerPath = null;
            FileUtil.TryLaunchUrl(tmpInstallerPath);
            // exit ourselves so that the installer can over-write the file
            Close();
        }

        private void Form1_FormClosing(object sender, FormClosingEventArgs e)
        {
            // TODO: weird. It looks like _backendProcess gets killed after we set _cleanFinish
            // but before Log.Line().
            _cleanFinish = true;
            if (_backendProcess != null)
            {
                Log.Line($"Form1_FormClosing: backend exited: {_backendProcess.HasExited}");
            }
            // TODO: if we have multiple forms, only when last form is being closed
            if (_backendProcess != null && !_backendProcess.HasExited)
            {
                Log.Line($"Form1_FormClosing: killing backend");
                _backendProcess.Kill();
            }
            Log.Line($"Loc: {Location}, Size: {Size}");
            Log.Close();
        }

        private async void Form1_Load(object sender, EventArgs e)
        {
            ReadUsage();
            if (!StartBackendServer())
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
            var s = _mainMenu.PreferredSize;
            var y = s.Height;
            var dy = area.Height - y;
            _webBrowser.SetBounds(0, y, area.Width, dy);
        }

        private void InitializeComponent2()
        {
            SuspendLayout();

            //Location = new Point(10, 10);
            //ClientSize = new Size(1033, 771);

            Layout += Form1_Layout;
            Load += Form1_Load;
            FormClosing += Form1_FormClosing;

            _webBrowser = new WebBrowser()
            {
                AllowNavigation = true,
                IsWebBrowserContextMenuEnabled = false,
            };
            Controls.Add(_webBrowser);

            CreateMenu();
            ResumeLayout(true);
        }

        public Form1()
        {
            InitializeComponent();
            InitializeComponent2();
        }
    }
}
