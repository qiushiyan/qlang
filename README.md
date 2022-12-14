Q Progrmaming Language
================

Q is a toy programming language with a mix of R and Python’s syntax. It
was written in Go and inspired by <https://interpreterbook.com/>.

## Usage

Go to [releases](https://github.com/qiushiyan/qlang/releases/tag/v0.0.1-lw) and download the corresponding binary.

Alternatively if you have Go installed, you can clone and cd into this repo and run

```
go run cmd/q/main.go
```

## Assignments

Both `=` and `<-` can be used for assignment. Variable names can contain
letters, numbers, and underscores, but must start with a letter.

``` q
x = 1
y <- x
x + y
```

    #> 2

There is also a `let` keyword for variable declaration.

``` q
let x <- 1
```

Q makes a difference between assignment and declartion in the context of
nested scopes. This is when we work inside functions, `if` and `for`
statements. See more details in the section on control structures below.

## Data structures

### Primitives

Primitive data structures include: numbers (integers and floats),
strings and booleans

``` q
1 + 1 + (10 * 2) / 4
```

    #> 7

``` q
"hello" + " " + "world"
```

    #> "hello world"

``` q
!false
```

    #> true

### Vectors

Both R and Python has 1-dimensional containers for storing a series of
values. Q offers the vector data structure as created by `[]`

``` q
[1, 2, 3]
```

    #> [1, 2, 3]

The `print()` helper shows the vector’s type as well as the number fo
elements.

``` q
print([1, 2, "hello"])
```

    #> Vector with 3 elements
    #> [1, 2, "hello"]

As in R, a vector is typed by its inner elements. Vectors containing
only numbers are numeric vectors, vectors with only string elements are
character vectors, and so on. A vector with mixed types is simply a base
`Vector` type, similar to a Python list. No type conversion is done
automatically, if the elements are heterogeneous the base type will be
used.

Vectors in Q have 1-based indexing: the first element starts at index 1,
not 0. Built-in functions for vectors include `len()`, `append()`,
`head()`, `tail()`

``` q
x = as_vector(1:10)

print(x[1:3])
print(append(x, [11, 12, 13], 14, "15"))
print(head(x, 10))
```

    #> NumericVector with 3 elements
    #> [1, 2, 3]
    #> Vector with 15 elements
    #> [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, "15"]
    #> NumericVector with 10 elements
    #> [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

Inspired by R, operators are vectorized element-wise. Or in numpy’s
terms, they are “broadcasted”.

``` q
[1, 2, 3] + [4, 5, 6]
```

    #> [5, 7, 9]

``` q
[1, 2, 3] * [4, 5, 6]
```

    #> [4, 10, 18]

``` q
["hello ", "good "] + ["world", "morning"]
```

    #> ["hello world", "good morning"]

Elements are recycled only if it has lenght 1 or is a scalar.

``` q
[1, 2, 3] * 2
```

    #> [2, 4, 6]

``` q
[1, 2, 3] + [4]
```

    #> [5, 6, 7]

``` q
[1, 2, 3] + [4, 5]
```

    #> ERROR: Incompatible vector lengths, left=3 and right=2

Boolean indexing works as well

``` q
s = random(10, 1, 3)
s[s > 1.5]
```

    #> [2.2093205759592394, 2.881018176090025, 2.3291201064369806, 1.8754283743739604, 1.8492749941425313, 2.373646145734219, 1.6018237211705741]

### Dictionaries

You can create a hash table structure in Q called a dictionary with a
pair of `{`, similar to Python, except that you don’t have to quote the
keys.

``` q
property = "functional"
q = {name: "Q", age: 0, property: true}
print(q)

q["age"] = q["age"] + 1

print(keys(q))
print(values(q))
```

    #> {"age": 0, "functional": true, "name": "Q"}
    #> CharacterVector with 3 elements
    #> ["age", "functional", "name"]
    #> Vector with 3 elements
    #> [1, true, "Q"]

### Control flows

Q supports `for ... in` loops and `if eles` conditions. Both the
iteration and condition need to be put in parentheses.

``` q
for (name in ["Q", "R", "Python"]) {
  if (name != "Python") {
    print(name)
  } else {
    print("I don't like Python")
  }
}
```

    #> "Q"
    #> "R"
    #> "I don't like Python"

`for in` and `if` blocks have their own scopes. Inside their inner
scope, we can (recursively) access and rebind a variable name defined in
the outer scope using assignment without the `let` keyword like so.

``` q
result = []
for (i in 1:3) {
  result = append(result, i)
}
result
```

    #> [1, 2, 3]

Or you can create vector with specified length with `vector()` and then
start filling in the elements with indexing.

``` q
result = vector(3)
for (i in 1:3) {
  result[i] = i
}
result
```

    #> [1, 2, 3]

However, if we declare a variable in the inner scope with the `let`
keyowrd, that variable will shadow the outer scope variable and get lost
when the block ends. For example, the following code will not work as
expected.

``` q
flag = true
if (flag) {
  let flag = false
}
flag
```

    #> true

### Functions

Function definition uses the R style, simply create a function object
with the `fn` keyword and then bind it to a name. Default arguments and
named arguments are supported.

``` q
add = fn(x, y = 1, z = 1) {
  x + y + z * 2
}

add(1, z = 2)
```

    #> 6

Functions in Q are first-class citizens. They can be passed around as
arguments and returned from other functions. There is a `return` keyword
but functions can also use implicit returns. Here we define a `map`
function that takes a function and a vector and applies the function to
each element of the vector.

``` q
map = fn(arr, f) {
    result = vector(len(arr))
    for (i in arr) {
        result[i] = f(arr[i])
    }
    result
}

[1, 2, 3] |> map(fn(x) x * 2)
```

    #> [2, 4, 6]

Of course the preferred the way to to double a vector is to simply use
the vectorized operator `*`.

Another example of implementing a `filter()` function that take a vector
and a predicate function and returns a vector with only the elements
that satisfy the predicate.

``` q
filter = fn(x, f) {
  result = []
  for (i in 1:len(x)) {
    if (f(x[i])) {
      result = append(result, x[i])
    }
  }
  result
}

[
  {name: "Ross",   job: "Paleontology"},
  {name: "Monoca", job: "Chef"}
] |> filter(fn(x) x["job"] == "Chef")
```

    #> [{"name": "Monoca", "job": "Chef"}]

## Next steps

- `...` for variadic arguments and spread operator

- index tests for vector and dict

- dataframe interface

- improve error message with token col and line

- more standard library functions
