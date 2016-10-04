# basic build sript

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd $DIR/frontend
go get  
gopherjs build ex.go
 
cp ex.js $DIR/backend
cp ex.js.map $DIR/backend

cd $DIR/backend
go run ex.go
