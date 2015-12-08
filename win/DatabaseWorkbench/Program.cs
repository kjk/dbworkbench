using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using System.Windows.Forms;

using Yepi;

namespace DatabaseWorkbench
{
    static class Program
    {
        static string LogPath()
        {
            var logDir = Util.AppDataLogDir();
            DateTime now = DateTime.Now;
            int y = now.Year;
            int m = now.Month;
            int d = now.Day;
            string logFilePath = Path.Combine(logDir, String.Format("log-{0:0000}-{1:00}-{2:00}-win.txt", y, m, d));
            return logFilePath;
        }

        [STAThread]
        static void Main()
        {
            Log.TryOpen(LogPath());
            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);
            Application.Run(new Form1());
        }
    }
}
