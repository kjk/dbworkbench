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

namespace DatabaseWorkbench
{
    public partial class Form1 : Form
    {
        WebBrowser _webBrowser;
        Process _backendProcess;
        string _websiteURL = "http://localhost:5555";
        //string _websiteURL = "http://databaseworkbench.com";
        bool _cleanFinish = false;
        string _updateInstallerPath;

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
            if (_backendProcess != null)
            {
                Console.WriteLine($"Form1_FormClosing: backend exited: {_backendProcess.HasExited}");
            }
            // TODO: if we have multiple forms, only when last form is being closed
            if (_backendProcess != null && !_backendProcess.HasExited)
            {
                Console.WriteLine($"Form1_FormClosing: killing backend");
                _backendProcess.Kill();
            }
        }

        // Find directory where dbworkbench.exe is
        // Returns "" if not found
        private string FindBackendDirectory()
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

        private bool StartBackendServer()
        {
            var dir = FindBackendDirectory();
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

        public static string AppVer()
        {
            Assembly assembly = Assembly.GetExecutingAssembly();
            FileVersionInfo fvi = FileVersionInfo.GetVersionInfo(assembly.Location);
            return Utils.CleanAppVer(fvi.ProductVersion);
        }

        // must happen before StartBackendServer()
        string _backendUsage = "";
        public void LoadUsage()
        {
            // can happen because UsageFilePath() might not exist on first run
            // TODO: make it File.TryReadAllText()
            try
            {
                _backendUsage = File.ReadAllText(Util.UsageFilePath());
            }
            catch (Exception e)
            {
                Log.Le(e);
            }
            Utils.TryFileDelete(Util.UsageFilePath());
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
            var s = "";
            s += "program_ver: " + AppVer() + "\n";
            s += "os: " + "windows" + "\n"; // TODO: the exact os version

            var cardId = Util.GetNetworkCardId();
            if (cardId != "")
            {
                s += $"networkCardId: {cardId}\n";
            }
            // TODO: some unique id of the machine
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
                NotifyUpdateAvailable();
                return;
            }

            var myVer = AppVer();
            // TODO: make it a POST request
            var postData = BuildAutoUpdatePostData();
            Log.L(postData);

            var uri = _websiteURL + "/api/winupdatecheck?ver=" + myVer;
            var result = await Http.UrlDownloadAsStringAsync(uri);
            if (result == null)
            {
                // TODO: log to a file
                Console.WriteLine("AutoUpdateCheck(): result is null");
                return;
            }

            Console.WriteLine($"result: {result}");
            string ver, dlUrl;
            var ok = Utils.ParseUpdateResponse(result, out ver, out dlUrl);
            if (!ok)
            {

                return;
            }
            // TODO: only trigger auto-update if ver > myVer
            if (ver == "" || ver == myVer)
            {
                Log.L($"AutoUpdateCheck: latest version {ver} is same as mine {myVer}");
                return;
            }
            var d = await Http.UrlDownloadAsync(dlUrl);
            if (d == null)
            {
                Console.WriteLine($"AutoUpdateCheck: failed to download {dlUrl}");
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
            NotifyUpdateAvailable();
        }

        public void NotifyUpdateAvailable()
        {
            // TODO: show a nicer dialog
            var res = MessageBox.Show("Update available. Update?", "Update available", MessageBoxButtons.YesNo);
            if (res != DialogResult.Yes)
            {
                return;
            }
            Console.WriteLine($"should run an updater {_updateInstallerPath}");
            // move the installer to another, temporary path, so that when the installation is finished
            // and we restart the app, we won't think an update is available
            var tmpInstallerPath = Path.GetTempFileName();
            File.Delete(tmpInstallerPath);
            tmpInstallerPath += ".exe";
            File.Move(_updateInstallerPath, tmpInstallerPath);
            _updateInstallerPath = null;
            Utils.TryLaunchUrl(tmpInstallerPath);
            // exit ourselves so that the installer can over-write the file
            Close();
        }

        private async void Form1_Load(object sender, EventArgs e)
        {
            LoadUsage();
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
            _webBrowser.SetBounds(0, 0, area.Width, area.Height);
        }

        public Form1()
        {
            InitializeComponent();
            InitializeComponent2();
        }
    }
}
