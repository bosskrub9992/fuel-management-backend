#### generate mock

1. install mockery
```sh
brew install mockery
```

2. run at root
```sh
mockery --dir=internal/services --all --outpkg=mocks --output=internal/services/mocks --case=snake --with-expecter
```