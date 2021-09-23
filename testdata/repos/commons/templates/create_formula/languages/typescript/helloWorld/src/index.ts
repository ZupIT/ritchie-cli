import run from './formula/Formula'

const inputText: string = process.env.RIT_INPUT_TEXT
const inputBoolean: boolean = JSON.parse(process.env.RIT_INPUT_BOOLEAN.toLowerCase())
const inputList: string = process.env.RIT_INPUT_LIST
const inputPassword: string = process.env.RIT_INPUT_PASSWORD

run(inputText, inputBoolean, inputList, inputPassword)
