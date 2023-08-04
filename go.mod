module github.com/go-ricrob/server

go 1.20

require (
	github.com/go-ricrob/exec v0.0.12
	github.com/go-ricrob/game v0.0.7
)

require golang.org/x/exp v0.0.0-20230801115018-d63ba01acd4b // indirect

//replace (
//	github.com/go-ricrob/exec => ../exec
//	github.com/go-ricrob/exec/task => ../exec/task
//
//	github.com/go-ricrob/game => ../game
//	github.com/go-ricrob/game/board => ../game/board
//)
