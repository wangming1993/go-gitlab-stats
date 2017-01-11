run:
	go run main.go

clean:
	rm -rf htmls/*

init:
	[ -e htmls ] || mkdir htmls

fmt:
	gofmt -e -w lib/*.go
