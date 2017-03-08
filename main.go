package main

import (
	"./udp"
	"./master"
	. "./definitions"
	"./udp/peers"
	"time"
	"fmt"
)




func main(){
	//Make channels
	UDPoutChan := make (chan Message)
	UDPinChan := make (chan Message)
	peerChan := make(chan peers.PeerUpdate)
	peerMasterChan := make(chan peers.PeerUpdate)
	isMaster := make(chan bool)
	masterIDChan := make(chan string)
	masterMessage := make(chan Message)
	masterID := ""

	go master.MasterLoop(isMaster, masterMessage, peerMasterChan)

	localIP := udp.UDPInit(UDPoutChan, UDPinChan, isMaster, masterIDChan, peerChan)
	masterID = udp.MasterInit(peerChan, isMaster, localIP, UDPoutChan)

	go treatMessages(UDPinChan, UDPoutChan, masterMessage, masterIDChan, masterID, localIP)

	if(masterID == ""){
		fmt.Println("Waiting for ID...")
		select{
		case masterID = <- masterIDChan:
			fmt.Println("masteridchan" + masterID)
		}
	}
	go udp.UDPUpkeep(peerChan, peerMasterChan, isMaster, localIP, masterIDChan, masterID, UDPoutChan)
	if( localIP == "Bob"){
		requests := Orders{};
		array := [4]int{1, 1, 1, 1}
		queue := [4]int{3, 2, 0, 0}
		orders := Orders{array, array, array};
		elevators := make(map[string]Elevator)
		elevators["Alice"] = Elevator{true,1,1,requests, queue}
		queue = [4]int{1, 0, 0, 0}
		elevators["Bob"] = Elevator{true,3,1,requests, queue}
		message := Message{elevators, orders, "Bob", "Alice", 1}

		for{
			UDPoutChan <- message
			fmt.Println("I am in the ending loop")
			time.Sleep(time.Millisecond * 100)
		}
	}
	for{
		fmt.Println("I am in the ending loop")
		time.Sleep(time.Second*5)
	}
}


func treatMessages(UDPinChan chan Message, UDPoutChan chan Message, masterMessage chan Message, masterIDChan chan string, masterID string, localIP string){
	fmt.Println("Treat Messages")
	for{
		select{
		case message := <- UDPinChan:
			if (message.MsgType == 1 && localIP == masterID){
				fmt.Println("I got an order and my ID is " + localIP)
				masterMessage <- message
				break
			}

			if (message.MsgType == 3){
				fmt.Println("Someone asked if " + localIP + " is master")
				master.AmIMaster(message, masterID, UDPoutChan, localIP)
				break
			}
			if(message.MsgType == 4){
				fmt.Println("I was told that " + message.SenderID + " is the master")
				masterIDChan <- message.SenderID
				masterID = message.SenderID
				break
			}
		case masterID = <- masterIDChan:	
		}
	}
}


































//Clutter

	/*go func (){
		msg := Message{}
		msg.MsgType = 0
		for{
			UDPoutChan <- msg
			msg.MsgType++
			time.Sleep(time.Second)
		}
	}()*/


	/*fmt.Println("masterInit")
	companions := peerInfo.Peers
	companions := make([]string, len(peerInfo.Peers), (cap(peerInfo.Peers)+1))
	copy(companions, )
	for _, companion := range peerInfo.Peers{
		fmt.Println(companion)
	}
	if (len(companions) == 1 ){
		fmt.Println("heyo")
		masterID = localIP
		isMaster <- true
		return
	}
	for _, companion := range companions{
		msg := Message{}
		msg.MsgType = 3
		msg.SenderID = localIP
		msg.RecieverID = companion
		outgoing <- msg
	}
	select{
	case masterID = <-masterIDChan:
			return
	}



	timeChan := time.NewTimer(time.Second*5).C
	select{
	case masterID = <-masterIDChan:
		fmt.Println("masterInitFinished")
	case <- timeChan:
		for _,companion := range companions{
			if (companion < masterID){
				masterID = companion
			}
		}
	}

	*/