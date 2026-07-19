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
