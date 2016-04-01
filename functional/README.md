# Functional Test


Used to test the project with the frontend and client libraries.

-------

### Requeriments

Files:

* functional-config.json
* key.ssl
* cert.ssl

-------

### Running Test

Execute commands in `functional` folder.

```go
go run client.go
go run ../main.go -config ./functional-config.json -logtostderr
```

Open `http://localhost:5000/SpecRunner.html`