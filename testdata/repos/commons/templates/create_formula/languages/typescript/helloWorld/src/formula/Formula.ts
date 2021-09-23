import * as chalk from 'chalk'

function run(inputText: string, inputBoolean: boolean, inputList: string, inputPassword: string) {
    console.log('Hello World!')

    console.log(chalk.green(`My name is ${inputText}.`))

    if (inputBoolean) {
        console.log(chalk.blue(`I've already created formulas using Ritchie.`))
    } else {
        console.log(chalk.red(`I'm excited in creating new formulas using Ritchie.`))
    }

    console.log(chalk.yellow(`Today, I want to automate ${inputList}.`))

    console.log(chalk.cyan(`My secret is ${inputPassword}.`))
}

const Formula = run
export default Formula
