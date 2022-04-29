Get-ChildItem .\api\tasks\rpc -Recurse -Filter *.proto | ForEach-Object {
    $path = $_ | Resolve-Path -Relative
    Write-Output "protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $path"
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $path
}
