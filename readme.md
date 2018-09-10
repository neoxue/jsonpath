#jsonpath

-   requirement:
    -   golang 1.9
-   operators: TODO
    
    
*   source url：https://github.com/oliveagle
    TODO finish those uncompleted 
    
supported syntax:
- / 	$ 	the root object/element
- . 	@ 	the current object/element
- / 	. or [] 	child operator
- // 	.. 	recursive descent. JSONPath borrows this syntax from E4X.
- * 	* 	wildcard. All objects/elements regardless their names.
- [] 	[] 	subscript operator. XPath uses it to iterate over element collections and for predicates. In Javascript and JSON it is the native array operator.
- | 	[,] 	Union operator in XPath results in a combination of node sets. JSONPath allows alternate names or array indices as a set.
- n/a 	[start:end:step] 	array slice operator borrowed from ES4.
- [] 	?() 	applies a filter (script) expression.
unsupported syntax:
n/a 	() 	script expression, using the underlying script engine.

filter script expression only supports finit expression like:
@aaa op other
@.number >= $.number   // means filters the current node which number greater than root number 




