# Email Go Checker 

> Simple e-mail validation utility

* Validades e-mail syntax
* Validades e-dmail domain with MX Lookup
* Tries to valide e-mail "existence" using SMTP protocol

## Usage

### CLI
```
# go run
gom run . -email="email-that-needs-validation@doamin.com" -email-from-test="some-valid-email@email.com"

# compiled cli
email-go-checker -email="email-that-needs-validation@doamin.com" -email-from-test="some-valid-email@email.com"
```

### API Server
```
# go run
go rin . -serve

# compiled cli
email-go-checker -serve
```