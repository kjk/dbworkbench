using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.IO;
using System.Windows.Forms;

namespace DatabaseWorkbench
{
    public partial class Form1 : Form
    {
        WebBrowser _webBrowser;
        private void InitializeComponent2()
        {
            Layout += Form1_Layout;
            Load += Form1_Load;
            SuspendLayout();
            _webBrowser = new WebBrowser()
            {
                AllowNavigation = true,
            };
            Controls.Add(_webBrowser);
            ResumeLayout(true);
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
            var path = FindGoBackendDirectory();
            if (path == "")
            {
                // TODO: log
                return false;
            }
            // TODO: start dbworkbench.exe
            return true;
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
