package com.ritchie;

import com.ritchie.formula.Formula;

public class Main {

  public static void main(String[] args) {

    String inputText = System.getenv("RIT_INPUT_TEXT");
    boolean inputBoolean = Boolean.parseBoolean(System.getenv("RIT_INPUT_BOOLEAN"));
    String inputList = System.getenv("RIT_INPUT_LIST");
    String inputPassword = System.getenv("RIT_INPUT_PASSWORD");

    Formula formula = new Formula(inputText, inputBoolean, inputList, inputPassword);
    formula.Run();
  }
}
