use colored::*;

pub fn run(input_text: String, input_bool: bool, input_list: String, input_password: String) {
    println!("Hello World!");
    println!("{}", format!("My name is {}.", input_text).green());

    if input_bool {
        println!("{}", "I've already created formulas using Ritchie.".red())
    } else {
        println!(
            "{}",
            "I'm excited in creating new formulas using Ritchie.".red()
        )
    }

    println!(
        "{}",
        format!("Today, I want to automate {}.", input_list).yellow()
    );
    println!("{}", format!("My secret is {}.", input_password).cyan());
}
