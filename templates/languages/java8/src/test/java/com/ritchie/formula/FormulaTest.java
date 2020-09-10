package com.ritchie.formula;

import static org.junit.Assert.*;

import org.junit.Test;

public class FormulaTest {

  @Test
  public void run() {
    Formula formula = new Formula("Hello", "World", true);
    String excpeted =
        "Hello World!\n"
            + "You receive Hello in text.\n"
            + "You receive World in list.\n"
            + "You receive true in boolean.\n";
    assertEquals(excpeted, formula.Run());
  }
}
