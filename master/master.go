package master

import (
	."../definitions"
	."fmt"
	"strconv"
)



func MasterLoop(isMaster chan bool, masterMessage chan Message, peerChan chan PeerUpdate){
	Println("MasterLoop")
	slaves := PeerUpdate{}
	messageBackup := Message{}
	for{
		select{
		case  <- isMaster:
			Println("I AM MASTA")
			Master:
			for{
				select{
				case <- isMaster:
					Println("Break master")
					break Master

				case messageBackup = <- masterMessage:
					Println("master recieved message")
					senderID := messageBackup.SenderID
					if (messageBackup.MsgType == 1){
						messageBackup.Elevators = calculateOptimalElevator(messageBackup.Elevators, senderID)
					}
					Println("This is the calculated queue")
					Println(messageBackup)
				case slaves = <- peerChan:
					Println("maser peerupdate")
					Println(slaves)
					
				}
			}
		}
	}
}


func calculateOptimalElevator(slaves map[string]Elevator, sender string) (map[string]Elevator){
	Println("Calculating optimal elevator")
	leastCostID := ""
	firstZero := 0
	slavePointer := make(map[string]*Elevator)
	var slavetemp [ELEVATORS]Elevator
	i :=0

	//making a map with pointers that is possible to change
	for key := range slaves{
		slavetemp[i] = slaves[key]
		slavePointer[key] = &(slavetemp[i])
		i++
	}
	senderElevator := (slavePointer[sender])
	orders := senderElevator.Order

	//Calculate optimal elevator for external up orders
	for i,order := range orders.ExtUpOrders{
		if(order == 0){
			continue
		} else{
			leastCostID, firstZero = calculateOptimalElevatorAssignment(slavePointer, i+1)
			Println( "IP: " + leastCostID + " Order: " + strconv.Itoa(i+1) + " FirstZero: " + strconv.Itoa(firstZero))
			optimalSlave := slavePointer[leastCostID]
			(*optimalSlave).Light.ExtUpOrders[i] = 1
			if(firstZero == -1){
				continue
			} else{
				(*optimalSlave).Queue[firstZero] = i+1
			}
			Println(optimalSlave)
			Println(slavePointer[leastCostID])
		}
	}
	//Calculate optimal elevator for external down orders
	for i,order := range orders.ExtDwnOrders{
		if(order == 0){
			continue
		} else{
			leastCostID, firstZero = calculateOptimalElevatorAssignment(slavePointer, i)
			optimalSlave := *(slavePointer[leastCostID])
			optimalSlave.Light.ExtUpOrders[i] = 1
			if(firstZero == -1){
				continue
			} else{
				Println("Placing order in queue")
				optimalSlave.Queue[firstZero] = i+1
			}
		}
	}
	//Give internal orders to the right elevator
	senderElevator = (slavePointer[sender])	
	for i,order := range orders.IntOrders{
		if (order == 0){
			continue
		} else {
			senderElevator.Light.IntOrders[i] = 1
			for j,queueElement := range senderElevator.Queue{
				if( i+1 == queueElement){
					break
				} 
				if (queueElement == 0){
					senderElevator.Queue[j] = i+1
					break
				}
			}
		}
	}

	elementMap := make(map[string]Elevator)
	for key, element := range slavePointer{
		elementMap[key] = *element
	}
	return elementMap
}

func calculateOptimalElevatorAssignment(slaves map[string]*Elevator, order int) (string, int){
	leastCostID := ""
	leastCost := -1
	cost := 0
	firstZero := -2
	leastCostFirstZero := 0
	lastElement := 0
	for ip,slave := range slaves{
		for i,queueElement := range slave.Queue{
			if( order == queueElement){
				return ip, -1
			}
			if(queueElement == 0){
				if(firstZero == -2){
					firstZero = i
				}
				continue
			} else{
				if(i == 0){ lastElement = slave.Floor}else
				{lastElement = slave.Queue[i-1]}
				cost = cost + abs(queueElement - lastElement)
			}
		}
		cost = cost + abs(slave.Floor-order)
		if(leastCost == -1 || leastCost > cost){
			leastCost = cost
			leastCostID = ip
			leastCostFirstZero = firstZero
		}
		firstZero = -2
		cost = 0
	}
	return leastCostID, leastCostFirstZero
}

func AmIMaster(message Message, masterID string, UDPoutChan chan Message, localIP string){
	if(masterID == localIP){
		Println("I am a master and my ID is " + localIP)
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