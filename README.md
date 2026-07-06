# lig

## Motivation

While reading [Crafting Interpreters](https://craftinginterpreters.com/), I felt like I was just blindly copy and pasting code.
So, using Golang (which has GC, which is crucial for jlox implementation), I decided to make a language of mine, composed with my choices and philosophy, strongly inspired by jlox.

## Structure

Since, I want to feel the 'building up' process of a language, I will slice it up into stages. If I were to really implement a language of mine,
I wouldn't start by defining the perfect EBNF that expresses the whole language. So, I'm starting small, then going to add more and more features to the bare bone.

### Stage 1. A calculator

numbers, binary ops(+, -, *, /), and grouping is only we got for now. Also only REPLs work for now.
