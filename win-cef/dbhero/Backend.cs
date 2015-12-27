using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;
using Yepi;

namespace DbHero
{
    public static class Backend
    {
        public static Process Process;
        public static EventHandler exitHandler;

        // Find directory where dbherohelper.exe is
        // Returns "" if not found
        private static string FindBackendDirectory()
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


        public static bool Start(EventHandler Exited)
        {
            if (Process != null)
            {
                Log.Line("Backend.Start: has already been started");
                return false;
            }
            var dir = FindBackendDirectory();
            if (dir == "")
            {
                Log.S("Backend.Start: FindBackendDirectory() failed\n");
                return false;
            }
            exitHandler = Exited;
            // explanation of StartInfo flags:
            // http://blogs.msdn.com/b/jmstall/archive/2006/09/28/createnowindow.aspx
            var p = new Process();
            Process = p;
            p.EnableRaisingEvents = true;
            p.StartInfo.WorkingDirectory = dir;
            p.StartInfo.FileName = "dbherohelper.exe";
            p.StartInfo.UseShellExecute = true;
            p.StartInfo.CreateNoWindow = true;
            p.StartInfo.WindowStyle = ProcessWindowStyle.Hidden;
            p.Exited += P_Exited;
            var ok = p.Start();
            return ok;
        }

        private static void P_Exited(object sender, EventArgs e)
        {
            if (exitHandler != null)
            {
                exitHandler(sender, e);
            }
        }

        public static void Stop()
        {
            if (Process == null)
                return;
            // don't notify about exit if it was killed intentionally by us
            exitHandler = null;
            if (!Process.HasExited)
            {
                Process.Kill();
            }
            Process = null;
        }

    }
}
