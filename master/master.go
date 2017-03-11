package master

import (
	."../definitions"
	."fmt"
	"strconv"
)



func MasterLoop(isMaster 	chan bool, 			masterMessage 	chan Message, 
				peerChan 	chan PeerUpdate, 	UDPoutChan 		chan Message){
	
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
					change1 := false
					change2 := false
					if (messageBackup.MsgType == 1){
						messageBackup.Elevators, change1 = isTheElevatorFinished(messageBackup.Elevators, senderID)
						messageBackup.Elevators, change2 = calculateOptimalElevator(messageBackup.Elevators, senderID)
						Println(change1)
						Println(change2)
						if (change1 || change2) {
							Println("There was a change!")
							for _,slave := range slaves.Peers{
								messageBackup.RecieverID = slave
								messageBackup.MsgType = 2
								Println(messageBackup)
								UDPoutChan <- messageBackup
							}
						}
					}
				case slaves = <- peerChan:
					Println("maser peerupdate")
					Println(slaves)
				}
			}
		}
	}
}

func isTheElevatorFinished(slaves map[string]Elevator, senderIP string) (map[string]Elevator, bool){
	slavePointer := make(map[string]*Elevator)
	var slavetemp [ELEVATORS]Elevator
	change := false
	i :=0

	//making a map with pointers that is possible to change
	for key := range slaves{
		slavetemp[i] = slaves[key]
		slavePointer[key] = &(slavetemp[i])
		i++
	}
	
	sender := (slavePointer[senderIP])


	if (sender.Position != 0){
		if(sender.Position == sender.Queue[0]){
			change = true 
			for i := range sender.Queue{
				if (i == (len(sender.Queue)-1)){
					sender.Queue[i] = 0
 				} else{
 					Println(strconv.Itoa(i))
					sender.Queue[i] = sender.Queue[i+1] 
				}
			}
		}
	}


	//converting the map back to normal map without pointers
	elementMap := make(map[string]Elevator)
	for key, element := range slavePointer{
		elementMap[key] = *element
	}
	return elementMap, change
}

func calculateOptimalElevator(slaves map[string]Elevator, senderIP string) (map[string]Elevator, bool){
	Println("Calculating optimal elevator")
	leastCostID := ""
	firstZero := 0
	slavePointer := make(map[string]*Elevator)
	change := false
	var slavetemp [ELEVATORS]Elevator
	i :=0

	//making a map with pointers that is possible to change
	for key := range slaves{
		slavetemp[i] = slaves[key]
		slavePointer[key] = &(slavetemp[i])
		i++
	}
	senderElevator := (slavePointer[senderIP])
	orders := senderElevator.Order

	//Calculate optimal elevator for external up orders
	for i,order := range orders.ExtUpButtons{
		if(order == 0){
			continue
		} else{
			change = true
			leastCostID, firstZero = calculateOptimalElevatorAssignment(slavePointer, i+1)
			optimalSlave := slavePointer[leastCostID]
			(*optimalSlave).Light.ExtUpButtons[i] = 1
			if(firstZero == -1){
				continue
			} else{
				(*optimalSlave).Queue[firstZero] = i+1
			}
		}
	}
	//Calculate optimal elevator for external down orders
	for i,order := range orders.ExtDwnButtons{
		if(order == 0){
			continue
		} else{
			change = true
			leastCostID, firstZero = calculateOptimalElevatorAssignment(slavePointer, i+1)
			optimalSlave := *(slavePointer[leastCostID])
			optimalSlave.Light.ExtUpButtons[i] = 1
			if(firstZero == -1){
				continue
			} else{
				optimalSlave.Queue[firstZero] = i+1
			}
		}
	}
	//Give internal orders to the right elevator
	senderElevator = (slavePointer[senderIP])	
	for i,order := range orders.IntButtons{
		if (order == 0){
			continue
		} else {
			change = true
			senderElevator.Light.IntButtons[i] = 1
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
	return elementMap, change
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