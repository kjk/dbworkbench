using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.IO;
using System.Net.Http;
using System.Diagnostics;

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

        // per some info on the web, Process.Start(url) might throw
        // an exception, so swallow it
        public static void TryLaunchUrl(string url)
        {
            try
            {
                Process.Start(url);
            }
            catch (Exception e)
            {
                Console.WriteLine($"{e.ToString()}");
            }
        }

        // returns ver, downloadUrl
        public static Tuple<string, string> ParseAutoUpdateCheckResponse(string s)
        {
            var downloadUrl = "";
            var ver = "";
            var lines = s.Split(new string[] { "\n" }, StringSplitOptions.RemoveEmptyEntries);
            foreach (var line in lines)
            {
                var parts = line.Split(new char[] { ':' }, 2);
                if (parts.Length != 2)
                {
                    Console.WriteLine($"AutoUpdateCheck: unexpected line: {line}");
                    continue;
                }
                var name = parts[0];
                var val = parts[1];
                if (name == "ver")
                {
                    ver = val;
                }
                else if (name == "url")
                {
                    downloadUrl = val;
                }
            }
            if (ver == "" || downloadUrl == "")
            {
                Console.WriteLine($"AutoUpdateCheck: unexpected response, missing 'ver' or 'url' lines");
            }
            return new Tuple<string, string>(ver, downloadUrl);
        }

        // returns null if failed to download
        public static async Task<string> UrlDownloadAsStringAsync(string uri)
        {
            try
            {
                using (HttpClient client = new HttpClient())
                {
                    var res = await client.GetStringAsync(uri);
                    return res;
                }
            }
            catch
            {
                Console.WriteLine($"UrlDownloadAsStringAsync: exception happened");
                return null;
            }
        }

        // returns null if failed to download
        public static async Task<byte[]> UrlDownloadAsync(string uri)
        {
            try
            {
                using (HttpClient client = new HttpClient())
                {
                    var res = await client.GetByteArrayAsync(uri);
                    return res;
                }
            }
            catch
            {
                Console.WriteLine($"UrlDownloadAsync: exception happened");
                return null;
            }
        }
    }
}
