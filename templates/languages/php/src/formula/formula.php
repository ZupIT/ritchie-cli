<?php
use Codedungeon\PHPCliColors\Color;

function Run($input1, $input2, $input3) {
    echo "Hello World! \n";
    echo Color::GREEN, "You receive $input1 in text. \n";
    echo Color::RED, "You receive $input2 in list. \n";
    echo Color::BLUE, "You receive $input3 in boolean. \n";
}
?>
