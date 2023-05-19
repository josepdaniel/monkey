# Lexing


We lex/tokenize the input to represent source code in a form that is easier to work with.

```mermaid
graph LR
    A("Source Code\n`let x = 5`") -- Lexing --> B("Tokens\n[LET, IDENT.x, EQ, INT.5]") -- Parsing --> C("AST\nASSIGN(ident='x', val=int(5), in=nil)")
```


The input to the lexer is the source code, the output is a stream of tokens. A token is an interpreted atom of source code (it has converted the raw text to something meaningful).