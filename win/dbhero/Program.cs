using System;
using System.IO;
using System.Threading;
using System.Runtime.ExceptionServices;
using System.Security.Permissions;
using System.Windows.Forms;

using Yepi;
using System.Threading.Tasks;

namespace DbHero
{
    static class Program
    {
        static Mutex mutex = new Mutex(true, "dbheroapp.com/dbhero");

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

        static void RunApp()
        {
            Application.ThreadException += new ThreadExceptionEventHandler(ThreadExceptionHandler);
            Application.SetUnhandledExceptionMode(UnhandledExceptionMode.CatchException);
            AppDomain.CurrentDomain.UnhandledException += new UnhandledExceptionEventHandler(UnhandledExceptionHandler);
            // http://stackoverflow.com/questions/20542212/winforms-app-still-crashes-after-unhandled-exception-handler
            TaskScheduler.UnobservedTaskException += TaskScheduler_UnobservedTaskException;

            Log.TryOpen(LogPath());
            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);
            Application.Run(new Form1());
        }

        [STAThread]
        [SecurityPermission(SecurityAction.Demand, Flags = SecurityPermissionFlag.ControlAppDomain)]
        static void Main()
        {
            if (mutex.WaitOne(TimeSpan.Zero, true))
            {
                RunApp();
                mutex.ReleaseMutex();
            }
            else
            {
                // send our Win32 message to make the currently running instance
                // jump on top of all the other windows
                NativeMethods.PostMessage(
                    (IntPtr)NativeMethods.HWND_BROADCAST,
                    NativeMethods.WM_SHOWME,
                    IntPtr.Zero,
                    IntPtr.Zero);
            }
        }

        // TODO: send crash report to the website
        static void ShowCrash(Exception e)
        {
            var msg = e.Message;
            if (msg.Length > 0)
                msg += "\n\n";
            if (e.StackTrace == null)
            {
                if (e is AggregateException)
                {
                    var ae = e as AggregateException;
                    var baseException = ae.GetBaseException();
                    if (baseException.StackTrace != null)
                    {
                        msg += baseException.StackTrace.ToString();
                    }
                }
            } else
            {
                msg += e.StackTrace.ToString();
            }
            MessageBox.Show("We're sorry, we crashed!\n\n" + msg, "dbHero crashed", MessageBoxButtons.OK);
        }

        private static void TaskScheduler_UnobservedTaskException(object sender, UnobservedTaskExceptionEventArgs e)
        {
            Log.E(e.Exception);
            Log.Close();
            ShowCrash(e.Exception);
            Application.Exit();
        }

        [HandleProcessCorruptedStateExceptionsAttribute]
        static void UnhandledExceptionHandler(object sender, UnhandledExceptionEventArgs args)
        {
            Exception e = args.ExceptionObject as Exception;
            if (e != null)
            {
                Log.E(e);
                ShowCrash(e);
            }
            Log.Close();
            Application.Exit();
        }

        [HandleProcessCorruptedStateExceptionsAttribute]
        private static void ThreadExceptionHandler(object sender, ThreadExceptionEventArgs t)
        {
            Log.E(t.Exception);
            Log.Close();
            ShowCrash(t.Exception);
            Application.Exit();
        }

    }
}
