using System;
using System.Collections.Generic;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Runtime.Serialization;

namespace DbHero
{ 
    [DataContract]
    public class Settings
    {
        [DataMember]
        public Point WindowPosition;

        [DataMember]
        public Size WindowSize;
    }
}
