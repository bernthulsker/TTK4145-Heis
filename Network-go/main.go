package main

import (
	"./udp"
	"./master"
	. "./definitions"
	"./udp/peers"
	//"time"
	"fmt"
)



func main(){
	//Make channels
	outgoingMsg := make (chan Message)
	incomingMsg := make (chan Message)
	peerChan := make(chan peers.PeerUpdate)
	amIMaster := make(chan bool)
	masterIDChan := make(chan string)
	masterID := ""

	go master.MasterLoop(amIMaster)

	localIP := udp.UDPInit(outgoingMsg, incomingMsg, amIMaster, masterIDChan, peerChan)

	go treatMessages(incomingMsg, outgoingMsg, masterIDChan, masterID, localIP)

	masterID = udp.MasterInit(peerChan, amIMaster, localIP, masterIDChan, outgoingMsg)

	go udp.UDPUpkeep(peerChan, amIMaster, localIP, masterIDChan, masterID)

	for{
		select{
		case masterID = <- masterIDChan:

		}
	}
}


func treatMessages(incomingMsg chan Message, outgoingMsg chan Message, masterIDChan chan string, masterID string, localIP string){
	fmt.Println("Treat Messages")
	for{
		select{
		case message := <- incomingMsg:
			//TREAT THE MESSAGE
			if (message.MsgType == 3){
				fmt.Println("Someone asked if I were master" + localIP)
				fmt.Println(masterID)
				if(masterID == localIP){
					fmt.Println("I am a master and my ID is" + localIP)
					message.MsgType = 4
					message.SenderID = localIP
					outgoingMsg <- message
				}
			}
			if(message.MsgType == 4){
				fmt.Println("I were told that this guy is master : " + message.SenderID)
				masterIDChan <- message.SenderID
			}	
		}
	}
}


































//Clutter

	/*go func (){
		msg := Message{}
		msg.MsgType = 0
		for{
			outgoingMsg <- msg
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
		amIMaster <- true
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
	}*/