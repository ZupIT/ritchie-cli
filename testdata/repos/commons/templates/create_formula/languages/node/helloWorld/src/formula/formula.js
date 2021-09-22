const clc = require("cli-color")

function Run(inputText, inputBoolean, inputList, inputPassword) {
    console.log("Hello World!")
    console.log(clc.green("My name is "+ inputText +"."))
    if(inputBoolean){
        console.log(clc.blue("I've already created formulas using Ritchie."))
    } else {
        console.log(clc.red("I'm excited in creating new formulas using Ritchie."))
    }
    console.log(clc.yellow("Today, I want to automate "+ inputList +"."))
    console.log(clc.cyan("My secret is " + inputPassword +"."))
}

const formula = Run
module.exports = formula
