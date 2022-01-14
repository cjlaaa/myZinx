module mmo_game_zinx

go 1.17

require (
	github.com/golang/protobuf v1.5.0
	zinx v0.0.0
)

require google.golang.org/protobuf v1.27.1 // indirect

replace (
	mmo_game_zinx => ./
	zinx => ../zinx
)
