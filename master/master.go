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
	firstZero := 0
	sender := message.SenderID
	orders  := message.Order

	//Calculate optimal elevator for external orders
	for _,order := range orders.ExtUpOrders{
		if(order == 0){
			continue
		} else{
			leastCostID, firstZero = calculateOptimalElevatorAssignment(message, companions, order)
			if(leastCostID == ""){
				continue
			} else{
				temp := message.Elevators[leastCostID]
				temp.Queue[firstZero] = order
				temp.Requests.ExtUpOrders[order] = 1
				message.Elevators[leastCostID]= temp
			}
		}
	}
	//Calculate optimal elevator for external orders
	for _,order := range orders.ExtDwnOrders{
		if(order == 0){
			continue
		} else{
			leastCostID, firstZero = calculateOptimalElevatorAssignment(message, companions, order)
			if(leastCostID == ""){
				continue
			}  else{
				temp := message.Elevators[leastCostID]
				temp.Queue[firstZero] = order
				temp.Requests.ExtDwnOrders[order] = 1
				message.Elevators[leastCostID]= temp
			}
		}
	}
	//Give internal orders to the right elevator
	queue := message.Elevators[sender].Queue
	for _,order := range orders.IntOrders{
		if (order == 0){
			continue
		} else {
			for i,queueElement := range queue{
				if( order == queueElement){
					break
				} 
				if (queueElement == 0){
					temp := message.Elevators[sender]
					temp.Queue[i] = order
					temp.Requests.IntOrders[order] = 1
					message.Elevators[sender] = temp
				}
			}
		}
	}
	return message
}

func calculateOptimalElevatorAssignment(message Message, companions []string, order int) (string, int){
	leastCostID := ""
	leastCost := 0
	cost := 0
	firstZero := 0
	for _,companion := range companions{
		companionFloor := message.Elevators[companion].Floor
		companionQueue := message.Elevators[companion].Queue
		for i,queueElement := range companionQueue{
			if( order == queueElement){
				leastCostID = ""
				break
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
	return leastCostID, firstZero
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

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}