package com.ritchie;

import com.ritchie.formula.Formula;

public class Main {

  public static void main(String[] args) {
    String input1 = System.getenv("SAMPLE_TEXT");
    String input2 = System.getenv("SAMPLE_LIST");
    boolean input3 = Boolean.parseBoolean(System.getenv("SAMPLE_BOOL"));
    Formula formula = new Formula(input1, input2, input3);
    System.out.println(formula.Run());
  }
}
