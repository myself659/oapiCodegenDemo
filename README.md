## install oapi-codegen

```
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
oapi-codegen -version
```

## How to generate the code

### client
```
oapi-codegen --package=main --generate types,client openapi.yaml > openapi.gen.go
```

### server


```
oapi-codegen --package=main  --generate types,server openapi.yaml > generatedserver.go
```
