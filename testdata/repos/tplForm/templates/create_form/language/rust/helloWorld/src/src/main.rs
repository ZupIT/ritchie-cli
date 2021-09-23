mod formula;

use std::env;

fn main() {
    let input_text = string_from_env("RIT_INPUT_TEXT");

    let input_bool = bool_from_env("RIT_INPUT_BOOLEAN");

    let input_list = string_from_env("RIT_INPUT_LIST");
    let input_password = string_from_env("RIT_INPUT_PASSWORD");

    formula::run(input_text, input_bool, input_list, input_password);
}

fn string_from_env(key: &str) -> String {
    match env::var(key) {
        Ok(val) => val,
        Err(_) => "none".to_string(),
    }
}

fn bool_from_env(key: &str) -> bool {
    match env::var(key) {
        Ok(val) => val.parse().unwrap_or(false),
        Err(_) => false,
    }
}
