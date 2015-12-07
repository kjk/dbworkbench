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

#if false
        // TODO: move to yepi-utils
        // TODO: this doesn't seem to work without root priviledges
        // based on http://stackoverflow.com/questions/4084402/get-hard-disk-serial-number
        public static string[] GetHardDriveSerials()
        {
            var res = new List<string>();
            var query = new ManagementObjectSearcher("SELECT * FROM Win32_PhysicalMedia");

            foreach (ManagementObject o in query.Get())
            {
                var serial = o["SerialNumber"];
                if (serial != null)
                {
                    res.Add(serial.ToString());
                }
            }
            return res.ToArray();
        }

        // https://msdn.microsoft.com/en-us/library/windows/desktop/aa394132(v=vs.85).aspx
        public static string[] GetHardDriveSerials2()
        {
            var res = new List<string>();
            var query = new ManagementObjectSearcher("SELECT * FROM Win32_NetworkAdapter ");

            foreach (ManagementObject o in query.Get())
            {
                var model = o["Model"];
                var typ = o["InterfaceType"];
                var caption = o["Caption"];
                var id = o["DeviceID"];
                if (id != null)
                {
                    res.Add(id.ToString());
                }
            }
            return res.ToArray();
        }
#endif

        // heuristic to pic 
        internal struct NetworkCardInfo
        {
            public int typ;
            public string guid;
            public string name;
            // heuristic based on what I saw
            public bool IsBluetooth()
            {
                var s = name.ToLowerInvariant();
                return s.Contains("bluetooth");
            }

            // the higher the priority, the more important the value
            // heuristic based on a guesswork
            /* types:
                Ethernet 802.3 (0)
                Token Ring 802.5 (1)
                Fiber Distributed Data Interface (FDDI) (2)
                Wide Area Network (WAN) (3)
                LocalTalk (4)
                Ethernet using DIX header format (5)
                ARCNET (6)
                ARCNET (878.2) (7)
                ATM (8)
                Wireless (9)
                Infrared Wireless (10)
                Bpc (11)
                CoWan (12)
                1394 (13)
            */
            public int TypePriority()
            {
                // 10 - ethernet that is not bluetooth
                // 9 - wireless
                // 8 - ethernet that is bluetooth
                // 7 - everything else
                // 6 - -1
                if (typ == -1)
                {
                    return 6;
                }
                bool isEthernet = (typ == 0) || (typ == 5);
                bool isBluetooth = IsBluetooth();
                if (isEthernet && !isBluetooth)
                {
                    return 10;
                }
                if (isEthernet && isBluetooth)
                {
                    return 8;
                }
                if (typ == 9)
                {
                    return 9;
                }
                return 7;
            }
        }

        // return true if c1 is more important than c2
        public static bool NetworkAdapterGt(NetworkCardInfo c1, NetworkCardInfo c2)
        {
            return c1.TypePriority() > c2.TypePriority();
        }

        // https://msdn.microsoft.com/en-us/library/aa394216(v=vs.85).aspx
        // available since Vista
        // Return a guid of a network card. If there is more than one network card,
        // try to pick the best one.
        // This value is meant as a unique id of the computer
        public static string GetNetworkCardId()
        {
            NetworkCardInfo card;
            card.typ = -1;
            card.name = "";
            card.guid = "";

            var query = new ManagementObjectSearcher("SELECT * FROM Win32_NetworkAdapter ");

            foreach (ManagementObject o in query.Get())
            {
                NetworkCardInfo card2;
#if DEBUG
                var typ = o["AdapterType"];
                var id = o["DeviceID"];
                var mac = o["MACAddress"];
                var phys = o["PhysicalAdapter"];
                var pnpid = o["PNPDeviceID"];
#endif
                var guid = o["GUID"];
                var name = o["Name"];
                var caption = o["Caption"];

                if (guid == null)
                {
                    continue;
                }
                UInt16? typid = o["AdapterTypeID"] as UInt16?;
                if (typid == null)
                {
                    card2.typ = 20; // bogus value different than documented types
                }
                else
                {
                    card2.typ = (int)typid;
                }
                card2.guid = guid.ToString();

                card2.name = "";
                if (name != null)
                {
                    card2.name = name.ToString();
                } else if (caption != null)
                {
                    card2.name = caption.ToString();
                }

                // remember this card if more important than previous
                if (NetworkAdapterGt(card2, card))
                {
                    card = card2;
                }
            }
            return card.guid;
        }

        // basic information about OS and user
        public struct ComputerInfo
        {
            public string UserName;
            public string OsVersion;
            public string MachineName;
            public string NetworkCardId;
        }

        // consider returning more info from:
        // Win32_OperatingSystem  https://msdn.microsoft.com/en-us/library/aa394239(v=vs.85).aspx
        // Win32_ComputerSystem  https://msdn.microsoft.com/en-us/library/aa394102(v=vs.85).aspx
        // Win32_Processor https://msdn.microsoft.com/en-us/library/aa394373(VS.85).aspx
        // Win32_MotherboardDevice https://msdn.microsoft.com/en-us/library/aa394204(v=vs.85).aspx
        public static ComputerInfo GetComputerInfo()
        {
            ComputerInfo i;
            i.NetworkCardId = GetNetworkCardId();
            i.UserName = Environment.UserName;
            i.OsVersion = Environment.OSVersion.Version.ToString();
            i.MachineName = Environment.MachineName;
            return i;
        }
    }
}
