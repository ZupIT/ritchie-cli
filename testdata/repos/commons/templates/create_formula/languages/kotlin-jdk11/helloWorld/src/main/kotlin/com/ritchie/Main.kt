package com.ritchie

import com.ritchie.formula.Formula

object Main {
    @JvmStatic
    fun main(args: Array<String?>?) {
        val inputText: String = System.getenv("RIT_INPUT_TEXT")
        val inputBoolean: Boolean = (System.getenv("RIT_INPUT_BOOLEAN").toBoolean())
        val inputList: String = System.getenv("RIT_INPUT_LIST")
        val inputPassword: String = System.getenv("RIT_INPUT_PASSWORD")
        val formula = Formula(inputText, inputBoolean, inputList, inputPassword)
        formula.Run()
    }
}
