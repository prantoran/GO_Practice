Creating, building and running Go application:

goto desired workspace

in the terminal,

~/$ export GOPATH=`pwd`
//pwd = current directory
//`pwd`=find current directory and use the string for assignment

~/$ echo $GOPATH
//print GOPATH, preceding $ indicates environment variable
 
~/$ mkdir -p src/github.com/prantoran/hello && cd $_
// -p flag creates missing directories in the path
// $_ returns the last created/executed path string

//create go files

~/$ go run main.go
//outputs in terminal

~/$ go install ./... && $GOPATH/bin/hello
//create binary of the go files recursively, starting from the main.go
// && execute the binary created in $GOPATH/bin folder


