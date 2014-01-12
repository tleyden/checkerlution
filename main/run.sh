
rm -f out.txt
GOMAXPROCS=8 go run main.go -feed longpoll -syncGatewayUrl http://localhost:4984/checkers -team RED -randomDelayBeforeMove 0 2>&1 | tee out.txt
