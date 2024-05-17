# Errs
Is a library enhancing the error handling by providing set of functionalities to track, format and enrich errors with extra data in a standarized way.

## Motivation
We have our beautiful application where we used idiomatic error handling. Three necessary
lines after almost each function:

```
if err != nil {
    return err
}
```
We run it and...
```
invalid character 'a' looking for beginning of value
```
Where? How? What? Nothing in the error indicates where this error comes from. An issue with
go's errors is that they don't track any data except the message. If we have a function used
in many places that returns a specific error. It is obvious what function returned the error,
but it is hard to figure out the function that called this function.

Of course there is a way to track the error with proper discipline. We could always pass the error to client code more or less like this:
```
func g() error {
    err := f()
    if err != nil {
        return fmt.Errorf("the error %w was returned in function g")
    }
    return nil
}
```
The we could use errors unwrapping. Tracking the path the error went from the place it's been created to the place it's supported now is easier to track. 
However, the irregular and unstructurized way of passing the error may not be very elegant or unpractical. Again, it is not something that can't be handled with proper discipline.

The **errs** library, helps in returning structurized errors and tracking it's origin by
wrapping the errors and gathering data on each stage. Then the error may be formatted by  predefined or custome formatters.


## Usage

Let's take the sample code again:
```
func g() error {
    err := f()
    if err != nil {
        return fmt.Errorf("the error %w was returned in function g")
    }
    return nil
}
```
Using the package we could write it this way:
```
func g() error {
    err := f()
    if err != nil {
        return errs.Wrap(err)
    }
    return nil
}
```
It requires a minimal effort but now we can track the information of the file and line
when the error was wrapped. Using the defualt formatter when we print the error the standard
way:
```
fmt.Printf("%v\n", err)
```
The result will be similar to that:
```
mymodule/myfile.go:3: original error returned from f()
```

if we want to use g() in another function:
```
func h() error {
    err := g()
    if err != nil {
        return errs.Wrap(err)
    }
    return nil
}
```
The printer error will look similar to this:
```
mymodule/myotherfile.go:15:
    mymodule/myfile.go:3: original error returned from f()
```


### Adding messages and arguments
It is possible to add additional message while wrapping the error or put an argument.
Let's change functions g and h to this:
```
func g(i int) error {
    err := f()
    if err != nil {
        return errs.Wrap(err).Arg("argument", i)
    }
    return nil
}

func h(float s) error {
    err := g(5)
    if err != nil {
        return errs.Wrap(err).Msg("something went wrong").Arg("argument", s)
    }
    return nil
}
```
The error will print message like:
```
mymodule/myotherfile.go:15: something went wrong: params: argument=4.5
    mymodule/myfile.go:3: original error returned from f(): params: argument: 5
```

### Custom formatters
The error predefined error formattings don't need to suit everybody's needs. It is possible
to define own formatter and provide it as the default one.
It takes implementing the interface:
```
type Formatter interface {
	Format(errData []ErrorData) string
}
```
The errData is the list a data struct providing informations about the error. The list contains error's chain from top to bottom. 

Then the formatter may be set using:
```
errs.SetFormatter(NewMyFormatter())
```

