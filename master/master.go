package master

import (
	."../definitions"
	"fmt"
)



func MasterLoop(isMaster chan bool, masterMessage chan Message, peerChan chan PeerUpdate){
	fmt.Println("MasterLoop")
	companions := PeerUpdate{}
	messageChan := make(chan Message)
	for{
		select{
		case  <- isMaster:
			fmt.Println("I AM MASTA")
			for{
				go calculateOptimalElevator(messageChan. companions.Peers)
				select{
				case <- isMaster:
					break

				case message := <- masterMessage:
					if (message.MsgType == 1){
						messageChan <- message
					}

				case message := <- messageChan:
					fmt.Println()
				case companions = <- peerChan:
					fmt.Println(companions)

				default:
					//masterloop
					
				}
			}
		}
	}
}

func calculateOptimalElevator(messageChan chan Message, companions []string){
	for{
		select{
		case message := <-messageChan:
			fmt.Println("Calculating optimal elevator")
			leastCostID := ""
			firstZero := 0
			sender := message.SenderID
			orders  := message.Order

			//Calculate optimal elevator for external up orders
			for order := range orders.ExtUpOrders{
				if(order == 0){
					continue
				} else{
					//leastCostID, firstZero = calculateOptimalElevatorAssignment(message, companions, order)
					temp := message.Elevators[leastCostID]
					temp.Requests.ExtUpOrders[order] = 1
					if(firstZero == -1){
						continue
					} else{
						temp.Queue[firstZero] = order
						message.Elevators[leastCostID]= temp
					}
				}
			}
			//Calculate optimal elevator for external down orders
			for order := range orders.ExtDwnOrders{
				if(order == 0){
					continue
				} else{
					//leastCostID, firstZero = calculateOptimalElevatorAssignment(message, companions, order)
					temp := message.Elevators[leastCostID]
					temp.Requests.ExtDwnOrders[order] = 1
					if(firstZero == -1){
						continue
					} else{
						temp.Queue[firstZero] = order
						message.Elevators[leastCostID]= temp
					}
				}
			}
			//Give internal orders to the right elevator
			queue := message.Elevators[sender].Queue
			for order := range orders.IntOrders{
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
			messageChan <- message
		}
}

func calculateOptimalElevatorAssignment(message Message, companions []string, order int) (string, int){
	leastCostID := ""
	leastCost := 0
	cost := 0
	firstZero := 0
	lastElement := 0
	for _,companion := range companions{
		companionFloor := message.Elevators[companion].Floor
		companionQueue := message.Elevators[companion].Queue
		for i,queueElement := range companionQueue{
			if( order == queueElement){
				return companion, -1
			}
			if(queueElement == 0){
				if(firstZero != 0){
					firstZero = i
				}
				continue
			} else{
				if(i == 0){ lastElement = companionFloor}else
				{lastElement = companionQueue[i-1]}
				cost = cost + abs(queueElement - lastElement)
			}
		}
		cost = cost + abs(companionFloor-order)
		fmt.Println(companion)
		fmt.Println(order)
		fmt.Println(cost)
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

func AmIMaster(message Message, masterID string, UDPoutChan chan Message, localIP string){
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