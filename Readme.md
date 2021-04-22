# ZError

## Notice

v1 is deprecated, please use v2

import path : `github.com/EchoUtopia/zerror/v2`

example import path: `github.com/EchoUtopia/zerror/examples/v2`

[中文文档](./Readme_zh.md)

## Terminology
- Error Definition(zerror.Def): error definition to generate or wrap error

## Features

- Predefine errors
- Make error code predefined and standard
- Specified error definition Status Code for Http, Grpc and so on
- Classify errors with groups
- Some built-in error definitions, like `zerror.NotFound`, `zerror.Forbidden` and so on
- List all error codes and corresponding infos(descriptions, status code, etc)
- Extension for every error definition to do many cool things you want.


## Examples:

```go


func DoSomething()error {
    ***
    // return zerror.Internal.Wrap(errors.New(`***`))
    return zerror.BadRequest.New()
}

```


## full examples 

see in examples directory


## Why is zerror made


I want errors to be clearly organized and reusable and classified in groups

I want hide errors from client and return code instead

I want to respond predefined code and log error in one line code

I want to respond with different Status Code for Grpc/Http for error generated/wrapped by  zerror

so Zerror is here
