package com.ritchie.formula;

public class Formula {

  private String input1;
  private String input2;
  private boolean input3;

  public String Run() {
    return String.format(
        "Hello World!\n"
            + "You receive %s in text.\n"
            + "You receive %s in list.\n"
            + "You receive %s in boolean.\n",
        input1, input2, input3);
  }

  public Formula(String input1, String input2, boolean input3) {
    this.input1 = input1;
    this.input2 = input2;
    this.input3 = input3;
  }

  public String getInput1() {
    return input1;
  }

  public void setInput1(String input1) {
    this.input1 = input1;
  }

  public String getInput2() {
    return input2;
  }

  public void setInput2(String input2) {
    this.input2 = input2;
  }

  public boolean isInput3() {
    return input3;
  }

  public void setInput3(boolean input3) {
    this.input3 = input3;
  }
}
