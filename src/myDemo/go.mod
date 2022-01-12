module myDemo

go 1.17

require (
	github.com/golang/protobuf v1.5.2
	zinx v0.0.0
)

require google.golang.org/protobuf v1.26.0 // indirect

replace zinx => ../zinx
