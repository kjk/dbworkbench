using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Drawing;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Diagnostics;
using System.IO;
using System.Windows.Forms;
using System.Net.Http;
using System.Reflection;
using System.Runtime.Serialization.Json;

using CefSharp;
using CefSharp.WinForms;

using Yepi;

namespace DbHero
{
    public partial class Form1 : Form, IKeyboardHandler
    {
        ChromiumWebBrowser _webBrowser;
        Panel _webBrowserPanel;
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

        public bool OnKeyEvent(IWebBrowser browserControl, IBrowser browser, KeyType type, int windowsKeyCode, int nativeKeyCode, CefEventFlags modifiers, bool isSystemKey)
        {
            //Console.WriteLine($"OnKeyEvent: type={type}, windowsKeyCode={windowsKeyCode}, nativeKeyCode={nativeKeyCode}, modifiers={modifiers}, isSystemKey={isSystemKey}");

            // Ctrl - +
            if (type == KeyType.KeyUp && windowsKeyCode == 187 && modifiers == CefEventFlags.ControlDown)
            {
                ViewZoomIn_Click(null, null);
                //Console.WriteLine("Ctrl +");
                return true;
            }

            // Ctrl - -
            if (type == KeyType.KeyUp && windowsKeyCode == 189 && modifiers == CefEventFlags.ControlDown)
            {
                ViewZoomOut_Click(null, null);
                //Console.WriteLine("Ctrl -");
                return true;
            }

            // Ctrl - 0
            if (type == KeyType.KeyUp && windowsKeyCode == 48 && modifiers == CefEventFlags.ControlDown)
            {
                ViewZoom100_Click(null, null);
                //Console.WriteLine("Ctrl 0");
                return true;
            }

            return false;
        }

        public bool OnPreKeyEvent(IWebBrowser browserControl, IBrowser browser, KeyType type, int windowsKeyCode, int nativeKeyCode, CefEventFlags modifiers, bool isSystemKey, ref bool isKeyboardShortcut)
        {
            return false;
        }

