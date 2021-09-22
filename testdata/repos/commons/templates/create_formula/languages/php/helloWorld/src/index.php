<?php
require __DIR__ . '/vendor/autoload.php';
include 'formula/formula.php';

$input1 = getenv('RIT_INPUT_TEXT');
$input2 = getenv('RIT_INPUT_BOOLEAN');
$input3 = getenv('RIT_INPUT_LIST');
$input4 = getenv('RIT_INPUT_PASSWORD');
Run($input1, $input2, $input3, $input4);
?>
