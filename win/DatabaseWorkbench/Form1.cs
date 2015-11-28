using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
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

        private void Form1_Load(object sender, EventArgs e)
        {
            Console.WriteLine($"{_webBrowser.AllowNavigation}");
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
