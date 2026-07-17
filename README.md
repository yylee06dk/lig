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

