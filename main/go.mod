module main

go 1.20

// директивой replace указываем положение корня
// модуля server относительно main/go.mod
replace server => ../server

require server v0.0.0-00010101000000-000000000000

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
)
