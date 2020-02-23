export GOOS=js
export GOARCH=wasm
go build -o main.wasm main.go
cp "c:/Program Files/Go/misc/wasm/wasm_exec.js" "./wasm_exec.js"
echo "Done"
sleep 9999
