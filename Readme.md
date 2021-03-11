# ZError

[中文文档](./Readme_zh.md)

## Features

- Predefine errors
- Classify errors with groups
- Make error code predefined and standard
- Code Name automatically generated by group name and error name, which can be customizable
- Hide -or not error infos from client
- Some built-in error defs
- Standard response format: {"code": "${groupName}:${errorName}", "data": ${interface}}
- list all error codes and corresponding infos(descriptions, http code, etc)


## Examples:

```go


func DoSomething()error {
    ***
    return zerror.BadRequest.New()
}

```


then you can do many fantastic things:

you can return http code predefined by error def in web request, or return corresponding grpc code in grpc request

you can log any thing 

full examples see in examples directory


## Why is zerror made


I want errors to be clearly organized and resuable and classified in groups

I want hide erors from client and return code instead, clients use code to do interactive things

I want to respond predefined code and log error in one line code

so Zerror is here



it's recommended that error will only be wrapped once, otherwise which error we should determine it is?