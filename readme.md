# Email Go Checker 

> Simple e-mail validation utility

* Validates e-mail syntax
* Validates e-mail domain through MX Lookup
* Tries to valide e-mail "existence" using SMTP protocol

## Usage

### CLI
```
# go run
go run . -email="email-that-needs-validation@doamin.com" -email-from-test="some-valid-email@email.com"

# compiled bin
./email-go-checker -email="email-that-needs-validation@doamin.com" -email-from-test="some-valid-email@email.com"
```

### API Server
```
# go run
go run . -serve

# compiled bin
./email-go-checker -serve
```
