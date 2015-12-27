using System;
using System.IO;
using System.Runtime.InteropServices;

namespace DbHero
{
    // this class just wraps some Win32 stuffthat we're going to use
    internal class NativeMethods
    {
        public const int HWND_BROADCAST = 0xffff;
        public static readonly int WM_SHOWME = RegisterWindowMessage("WM_SHOW_DBHERO");
        [DllImport("user32")]
        public static extern bool PostMessage(IntPtr hwnd, int msg, IntPtr wparam, IntPtr lparam);
        [DllImport("user32")]
        public static extern int RegisterWindowMessage(string message);
    }

    class Util
    {
        public static string AppDataDir()
        {
            var dir = Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData);
            dir = System.IO.Path.Combine(dir, "dbHero");
            Directory.CreateDirectory(dir);
            return dir;
        }

        public static string AppDataLogDir()
        {
            var dir = System.IO.Path.Combine(AppDataDir(), "log");
            Directory.CreateDirectory(dir);
            return dir;
        }

        public static string AppDataTmpDir()
        {
            var dir = System.IO.Path.Combine(AppDataDir(), "tmp");
            Directory.CreateDirectory(dir);
            return dir;
        }

        public static string UpdateInstallerTmpPath()
        {
            return Path.Combine(AppDataTmpDir(), "dbhero-installer-tmp.exe");
        }

        public static string UsageFilePath()
        {
            return Path.Combine(AppDataDir(), "usage.json");
        }

        public static string SettingsFilePath()
        {
            return Path.Combine(AppDataDir(), "settings.json");
        }

    }
}
