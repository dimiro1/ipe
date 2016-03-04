# Testes Funcionais


Utilizado para testar o projeto com as bibliotecas frontend e clientes.

-------

### Requerimentos

Files:

* config.json
* ssl.key
* ssk.crt

-------

### Running Test

Executar os comandos na pasta `functional`.

```go
go run client.go
go run ../main.go -config ./config.json -logtostderr
```

Open `http://localhost:5000/SpecRunner.html`