using System;

namespace formula
{
    public class Hello
    {
        public string input1;
        public string input2;
        public string input3;
        public string input4;
        public Hello(string intext, string inlist, string inbool, string insecret)
        {
            this.input1 = intext;
            this.input2 = inlist;
            this.input3 = inbool;
            this.input4 = insecret;

            Console.WriteLine("Hello World!");
            Console.ForegroundColor = ConsoleColor.Green;
            Console.WriteLine($"My name is {input1}.");
            Console.ForegroundColor = ConsoleColor.Blue;
            if (input3 == "true")   {
                Console.WriteLine("I've already created formulas using Ritchie.");
            } else  {
                Console.WriteLine("I'm excited in creating new formulas using Ritchie.");
            }
            Console.ForegroundColor = ConsoleColor.Yellow;
            Console.WriteLine($"Today, I want to automate {input2}.");
            Console.ForegroundColor = ConsoleColor.Cyan;
            Console.WriteLine($"My secret is {input4}.");
            Console.ResetColor();
        }

    }
}
