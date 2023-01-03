go build -o bin/ccproxy-server cmd/server/main.go
scp -P 46517 bin/ccproxy-server root@proxy01.ha-du32h098e72.medicalstus.ir:"/tmp/ccproxy-server"