all:
#	gcc -o websocket websocket.c -I../include -L../build -lwebsite -lev -lcrypto
#	6g client.go
#	6l -o client client.6
	6g conn.go
	6l -o conn conn.6
	6g test.go
	6l -o test test.6
