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


## Associativity
Our new grammar causes an issue with **associativity**. We want our expressions to evaluate from left to right (we are going to ignore operator precedence).

```
3 + 4 - 5

# Left associative (what we want)
(3 + 4) - 5  ==>  2

# Right associative (what our grammar gives us)
3 + (4 - 5)  ==>  4
```


We can further modify our grammar to maintain left-associativity:

```
EXPRESSION := START {END}
START := IDENT | INT
END := '+' START | '-' START
```

The `{curly braces}` indicate 0-or-more semantics. Using this grammar, our parser will produce this tree for the statement `a + b - c`:

```
( a + b ) - c 
```

which is what we want! Our parser can have the following structure:

```go
// expression := start {end}
func parseExpression(input string) (Ast, string){
    ast, rest = parseStart(input)
    ast, rest = parseEnd(rest, ast)
    retrurn ast, eRest
}

// start := ident | int
func parseStart(input string)  (Ast, string) {
    ast, rest = parseIdent(input) orElse parseInt(input)
    return ast, rest
}

// {end} := {'+' start | '-' start}
func parseEnd(input string, start ast) {
    if input.startswith("+") {
        ast, rest = ast.Add(start, parseStart(input[1:]))
        return parseEnd(rest, ast)
    }
    ...
    // recursive base case
    return input, start
}
```


## Parentheses

We will make operator precedence explicit by requiring parentheses around expressions that should be evaluated together. An expression can now be surrounded by braces. The 'enclosed expression' will live in the 'start' grammar. 

```
EXPRESSION := START {END}
START := ENCLOSED_EXPR | IDENT | INT
ENCLOSED_EXPR := '(' EXPRESSION ')'
END := '+' START | '-' START
```


## Operator precedence

https://www.engr.mun.ca/~theo/Misc/exp_parsing.htm#climbing

