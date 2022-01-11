module mmo_game_zinx

go 1.17

require (
	"zinx" v0.0.0
	"mmo_game_zinx" v0.0.0
)
replace (
	"zinx"  => ../zinx
	"mmo_game_zinx" => ./
)