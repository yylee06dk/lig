# lig

## Motivation

While reading [Crafting Interpreters](https://craftinginterpreters.com/), I felt like I was just blindly copy and pasting code.
So, using Golang (which has GC, which is crucial for jlox implementation), I decided to make a language of mine,
composed with my choices and philosophy, strongly inspired by jlox.

## Structure

Since, I want to feel the 'building up' process of a language, I will slice it up into stages. If I were to really implement a language of mine,
I wouldn't start by defining the perfect EBNF that expresses the whole language. So, I'm starting small, then going to add more and more features to the bare bone.
And, as a result I desire to end up with a somewhat smaller, compact version of lox with slightly different selections made within the process
Also, this is the literal "first" project written in Golang, and the whole README will have my impressions on Golang, and the difference of transplanting the thought process
made in Java to Go, which is quite different from the perspective of OOP.

### Stage 1. A calculator

numbers, binary ops(+, -, *, /) is only we got for now. Also only REPLs work for now.  

Notes on what I learned/felt on each stage of the interpreter

###### Scanner
 Exploring the Go's error reporting system, it is very nice. Only logging at the top level. This makes error reporting modularized in each function~package level, and reporting
to the user is only done at the top level. A nice usage of the stack-like structure.  
Also, found out OOP naturally allows static global variables and I never knew it was a previlege. Without such object-ness, I should always give the source string(code) and
the current position where I'm reading(in byte-level) to every functions in the scanner package. Of course, this can be fixed through structs and method recievers, but I wanted
to try coding without them.

###### Parser
 Found out how the jlox's parser avoids the problem that normally occurs with right-recursive grammar. (ex. 1 - 2 - 3)  
 Also, found out when coding with techniques like recursive descent(for parser), we should carefully place incrementing codes. Best practice seems to
be always trying the best to keep in a single level. In the code given, incrementing happens at the most lowest level(literal)(Also, this is the part
where an actual expression is made)
 Hmm. The parser I coded from the book in chapter 6 doesn't seem to report issues of parsing 1 + 2 8 3. It just interprets 1 + 2 then ends it.
So, I added a check of EOF before terminating the parse process, so when such issue happens, error is reported

###### Interpreter
 With type switches, I could easily make a functional, ML-like approach. It just felt like pattern matching since using `v := expr.(type)` just
casted the desired subtype automatically.
 Also, for future convenience, decided to return any type values. Found out go has a nice tools to deal with any types such as type switches and
type assertions. 

### Stage 1-1. Adding more data we can use

Currently we only have numbers. Let's add true, false, and strings! Also, let's add operations we can use to manipulate them.

The key point in this section is how we will implement strings. Implementing strings is strongly affected by how the host language deals with them.
First, let's learn how `Go` deals with it.

##### String implementations in several languages

###### Go
In the Go's builtin package, it has type string:
> string is the set of all strings of 8-bit bytes, conventionally but not necessarily representing UTF-8-encoded text.
A string may be empty, but not nil. Values of string type are immutable

So this is the string literal. Also, it has raw string literals, but we won't dive deep into those.
Go's strings are stored in the read-only section. So, if we want mutable strings we type cast to byte slices or rune slices(if needed)
which forces allocation to heap.

###### python
In python, strings are arrays right away. It is quite similar to C, but with no pointers & types. Also, it is immutable by default.

###### C
Every strings are char* with null termination. But string literals like `char *a = "Hello";`, is allocated in the read-only section which
disallows mutability. But allocation in the stack or heap makes it mutable.

###### Java
It uses string pool for efficiency. Since string literals are strictly immutable, it can reuse the same memory address again without worrying about
side-effects. And for mutable strings, different classes such as StringBuffer or StringBuilder must be used.

Viewing such implementations show me that a type system and an ADT is needed to fully implement mutable strings via arrays(since array
needs to pre-allocate memory with fixed size, type system, or at least a size of character should be given)

###### Summary
I will implement strings with Golang's strings for now. But later on, maybe in another language, I might try a type system and a C-like string

#### Implementing!

Notes on what I learned/felt while coding LIG

###### Scanner
Adding new tokens to the scanner was relatively easy. Tried out testing, and unit testing was easy to try out, but I immediately found
out that a great test table is where the real value comes from. It's hard to make a good test table for every scanning rule and edge cases

###### Parser
Before starting with the parser, we should definitely take a look at the grammar. When reading the book, the first place I bumped into
was the abstract grammar using BNF. I understand that we parse by "less-bonding" or lower precedence to "strongly-bonding" or higher precedence
because that's how it should work. But the real question was, where did the precedence first come from. As a person who's first language
is C, the given precedence seems intuitive. But I'm definitely sure, that it wasn't **that** intuitive back in the 80s. So I decided to
build up my own precedence.

###### Keypoints of my precedence :
Since we're currently building a dynamic language, let's reduce runtime errors
And, make it seem intuitive as possible. (The judge of intuitive-ness is me!)

Operators we have: +, -, *, /, ++, !, != , =, ==, >, >=, <, <=, &&, ||
Of course some operators like =, &&, || are not implemented in this stage, so we will skip for now

first let's look at some intuitive ones.
+, - have same precedence and it is lower than *, / (basic math)

Also, for >, >=, <, <=(the comparison operators) we need to remark that
they do not chain (like 1 + 2 + 3). Chaining these operators will end up in a
quite weird state since calculating one side will return a boolean value and
boolean values are not comparable(we can make it so, but it wouldn't be a intuitive one)
So, let's first group them with a name `compops`, with `compops -> ">" | ">=" | "<" | "<="`

Now we need to actually decide which one has higher precedence
First, let's compare == with `compops`.

Look at `a == b <= c` where a,b,c are any-type values.
== or `compops` are binary operators who spit out a boolean value. == can have boolean values as operands, but
`compops` cannot. So, by our thought of reducing runtime errors, we should give `compops` higher precedence

Look at `a <= b + c` where a,b,c are any-type values.
Also, with the same logic, we see that `(+)` operator takes(also `(-)`) integer operands only.
But the comparison doesn't. So we should give `(+)` higher precedence.

Look at `!a + b` where a, b are any-type values.
This is where I thought about it alot. For intuitive grammar it is better to give unary operators higher precedence.
But, with reducing runtime errors we see that no one will try to write `{boolean} + {int}`.
Then, if we look at `-a + b`, it is much better to think of it as `(-a) + b`
So, despite the risk of `!` operator, we place it below plus/minus operators

The same applies to the mult/div operators

The current precedence is like so:
`Expression -> Equality -> Comparison -> Term -> Factor -> Unary -> Primary` (Here Primary means true, false, literals, identifiers)

Now we need to find a place to squeeze in the `++` operator. (I did regret a bit adding this operator and not using overloading)

Looking at `a ++ b < c` shows that it is better to have higher precedence than comparison operators

So, now looking at `a ++ b + c`, it just didn't make sense. Such code will 100% spit out a runtime error since `lig` is a dynamic language, but
the host language, `Golang` is not. So, for my personal intuition, I decided it to have higher precedence than plus/minus operators.

So, we have the final precedence, so let's right the grammar rules out.

```
expression -> equality
eqaulity -> comparison ("==" | "!=") equality | comparison
comparison -> term compops term | term   //  Where compops -> ">" | ">=" | "<" | "<="
term -> concat ("+" | "-") term | concat
concat -> factor "++" concat | factor // This does feel really weird
factor -> unary ("*" | "/") factor | unary
unary -> ("!" | "-") primary | primary
primary -> "true" | "false" | Number | String | Identifier
```

We see that all rules are right-recursive (except for comparison where we cannot chain the operators) since we are going to use a
recursive descent parser which needs to parse from the left. So, we should be able to parse the "left-part" then parse the rest.
But, using left-recursive grammar, we cannot deal with parsing the "left-part". So, right-recursion is key in the grammar

Addressing the grammar was quite of a pain(Since I read the book once, I passed the real pain of amibiguity and more)
But, after the grammar was given, the implementing the parser was so easy.

###### Interpreter
Again, it was quite easy to deal with, but the go's type assertion was a bit of a pain, but it was doable, and still an enjoyable process.
Also, learned how Go deals with return statements especially mixed with switch statements. It seems safe and robust.

Learned what an **interface** really is... I first found the go's interface a bit confusing, especially with value receivers and pointer receivers. And
kept on trying to pass pointers to interfaces with the type `*datatypes.Expr` which is soooo wrong. This is a pointer to a interface. So, we cannot use polymorphism.
Rather, just `datatypes.Expr` alone means that every type that implements this interface. And using pointer receivers, we implement the interface with pointer types.
So, `datatypes.Expr` works as a big set of `*datatypes.Binary`, `*datatypes.Unary`, `*datatypes.Literal`. (Here, it is not weird to have the pointer since they are structs not interfaces)
With this understood, changing the code was so easier.
And now, I have a faster and more memory efficient interpreter(compared to the before one)

Decided to add Name field to tokens --> Why?
I tried to implement scanning/parsing identifiers. Then, identifiers has names not values. So, having a any-type value for a 100% string means meaningless code
needed to be written every time we deal with identifiers(type assertions). So, I added a Name field of type string and a Variable expression with Name field

Also, following from above, we see that we need to define keywords somewhere. I decided to define it in the datatypes package, since
- Identifier Expression needs to see it to know if itself is a variable/keyword
- We can import it from the interpreter to see if its a keyword.
> Note that we are not checking this from the parser since parser is only about grammar.

### Stage 1-2
Adding some convenience features to the language. (File as input, better error printing)
And, between that I found out why unsafe-features are needed, and why the book decided to have lexemes inside tokens.
The first one is because, sometimes the compiler's safe checks are a limit for efficiency and speed. When human can 100% ensure some properties of the data, we can there are cases
where we can perform better in unsafe ways. For example, to add files as input for the interpreter, I changed every string-types for source code to []byte types. And, changing a
file to a []byte type means allocation in the heap. If the allocation is very large, this will become a very expensive function call. So, in such interpreters, we know that we will
only read the source code. So, we can use unsafe function calls such as unsafe.Slice that immediately returns the existing memory address without allocation(So, the memory of the
source file)
For the second one, it was for error printing. I feel like I'm (once again) losing to the book since, I got rid of it because lexeme seemed to be used only for identifiers, and
we can deal the identifier's names with any-typed field, `Value`. But, yepp to have better error printing, we need to deal with lexemes

Not exactly the fun part for user convenience. I feel like I'm digging more into single functions to see if it can satisfy what I'm looking for. Since, I haven't dealt with
user convenience in coding, (I only coded for CS projects that normally try to conceptualize what we learned) this was quite of a hard part.
Feels like if I started with a better structure, it would've been a little easier. Starts feel like the code is turning into spaghetti. I'm looking for more hacks to just squeeze
in one functionality(like skipping comments) and more. I guess this is the way to "good architecture". Modularizing code and especially determinating return types.
Like I see many cases that I can deal with some explicit return value, but this makes the code worse. In go, I could've tried something like adding a struct field or adding a
return value and more. But, this makes the need of total change in the code. So I chose to just return something weird and check if it's weird in the caller.

Finally dealed with all of the things I've aimed!!
The biggest targets were, error handling with lines, and File as input + debug mode with flags.
So happy to see things are working as intended... This definitely is the most loveable project in my life.

## First Milestone.
The code I've written feels like the first milestone of the project. I've nearly got every functionality jlox provides until Ch.7. It results as a nicely written 1000~ lines of 
code, which follows idiomatic Go-like designs. 

Now I've got a decision to make. I could march on to variables and scope or do some more detail-related things like ternary operators(`? :`), or multi line comments.
But, since I'm currently thinking of this project as a subset of jlox, and I'm trying to be more concise in such manner, so I'll just let such details be implemented in a 
later language and march on for now.

### Stage 2
Now, it is time to move on. The next stage is adding statements. Like many procedural languages, we have statements which do not interpret to values. But rather have side-effects
In this stage we aim to add statements. I found out if my grammar looks like `code -> (statement)* EOF` then, we cannot use the current 
form of checking the result value of just printing out the return value since statements don't return values. This means we need to 
explicitly follow the values given by each statement then return the last-value we saw or something like that. So, now I understand
why we had to implement a `printStmt` in the book. 

###### Scanner
I only need to add datatype Semicolon, since in this phase we don't even take a look at the grammars.

###### Parser
Need a structural re-design.
Token list is parsed into list of trees not a single tree. A single statement is parses to a single tree.

Finished with statements. For interpreter too. It rather feels a bit odd.. After the minimal structure is given, adding single
features are sooo easy. I think this is the power of goo modularizing. Adding one function means adding one functionality. This is peek 
coding.. Thanks to the book, I could start with a nicely structured code.

### Stage 2.1 
Now adding environments and variables. The book uses grammar like `code -> declaration* EOF` and `declaration -> varDecl | statement`
This avoids variable declarations to be only available at top-level. This is due to the existence of the blockStmt. Like I currently don't understand why this is so.
But, after implementing it myself, I might change!

###### Interpreter
I understand why the book had to deal with environments by making a instance of the interpreter. I first thought it is doable
just by passing environment as an argument, but I wasw wrong. This is possible when we decide to make the interpret function 
to a whole big function. (Easily seen when we have pattern matching in ML-like languages) But, in go we need to divide logics 
into functions because we can, but is not recommended to use multiple nested switches/or if statements.(It looks much worse
than nested pattern matching)
 So, using a somewhat procedural language, I will follow the books decision.
And the main point of the interpreter struct is that we only make a new one when we need to deal with a code block.
A code block can be thought as the whole code, a block normally defined with curly braces.

This stage was a bit buggy due to my REPL implementation, but solved it with a global variable that holds the starting empty environment

### Stage2.2 - Adding assignment
I first had to decide if assignment is a expression or a statement. As I said above, the biggest difference between the two is if it 
interprets to a value. If I choose it to be a statement I can't do things like `a=b=1` where it means `a = (b=1)` and `b=1` interprets to 
value `1` and so on. (Right recursiveness is assumed) Choosing it as an expression is a C-like approach but it is error-prone. 
such as `if (x=3)` will automatically proceed tothe then-branch in C since assignment results in a value `3` which is true in C.
So, I decided to make it a statement.

Again, adding was relatively simple, just under a hour work of just reading the code, choosing precedence then done!

### Stage 2.3 - Code blocks
Now having working environments, we need to change the update the environment as we go. 

### Stage 3 - Control Flow
Now, let's add if-then-else statements and while & for loops. 
First, let's choose the keywords. We won't use a `then` keyword since we have statements. If every thing was an expression, everything
should result(interpret to) a vaule. So, only having if-statements is not allowed in such languages. Because, if the condition for the 
statement is not matched, we don't know what the value of the expression is. But, this results in very ugly code when we don't have 
pattern matching or other types of control flow available. So, for lig, I will only use `if` & `else` keywords along with `for` and `while`

###### Parser
Grammar :
```
Stmt -> IfStmt | WhileStmt | ForStmt | ~~~
IfStmt -> "if" "(" expression ")" statement ("else" statement)?
WhileStmt -> "while" "(" expression ")" statement
ForStmt -> "for "(" initializer? ";" expression? ";" updater? ")" statement
```

Currently finished with if-else statements. Again felt so easy. With grammar and a uniform architecture to the code, I can just easily add stuff.

Now, let's take a look at the for statement's grammar
First, what should we allow for the initializer?
- variable declaration
- assignment statements
- expressions --> These do not have any real side-effects, but I just added them to easily code it. (assignOrExprStatement exactly covers
the two statements above)

Secondly, how should the updater look like?
- assignment statements
- expressions --> The same reason

With this, we can easily write the grammar for `ForStmt`
```
initializer -> varDecl | assignStmt | expression 
updater -> assignStmt | expression
```

And, while implementing the parser for the for statement, I had to choose if I will desugar the statement into a while
statement. I decided to not, since in the book I already tried out desugaring, and the advantages of desugaring comes from
optimization code written for the `while` loop can be used again for `for` loops. But in such simple language(as said in the 
book), it doesn't have much advantages both way.

###### Interpreter
While implementing the interpreteer for while statements, I faced another subtle detail about lexical scopes. First, I tried to implement 
it simply just like the book. But, soon found out if the `statement` of the while-grammar was a block statement, re-executing the block 
statement over and over again will not conserve the environment. At that time, I was thinking that the environment of the iterated part
of the `while-statement` should have a single environment that accumulates data as it goes. 
But, I was wrong. `C` also, just kept on iterating block statements. That's becuase both semantically doesn't matter. 
If I were to make a declaration in the iterated block, that doesn't need to be accumulated since, the declaration will happen again
returning to the same ground.
> Or is it?
Then, a question arises. Actually when we code, we normally don't have variable declarations inside loops. If I were to look at the side
effects happening by looping, why would you declare it again returning to the same ground? So, we see that `var a = 1;` is not a big
issue if we view from the semantic perspective. But, if wee see `var a = a + 1;`, now what? Should we accumulate? What does that even mean?
To understand this issue, we need to address some things. 
- When are variables declared to the scope? 
In lox, when `var a = 1;` is seen, `a` is introduced to the environment with the initializer. So, the right side is evaluated with the
original(`a` not yet introduced) environment. But for example, `C` is different. `C` introduces the variable without checking the
initializer. So, when `int a = 1; int a = a + 1;` is seen, this is undefined behavior since `a` in the right side of the second statement is not yet initiated, and will normally contain garbage value. On this topic, many languages chose a `C`-like approach but not always
occurring undefined behavior.(Javascript throws an ReferenceError) Or, language like `Rust` works like lox. So, the second statement will assign value 2 to `a`. And, I personally think lox & rust's aprroach is a bit better, and my code currently works this way, so I'll just let it be.
--> Exact same thing is expressed in Ch.11, and lox chooses to make it an error. (I personally do think initializer should be interpreted in the original envrionment)
- What is the interpreting environment of the `condition Expression`? 
We need to interpret the condition everytime to know if we will iterate once more. So, it is very important to select this.
For the first iteration, we will not use/can not approach the inner stuff(the code inside the while block), so we should be able to interpret right away from outer resources. Then, we can think of two cases. In the inner block we re-declare(shadowing) the resouces used in the condition expression, so that's how value updates are done. This seems a bit off. This intuitively seems like it will cause alot of unintended errors/outputs, so we should always interpret the condition with the outer environment(inner code can update the outer env by assigning values(`a = a+1;`))
- Should the environment of inner block be accumulated?
So, if we set the environment when evaluating the expression should be kept as the outer environment, we see that accumulation is not needed in a semantical perspective. 

With this, we are finally ready to implement the interpreting mechanism of while

Now, we encounter a very similar argument when we try to interpret for statements. And, I once more applied the same logic while interpreting it.
The init statement is done in the outer environment which seems quite intuitive to me(you shouldn't do it inside the iterated part right?). Then, with the condition we also interpret it in the outer environment due to the reason explained above. And, the incremental statement can be done separately with the block because there is no reason to modify the inner scope variables.
 I do think actually the updating statement can be related with the inner scope, but this would be error prone and make the code harder to read.
