package com.ritchie.formula;

import org.junit.Test;
import static org.junit.Assert.*;

public class HelloTest {

    @Test
    public void run() {
        Hello hello = new Hello("Hello", "World", true);
        String excpeted = "Hello World!\n" +
                "You receive Hello in text.\n" +
                "You receive World in list.\n" +
                "You receive true in boolean.\n";
        assertEquals(excpeted, hello.Run());
    }
}