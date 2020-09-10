<?php
require __DIR__ . '/vendor/autoload.php';
include 'formula/formula.php';

$input1 = getenv('SAMPLE_TEXT');
$input2 = getenv('SAMPLE_LIST');
$input3 = getenv('SAMPLE_BOOL');
Run($input1, $input2, $input3);
?>
