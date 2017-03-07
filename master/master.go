package master

import (
	."../definitions"
	"../udp/peers"
	"fmt"
)



func MasterLoop(isMaster chan bool, masterMessage chan Message, peerChan chan peers.PeerUpdate){
	fmt.Println("MasterLoop")
	companions := peers.PeerUpdate{}
	for{
		select{
		case  <- isMaster:

			select{
			case <- isMaster:
				break

			case message := <- masterMessage:
				if (message.MsgType == 1){
					message = calculateOptimalElevator(message, companions.Peers)
				}

			case companions = <- peerChan:

			default:
				//masterloop
				fmt.Println("I AM MASTA")
			}
		}
	}
}



func calculateOptimalElevator(message Message, companions []string) (Message) {
	leastCostID := ""
	leastCost := 0
	cost := 0
	sender := message.SenderID
	orders  := message.Orders

	//Calculate optimal elevator for external orders
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
							if(firstZero != 0){
								firstZero = i
							}
							continue
						} else{
							cost = cost + abs(queueElement - companionQueue[i-1])
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
				message.Elevators[leastCostID].Requests.ExtUpOrders[order] = 1
			}
		}
	}
	//Calculate optimal elevator for external orders
	for _,order := range orders.ExtDwnOrders{
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
							if(firstZero != 0){
								firstZero = i
							}
							continue
						} else{
							cost = cost + abs(queueElement - companionQueue[i-1])
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
				message.Elevators[leastCostID].Requests.ExtDwnOrders[order] = 1
			}
		}
	}
	//Give internal orders to the right elevator
	queue = message.Elevators[sender].Queue
	for _,order := range orders.IntOrders{
		if (order == 0){
			continue
		} else {
			for i,queueElement := range queue{
				if( order == queueElement){
					break
				} 
				if (queueElement == 0){
					message.Elevators[sender].Queue[i] = order
					message.Elevators[sender].Requests.IntOrders[order] = 1
				}
			}
		}
	}
	return message
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