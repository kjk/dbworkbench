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

namespace DatabaseWorkbench
{
    public partial class Form1 : Form
    {
        WebBrowser _webBrowser;
        Process _backendProcess;
        bool _cleanFinish = false;

        private void InitializeComponent2()
        {
            Layout += Form1_Layout;
            Load += Form1_Load;
            this.FormClosed += Form1_FormClosed;
            SuspendLayout();
            _webBrowser = new WebBrowser()
            {
                AllowNavigation = true,
            };
            Controls.Add(_webBrowser);
            ResumeLayout(true);
        }

        private void Form1_FormClosed(object sender, FormClosedEventArgs e)
        {
            if (_backendProcess != null && !_backendProcess.HasExited)
            {
                _cleanFinish = true;
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
            var p = new Process();
            _backendProcess = p;
            p.StartInfo.WorkingDirectory = dir;
            p.StartInfo.FileName = "dbworkbench.exe";
            p.StartInfo.UseShellExecute = true;
            p.StartInfo.CreateNoWindow = true;
            p.StartInfo.WindowStyle = ProcessWindowStyle.Hidden;
            p.Exited += Process_Exited;
            var ok = p.Start();
            return ok;
        }

        // happens when go backend process has finished e.g. because of an error
        private void Process_Exited(object sender, EventArgs e)
        {
            // TODO:
            // - show error message
            // - exit the application
            if (_cleanFinish)
            {
                // we killed the process ourselves
                return;
            }
        }

        private void Form1_Load(object sender, EventArgs e)
        {
            if (!StartGoBackend())
            {
                // TODO: show error message
                Close();
                return;
            }
            _webBrowser.Navigate("http://127.0.0.1:5444");
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
