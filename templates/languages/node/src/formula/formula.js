const clc = require("cli-color")


function Run(input1, input2, input3) {
    console.log("Hello World!")
    console.log(clc.green("You receive "+ input1 +" in text."));
    console.log(clc.red("You receive "+ input2 +" in list."));
    console.log(clc.yellow("You receive "+ input3 +" in boolean."));
}

const formula = Run
module.exports = formula
