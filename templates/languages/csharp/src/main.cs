using System;
using formula;

namespace main
{
    static class main
    {
        static void Main(string[] args)
        {
            string TEXT     = Environment.GetEnvironmentVariable("SAMPLE_TEXT");
            string TLIST    = Environment.GetEnvironmentVariable("SAMPLE_LIST");
            string TBOOL    = Environment.GetEnvironmentVariable("SAMPLE_BOOL");
            new formula.Hello(TEXT, TLIST, TBOOL);
        }
}
}
