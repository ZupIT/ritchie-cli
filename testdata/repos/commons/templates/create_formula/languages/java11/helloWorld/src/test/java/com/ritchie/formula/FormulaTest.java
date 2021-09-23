package com.ritchie.formula;

import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;

import static org.junit.Assert.assertEquals;

public class FormulaTest {

  private final ByteArrayOutputStream outContent = new ByteArrayOutputStream();
  private final PrintStream originalOut = System.out;

  @Before
  public void setUpStreams() {
    System.setOut(new PrintStream(outContent));
  }

  @After
  public void restoreStreams() {
    System.setOut(originalOut);
  }

  @Test
  public void runTrueInput() {
    new Formula("Hello", true, "world", "pass").Run();

    assertEquals("Hello World!" +
            "My name is Hello." +
            "I've already created formulas using Ritchie." +
            "Today, I want to automate world." +
            "My secret is pass.", outContent.toString().replaceAll("\\r|\\n", "").replaceAll("\u001B\\[[;\\d]*m", ""));
  }

  @Test
  public void runFalseInput() {
    new Formula("Hello", false, "world", "pass").Run();

    assertEquals("Hello World!" +
            "My name is Hello." +
            "I'm excited in creating new formulas using Ritchie." +
            "Today, I want to automate world." +
            "My secret is pass.", outContent.toString().replaceAll("\\r|\\n", "").replaceAll("\u001B\\[[;\\d]*m", ""));
  }

  @Test
  public void runNoSecretsInput() {
    new Formula("Hello", false, "world", "").Run();

    assertEquals("Hello World!" +
            "My name is Hello." +
            "I'm excited in creating new formulas using Ritchie." +
            "Today, I want to automate world." +
            "My secret is .", outContent.toString().replaceAll("\\r|\\n", "").replaceAll("\u001B\\[[;\\d]*m", ""));
  }
}
