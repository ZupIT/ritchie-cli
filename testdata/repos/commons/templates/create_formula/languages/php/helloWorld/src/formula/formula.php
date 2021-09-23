<?php
use Codedungeon\PHPCliColors\Color;

function Run($input1, $input2, $input3, $input4) {
    echo "Hello World!\n";
    echo Color::GREEN, "My name is $input1.\n";

    if ($input2 === "true") {
        echo Color::BLUE, "I've already created formulas using Ritchie.\n";
    } else {
        echo Color::RED, "I'm excited in creating new formulas using Ritchie.\n";
    }
    echo Color::YELLOW, "Today, I want to automate $input3.\n";
    echo Color::CYAN, "My secret is $input4.\n";
}
?>
