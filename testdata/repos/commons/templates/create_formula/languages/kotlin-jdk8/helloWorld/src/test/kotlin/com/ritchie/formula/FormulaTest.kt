package com.ritchie.formula

import org.junit.After
import org.junit.Before
import org.junit.Test
import java.io.ByteArrayOutputStream
import java.io.PrintStream
import kotlin.test.assertEquals

class FormulaTest {
    private val outContent: ByteArrayOutputStream = ByteArrayOutputStream()
    private val originalOut: PrintStream = System.out
    @Before
    fun setUpStreams() {
        System.setOut(PrintStream(outContent))
    }

    @After
    fun restoreStreams() {
        System.setOut(originalOut)
    }

    @Test
    fun runTrueInput() {
        Formula("Hello", true, "world", "pass").Run()
        assertEquals("Hello World!" +
                "My name is Hello." +
                "I've already created formulas using Ritchie." +
                "Today, I want to automate world." +
                "My secret is pass.", outContent.toString().replace(Regex("\\r|\\n"), "").replace(Regex("\u001B\\[[;\\d]*m"), ""))
    }

    @Test
    fun runFalseInput() {
        Formula("Hello", false, "world", "pass").Run()
        assertEquals("Hello World!" +
                "My name is Hello." +
                "I'm excited in creating new formulas using Ritchie." +
                "Today, I want to automate world." +
                "My secret is pass.", outContent.toString().replace(Regex("\\r|\\n"), "").replace(Regex("\u001B\\[[;\\d]*m"), ""))
    }

    @Test
    fun runNoSecretsInput() {
        Formula("Hello", false, "world", "").Run()
        assertEquals("Hello World!" +
                "My name is Hello." +
                "I'm excited in creating new formulas using Ritchie." +
                "Today, I want to automate world." +
                "My secret is .", outContent.toString().replace(Regex("\\r|\\n"), "").replace(Regex("\u001B\\[[;\\d]*m"), ""))
    }
}
