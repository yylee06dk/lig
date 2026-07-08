# lig

## Motivation

While reading [Crafting Interpreters](https://craftinginterpreters.com/), I felt like I was just blindly copy and pasting code.
So, using Golang (which has GC, which is crucial for jlox implementation), I decided to make a language of mine, composed with my choices and philosophy, strongly inspired by jlox.

## Structure

Since, I want to feel the 'building up' process of a language, I will slice it up into stages. If I were to really implement a language of mine,
I wouldn't start by defining the perfect EBNF that expresses the whole language. So, I'm starting small, then going to add more and more features to the bare bone.
Also, this is the literal "first" project written in Golang, and the whole README will have my impressions on Golang, and the difference of transplanting the thought process
made in Java to Go, which is quite different from the perspective of OOP.

### Stage 1. A calculator

numbers, binary ops(+, -, *, /) is only we got for now. Also only REPLs work for now.  

Notes on what I learned/felt on each stage of the interpreter

Scanner : 
Exploring the Go's error reporting system, it is very nice. Only logging at the top level. This makes error reporting modularized in each function~package level, and reporting
to the user is only done at the top level. A nice usage of the stack-like structure.  
Also, found out OOP naturally allows static global variables and I never knew it was a previlege. Without such object-ness, I should always give the source string(code) and
the current position where I'm reading(in byte-level) to every functions in the scanner package. Of course, this can be fixed through structs and method recievers, but I wanted
to try coding without them.

Parser :
Found out how the jlox's parser avoids the problem that normally occurs with right-recursive grammar. (ex. 1 - 2 - 3)

