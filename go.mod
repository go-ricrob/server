module github.com/go-ricrob/server

go 1.20

require (
	github.com/go-ricrob/exec v0.0.8
	github.com/go-ricrob/game v0.0.6
)

require golang.org/x/exp v0.0.0-20230725093048-515e97ebf090 // indirect

replace (
	github.com/go-ricrob/exec => ../exec
	github.com/go-ricrob/exec/task => ../exec/task
)
