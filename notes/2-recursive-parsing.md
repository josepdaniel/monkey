# Parsing

Parsers are recursive - parsing an expression like 

```
let x = 3
let y = x + 2
```

implies that both sides of an addition expression are themselves expressions.

So we can define a simple parser for a simple language with two main functions - *assign* and *add*:

```
EXPRESSION := ADDITION | IDENT | INT
INT := [0-9]
IDENT := [a-zA-Z]
ADDITION := EXPRESSION '+' EXPRESSION
ASSIGNMENT := 'let' IDENT '=' EXPRESSION
```

`ADDITION` is defined in terms of `EXPRESSION`, which is in turn defined in terms of `ADDITION`.

## Infinite recursion

- What happens if we try to parse an input like '3'? - We will drop into the 'addition' branch of the parser, because it *could* be an addition.
- What is an addition? It is an expression made up of `EXPRESSION '+' EXPRESSSION` - so we will have to try and parse the left-hand-side expression. 
- What is an expression? Well it could be an addition, so let's try parsing that. 
- It's borked.

## Re-ordering doesn't work

```
EXPRESSION := IDENT | INT | ADDITION
```

- Moving the 'addition' parser to the lowest priority fixes our infinite recursion problem. But then we will never be able to parse an expression like `1+1`. 
- The parser will first drop into the `int` branch and parse the first `1`. The leftover input `+1` doesn't match anything in our parser, so it will terminate there. 

## Eliminate the left-recursive path

- The infinite recursion arises because the 'left recursive' path has no base case - it will continue recursing to infinity.
- We can force it to have a base case by splitting the recursive action into two steps

```
EXPRESSION := START END
START := IDENT | INT
END := '+' EXPRESSION | '-' EXPRESSION | NOTHING
```

- An expression is now made up of a 'start' and an 'end'. The start has our simple non-recursive (terminal) parsers
- The end has our recursive (non-terminal) parsers (now also including subtraction). The end can also be 'Nothing'.
- The goal is to always have the recursive parsers match some token first (`+` or `-` in this case) before doing any recursion. That way a base case will always be hit. 