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
	peerMasterChan := make(chan peer.PeerUpdate)
	isMaster := make(chan bool)
	masterIDChan := make(chan string)
	masterMessage := make(chan Message)
	masterID := ""

	go master.MasterLoop(isMaster, masterMessage, peerMasterChan)

	localIP := udp.UDPInit(UDPoutChan, UDPinChan, isMaster, masterIDChan, peerChan)
	masterID = udp.MasterInit(peerChan, isMaster, localIP, UDPoutChan)

	go treatMessages(UDPinChan, UDPoutChan, masterIDChan, masterID, localIP)

	fmt.Println("Waiting for ID...")
	select{
	case masterID = <- masterIDChan:
		fmt.Println("masteridchan" + masterID)
	}

	go udp.UDPUpkeep(peerChan, peerMasterChan, isMaster, localIP, masterIDChan, masterID, UDPoutChan)

	for{
		fmt.Println("I am in the ending loop")
		time.Sleep(time.Second*5)
	}


func treatMessages(UDPinChan chan Message, UDPoutChan chan Message, masterMessage chan Message, masterIDChan chan string, masterID string, localIP string){
	fmt.Println("Treat Messages")
	for{
		select{
		case message := <- UDPinChan:
			if (message.MsgType == 1){
				masterMessage <- message
			}

			if (message.MsgType == 3){
				fmt.Println("Someone asked if " + localIP + " is master")
				Master.amIMaster(message, masterID, UDPoutChan, localIP)
			}
			if(message.MsgType == 4){
				fmt.Println("I was told that " + message.SenderID + " is the master")
				masterIDChan <- message.SenderID
				masterID = message.SenderID
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
	}*/