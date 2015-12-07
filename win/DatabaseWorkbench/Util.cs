using System;
using System.Collections.Generic;
using System.IO;
using System.Management;

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

        public static string UsageFilePath()
        {
            return Path.Combine(AppDataDir(), "usage.json");
        }

        // TODO: move to yepi-utils
        // TODO: this doesn't seem to work without root priviledges
        // based on http://stackoverflow.com/questions/4084402/get-hard-disk-serial-number
        public static string[] GetHardDriveSerials()
        {
            var res = new List<string>();
            var searcher = new ManagementObjectSearcher("SELECT * FROM Win32_PhysicalMedia");

            foreach (ManagementObject hd in searcher.Get())
            {
                var serial = hd["SerialNumber"];
                if (serial != null)
                {
                    res.Add(serial.ToString());
                }
            }
            return res.ToArray();
        }

    }
}
