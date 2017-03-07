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
	leastCostID := ""
	leastCost := 0
	cost := 0
	orders  = message.Orders
	for _,order := range orders.ExtUpOrders{
		if(order == 0){
			continue
		} else{
			CostLoop:
				for _,companion := range companions{
					companionFloor = message.Elevators[companion].floor
					companionQueue = message.Elevators[companion].Queue
					for i,queueElement := range companionQueue{
						if( order == queueElement){
							break CostLoop
							cost = -1
						}
						if(queueElement == 0){
							firstZero = i
							continue
						} else{
							cost = cost + abs(queueElement - order)
						}
					}
					cost = cost + abs(companionFloor-order)
					if(leastCost == 0 || leastCost > cost){
						leastCost = cost
						leastCostID = companion
					}else{
						firstZero = 0
					}
					cost = 0
				}
			if(cost == -1){
				continue
			} else{
				message.Elevators[leastCostID].Queue[firstZero] = order
			}
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