<?php
require_once __DIR__ . "/../vendor/autoload.php";

use Flow\JSONPath\JSONPath;
use Flow\JSONPath\JSONPathLexer;
use \Peekmo\JsonPath\JsonPath as PeekmoJsonPath;
test();

function test() {
        $json = '
            {
                "features": [{"name": "foo", "value": 1},{"name": "bar", "value": 2},{"name": "baz", "value": 1}]
            }
        ';
    $json = json_decode($json, 1);
    $results = (new JSONPath($json))->find('$[]');
    print_r($results);
    //$results = (new JSONPath($json))->find('$..features[?(@.value = 1)]');
    //print_r($results);
    //$results = (new JSONPath($json, 2))->find('$..features[1].name');
    //print_r($results);
}
