using System;
using main;

namespace formula
{
    public class Hello
    {
        public string Input1;
        public string Input2;
        public string Input3;
        public Hello(string sptx, string sptl, string sptb)
        {
            this.Input1 = sptx;
            this.Input2 = sptl;
            this.Input3 = sptb;
            Console.WriteLine("Hello World!");
            Console.WriteLine("You receive " + Input1 + " in text. ");
            Console.WriteLine("You receive " + Input2 + " in list. ");
            Console.WriteLine("You receive " + Input3 + " in boolean. ");
        }

    }
}
