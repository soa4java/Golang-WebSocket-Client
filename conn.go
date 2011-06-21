package main

import (
    "websocket"
    "fmt"
    "time"
    "flag"
    "runtime"
    "os"
    "container/ring"
)

type node struct{
    ip string
    port int
    ws *websocket.Conn
    msg chan []byte
    stopped bool
}

var num_ip, connPerIP, basePort, packetRate int
func init(){
    flag.IntVar(&num_ip, "n", 4, "number of IP address to use")
    flag.IntVar(&connPerIP, "c", 20000, "number of connections per IP address")
    flag.IntVar(&basePort, "b", 20000, "Base port, default 20000")
    flag.IntVar(&packetRate, "r", 500, "Sending packet rate, default 500/s")
}

func main() {
    host := ""
    msg_count := 0
    runtime.GOMAXPROCS(2)
    flag.Parse()

    fmt.Println("Starting...")
    fmt.Println(time.LocalTime())

    var err os.Error
    msgClick := make(chan int, 1000)
    wsList := ring.New(num_ip*connPerIP)
    connBuf := make(chan byte, 200)
    r := wsList
    for ipc:=2; ipc<2+num_ip; ipc++{
	host = fmt.Sprintf("ws://30.0.1.1:8080/")
	fmt.Println(host)
	for localPort:=basePort; localPort<basePort+connPerIP; localPort++{
	    n := new(node)
	    n.ip = fmt.Sprintf("30.0.1.%d",ipc)
	    n.port = localPort
	    n.msg = make(chan []byte, 2)
	    connBuf <- 1
	    go func(){
		n.ws , err = websocket.DialBind(host, "", "http://127.0.0.2/", n.ip, n.port);
		if err != nil {
		    fmt.Println("DialBind: " + err.String())
		    n.stopped = true
		    return
		}
		<-connBuf
		n.stopped = false

		// receive message
		var msg = make([]byte, 512);
		for {
		    if _, err := n.ws.Read(msg); err != nil {
			panic("Read: " + err.String())
		    }else{
			msgClick <- 1
		    }
		}
	    }()
	    r.Value = n
	    r = r.Next()
	}
    }
    go func(){
	for{
	    <-msgClick
	    msg_count++
	    if msg_count % 10000 == 0{
		fmt.Printf("%d message received, %v\n",msg_count, time.LocalTime())
	    }
	}
    }()
    fmt.Printf("all connections are established!\n")
    fmt.Println(time.LocalTime())
    time.Sleep(1000*1000*1000*2)
    fmt.Printf("now send some traffic!\n")
    go func(){
	fmt.Printf("Ring buffer length: %d\n", wsList.Len())
	for{
	    wsList.Do(sendMsg)
	}
    }()
    for{
	var a int
	fmt.Scanf("%d", &a)
	if a > 0 {
	    packetRate = a
	    fmt.Printf("Set packet rate to %d/s\n", packetRate)
	}
    }
}
func sendMsg(r interface{}){
    n := r.(*node)
    if n.stopped {
	fmt.Printf("stopped\n")
	return
    }
    var msg = make([]byte, 512);
    for i:=0; i<512; i++ {
	msg[i] = 'A' + uint8(i%52)
    }
    if nil == msg || nil == n.ws {
	fmt.Printf("msg: %p\n", msg)
	fmt.Printf("ws: %p\n", n.ws)
	panic("got you\n")
    }
    if _, err := n.ws.Write(msg); err != nil {
	panic("Write: " + err.String())
    }
    <-time.After((1000*1000*1000) / int64(packetRate))
}
