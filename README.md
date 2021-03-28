# P2P Matrix Go

This application is the backend for the [P2P Matrix project](https://github.com/marchellodev/p2p_matrix). it is a
testing environment for P2P algorithms that simulates the Internet connections (speed limitations and latency depending
on the geographical location)

## Running

`main.go` contains an example network - where every node knows everyone and has every file on the network. While not
efficient, let alone practical, it is a perfect testing network for this project.

Before running it, make sure that you have [Golang](https://golang.org/dl) installed.

Then, execute:
```shell
go run main.go
```
<br>

To create a shared library for P2P Matrix, run either of those depending on your system:
```shell
go build -buildmode c-shared -o lib.so            # Linux
go build -buildmode c-shared -o lib.dll           # Windows
go build -buildmode c-shared -o lib.dylib         # MacOS
```
<br>

## License

MIT
