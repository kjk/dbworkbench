using System;
using System.IO;

namespace DatabaseWorkbench
{
    class Util
    {
        public static string AppDataDir()
        {
            var dir = Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData);
            dir = System.IO.Path.Combine(dir, "Database Workbench");
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
            return Path.Combine(AppDataTmpDir(), "database-workbench-installer.exe");
        }
    }
}
