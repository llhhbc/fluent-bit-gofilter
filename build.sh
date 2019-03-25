
##go build -o call.la -buildmode=c-archive call.go

go build -o cgolib.so -buildmode=c-shared cgolib.go
