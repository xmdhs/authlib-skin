SET CGO_ENABLED=1
go build -trimpath -ldflags "-w -s" -tags="redis,sqlite"
