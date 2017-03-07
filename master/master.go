package master

import (
	."../definitions"
	"fmt"
)



func MasterLoop(isMaster chan bool, masterMessage chan Message, peerChan chan peers.PeerUpdate){
	fmt.Println("MasterLoop")
	var companions
	for{
		select{
		case  <- isMaster:

			select{
			case <- isMaster:
				break

			case message := <- masterMessage:
				if (message.MsgType == 1){
					calculateOptimalElevator(message, companions.Peers)
				}

			case companions = <- peerChan:

			default:
				//masterloop
				fmt.Println("I AM MASTA")
			}
		}
	}
}



func calculateOptimalElevator(message Message, companions []string){
	var leastCostID
	var cost
	for companion := range companions{
		companionFloor = message.Elevators[companion].floor
		if(companionFloor == 0){										//if the element doesn't exist in the map
			continue
		}
		companionQueue = message.Elevators[companion].Queue
		orderedFloor = message.order
		for queueElement :=range companionQueue{

		}
	}
}

func amIMaster(message Message, masterID string, UDPoutChan chan Message, localIP string){
	if(masterID == localIP){
		fmt.Println("I am a master and my ID is " + localIP)
		message.MsgType = 4
		message.RecieverID = message.SenderID
		message.SenderID = localIP
		UDPoutChan <- message
	}
}