        // could also use MainMenu http://stackoverflow.com/questions/2778109/standard-windows-menu-bars-in-windows-forms
        // in which case probably don't need to layout _mainMenu (it'll be part of non-client area)
        private void CreateMenu()
        {
            _mainMenu = new MenuStrip();

            var menuFile = new ToolStripMenuItem("&File");
            _mainMenu.Items.Add(menuFile);
            menuFile.DropDownItems.Add("&Exit", null, FileExit_Click);

            // TODO: the shortcuts don't work. Probably swollowed by the chrome control.
            // probably need to use KeyboardHandler on _webControl
            var menuView = new ToolStripMenuItem("&View");
            _mainMenu.Items.Add(menuView);
            var mi = new ToolStripMenuItem("Actual Size", null, ViewZoom100_Click);
            mi.ShortcutKeys = Keys.Control | Keys.D0;
            mi.ShortcutKeyDisplayString = "Ctrl+0";
            menuView.DropDownItems.Add(mi);
            mi = new ToolStripMenuItem("Zoom In", null, ViewZoomIn_Click);
            mi.ShortcutKeys = Keys.Control | Keys.Add;
            mi.ShortcutKeyDisplayString = "Ctrl-+";
            menuView.DropDownItems.Add(mi);
            mi = new ToolStripMenuItem("Zoom Out", null, ViewZoomOut_Click);
            mi.ShortcutKeyDisplayString = "Ctrl--";
            mi.ShortcutKeys = Keys.Control | Keys.OemMinus;
            menuView.DropDownItems.Add(mi);

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

        // same as Chrome
        double[] zoomLevels = new double[]{ 0.25, 0.33, 0.5, 0.67, 0.75, 0.9, 1, 1.1, 1.25, 1.5, 1.75, 2, 2.5, 3, 4, 5 };

        // returns an index at which zoomLevel[i] has the smallest difference from z
        private int FindClosestZoomLevelIdx(double z)
        {
            // silly control uses 0 for default size
            if (z == 0)
                z = 1.0;

            var min = zoomLevels[0];
            if (z <= min)
                return 0;

            var n = zoomLevels.Length;
            var max = zoomLevels[n - 1];
            if (z >= max)
            {
                return n - 1;
            }
            var prevDiff = Math.Abs(min - z);
            for (var i = 1; i < n; i++)
            {
                var diff = Math.Abs(zoomLevels[i] - z);
                if (diff > prevDiff)
                    return i - 1;

                prevDiff = diff;
            }
            return n - 1;
        }

        private double fixZoomOut(double z)
        {
            if (z == 1)
                return 0;
            return z;
        }

        private double FindNextZoomLevel(double z)
        {
            var idx = FindClosestZoomLevelIdx(z) + 1;
            if (idx >= zoomLevels.Length)
                idx = zoomLevels.Length - 1;

            return zoomLevels[idx];
        }

        private double FindPrevZoomLevel(double z)
        {
            var idx = FindClosestZoomLevelIdx(z) - 1;
            if (idx < 0)
                idx = 0;

            return zoomLevels[idx];
        }

        private async void ViewZoomIn_Click(object sender, EventArgs e)
        {
            var currZoom = await _webBrowser.GetZoomLevelAsync();
            var newZoom = FindNextZoomLevel(currZoom);
            _webBrowser.SetZoomLevel(newZoom);
            Console.WriteLine($"zoom in: currZoom: {currZoom}, newZoom: {newZoom}, newFixed: {fixZoomOut(newZoom)}");
        }

        private async void ViewZoomOut_Click(object sender, EventArgs e)
        {
            var currZoom = await _webBrowser.GetZoomLevelAsync();
            var newZoom = FindPrevZoomLevel(currZoom);
            _webBrowser.SetZoomLevel(newZoom);
            Console.WriteLine($"zoom out: currZoom: {currZoom}, newZoom: {newZoom}, newFixed: {fixZoomOut(newZoom)}");
        }

        private void ViewZoom100_Click(object sender, EventArgs e)
        {
            _webBrowser.SetZoomLevel(1.0);
        }

        private void HelpDiagnosticPage_Click(object sender, EventArgs e)
        {
            _webBrowser.Load("http://127.0.0.1:5444/diagnostic.html");
        }

        private void HelpMainPage_Click(object sender, EventArgs e)
        {
            _webBrowser.Load("http://127.0.0.1:5444");
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
            var computerInfo = Utils.GetComputerInfo();

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

        void SaveSettings()
        {
            var s = new Settings();
            s.WindowPosition = Location;
            s.WindowSize = ClientSize;
            var mem = new MemoryStream();
            var ser = new DataContractJsonSerializer(typeof(Settings));
            ser.WriteObject(mem, s);
            var d = mem.ToArray();
            FileUtil.TryWriteAllBytes(Util.SettingsFilePath(), d);
        }

        Settings LoadSettings()
        {
            var path = Util.SettingsFilePath();
            var d = FileUtil.TryReadAllBytes(path);
            if (d == null)
                return null;
            var ser = new DataContractJsonSerializer(typeof(Settings));
            try
            {
                using (var mem = new MemoryStream(d))
                {
                    return (Settings)ser.ReadObject(mem);
                }
            }
            catch (Exception e)
            {
                Log.E(e);
                // I assume it failed because the structure changed
                // so delete the file to not get this again
                FileUtil.TryFileDelete(path);
                return null;
            }
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
            SaveSettings();
            Log.Close();
            Cef.Shutdown();
        }

        private async void Form1_Load(object sender, EventArgs e)
        {
            var settings = LoadSettings();
            if (settings != null)
            {
                Location = settings.WindowPosition;
                ClientSize = settings.WindowSize;
            }

            ReadUsage();
            if (!StartBackendServer())
            {
                // TODO: better way to show error message
                MessageBox.Show("Backend didn't start. Quitting.");
                Close();
                return;
            }
            CreateBrowser("http://127.0.0.1:5444");
            // TODO: doing it here throws an exception
            // probably need to wait for the content to load
            //ViewZoom100_Click(null, null);
            await AutoUpdateCheck();
        }

        private void Form1_Layout(object sender, LayoutEventArgs e)
        {
            var area = ClientSize;
            var s = _mainMenu.PreferredSize;
            var y = s.Height;
            var dy = area.Height - y;
            _webBrowserPanel.SetBounds(0, y, area.Width, dy);
        }

        private void CreateBrowser(string url)
        {
            Cef.Initialize(new CefSettings());
            _webBrowser = new ChromiumWebBrowser(url);
            _webBrowser.KeyboardHandler = this;
            // set the zoom as 1.0. It's different than 0 that represents default
            // zoom, but if I try to use 0 intelligently, zoom in/out is unpredictable
            _webBrowserPanel.Controls.Add(_webBrowser);
        }

        private void InitializeComponent2()
        {
            SuspendLayout();

            Layout += Form1_Layout;
            Load += Form1_Load;
            FormClosing += Form1_FormClosing;
            //CreateBrowser("");
            CreateMenu();
            _webBrowserPanel = new Panel();
            //_webBrowserPanel.BackColor = Color.Red;
            Controls.Add(_webBrowserPanel);
            ResumeLayout(true);
        }

        public Form1()
        {
            InitializeComponent();
            InitializeComponent2();
        }
    }
}
