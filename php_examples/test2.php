<?php

$arr = array(
    "test",
    "value"
);

echo json_encode($arr, 1);
$arr = array(
    1 => "test"
);

echo $json = json_encode($arr, 1);
print_r(json_decode($json, 1));
