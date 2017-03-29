package main

import (
	. "./definitions"
	"./localLift"
	"./udp"
	"./master"
	."./watchDog"
	"fmt"
)




func main(){

	WatchDog()

	go stateMachine()
	go Heartbeat()

	StartBackup()
	fmt.Println("This is the end")
	

	select{}
}



func stateMachine(){
	//Make channels
	UDPoutChan 		:= make(chan Message)
	UDPinChan 		:= make(chan Message)
	masterMessage 	:= make(chan Message)
	peerChan 		:= make(chan PeerUpdate)
	peerMasterChan 	:= make(chan PeerUpdate)
	elevOut 		:= make(chan Elevator)
	elevIn 			:= make(chan Elevator)
	currentElevState:= make(chan Elevator)
	internetConnect := make(chan bool)
	isMaster 		:= make(chan bool)
	masterIDChan 	:= make(chan string)
	stateChan 		:= make(chan string)
	masterID 		:= ""
	localIP			:= ""
	state 			:= "Initialize elev"
	currentState 	:= Elevator{}
	initialized 	:= false

	for{
		StateMachine:
		switch state {	

		case "Initialize elev":

			fmt.Println(state)
			go udp.CheckInternetConnection(internetConnect)
			go localLift.Elev_driver(elevIn, elevOut)

			state = "Initialize"

		case "Initialize":

			fmt.Println(state)
			
			localIP = udp.UDPInit(UDPoutChan, UDPinChan, peerChan)
			if( localIP == ""){ 
				state = "No internet"
				break 
			}
			if (!initialized){
				go master.MasterLoop(isMaster, masterMessage, peerMasterChan, UDPoutChan)
				go treatMessages(UDPinChan, UDPoutChan, masterMessage, masterIDChan, elevIn, elevOut, currentElevState, stateChan, localIP)
				masterID = udp.MasterInit(peerChan, isMaster, peerMasterChan, localIP, UDPoutChan, masterIDChan)
				
				
				go udp.UDPUpkeep(peerChan, peerMasterChan, isMaster, masterIDChan, UDPoutChan, masterID, localIP)

				initialized = true
			}

			state = "Normal operation"
			stateChan <- state

		case "Normal operation":

			fmt.Println(state)
			currentStateCopy := currentState 							//Making a copy to avoid channel passing map pointers problems
			elevOut <- currentStateCopy

			for{
				select{
				case internet := <- internetConnect:
					if(!internet){
						state = "No internet"
						stateChan <- state
						break StateMachine
					}
				case currentState = <- currentElevState:
				}
			}

		case "No internet":

			fmt.Println(state)
			internetConnection 	:= make(chan bool)
			currentStateChan 	:= make(chan Elevator)
			currentState.Order 	 = currentState.Light 

			go localLift.LocalMode(internetConnection, currentStateChan, elevIn, elevOut, currentState)
			for{
				select{
				case internet := <- internetConnect:
					if(internet){
						fmt.Println("There is internet!")
						state = "Initialize"
						fmt.Println(state)
						fmt.Println("Do youever feel like a plastic bag drifting thought the wind wanting to start again? Do you ever feel just so paper thin")
						internetConnection <- true
						fmt.Println("FLY YOU FOOLS")
						select{
						case currentState = <- currentStateChan:
							fmt.Println(currentState)
							currentState.Order = currentState.Light 
							break StateMachine
						}
					}
				}		
			}
		}
	}
}


func treatMessages(	UDPinChan 			chan Message, 	UDPoutChan 		chan Message, 
					masterMessage 		chan Message, 	masterIDChan 	chan string, 
					elevIn 				chan Elevator, 	elevOut 		chan Elevator, 
					currentElevState 	chan Elevator,	stateChan 		chan string,
					localIP 		string){

	fmt.Println("Treat Messages")
	Elevators 				:= make(map[string]Elevator)
	messageBackup 			:= Message{Elevators, "", "", 0}
	masterID 				:= ""
	state 					:= ""
	messageBackup.Elevators[localIP] = Elevator{}
	for{
		if state == "No internet"{
			select{
			case state = <- stateChan:
			}
		}
		select{
		case tempMessage := <- UDPinChan:
			if(checkMessageValidity(tempMessage.Elevators[tempMessage.SenderID])){
				messageBackup = tempMessage
				if (messageBackup.MsgType == 1 && localIP == masterID){
					fmt.Println("I got an order and my ID is " + localIP)
					messageCopy := messageBackup    								//Making a copy to avoid channel passing map pointers problems
					masterMessage <- messageCopy
				} else if (messageBackup.MsgType == 2){
					elevIn <- messageBackup.Elevators[localIP]
					fmt.Println(localIP + " got a queue")
					currentElevState <- messageBackup.Elevators[localIP]
				} else if (messageBackup.MsgType == 3){
					fmt.Println("Someone asked if " + localIP + " is master")
					master.AmIMaster(messageBackup, masterID, UDPoutChan, localIP)
				} else if (messageBackup.MsgType == 4){
					fmt.Println("I was told that " + messageBackup.SenderID + " is the master")
					masterIDChan <- messageBackup.SenderID
					masterID = messageBackup.SenderID
				}
			}
		case masterID = <- masterIDChan:
			fmt.Println("I got a masterID" + masterID)

		case elev_status := <- elevOut:
			messageBackup.Elevators[localIP] = elev_status
			messageBackup.MsgType = 1
			messageBackup.RecieverID = masterID
			messageCopy := messageBackup    									//Making a copy to avoid channel passing map pointers problems
			UDPoutChan <- messageCopy

		case state = <- stateChan:
		}
	}
}


func checkMessageValidity(message Elevator) bool{

	if(message.Floor < 1) {return false} 
	if(message.Floor > 4) {return false}
	for _, element := range message.Order.IntButtons{
		if element > 1 || element < 0 {
			return false
		}
	}
	for _, element := range message.Order.ExtUpButtons{
		if element > 1 || element < 0 {
			return false
		}
	}
	for _, element := range message.Order.ExtDwnButtons{
		if element > 1 || element < 0 {
			return false
		}
	}
	return true
}




































//Clutter


	/*

	func sendElevator(elevOut chan Elevator){
	ones := 	[4]int{1, 1, 1, 1}
	order := 	Buttons{ones, ones, ones}
	light := 	Buttons{}
	queue := 	[4]int{3,2,0,0}
	elevator :=	Elevator{1,0,1,light,order,queue}
	for{
		elevOut <- elevator                               	
		fmt.Println("I am in the ending loop")
		time.Sleep(time.Millisecond*10)
	}
}
	go func (){
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

	if( localIP == "Bob"){
		lights := Orders{};
		array := [4]int{1, 1, 1, 1}
		queue := [4]int{0, 0, 0, 0}
		orders := Orders{array, array, array};
		elevators := make(map[string]Elevator)
		elevators["Alice"] = Elevator{true,1,1,0,lights, orders, queue}
		queue = [4]int{0, 0, 0, 0}
		elevators["Bob"] = Elevator{true,3,1,0,lights, orders, queue}
		message := Message{elevators, "Bob", "Alice", 1}
		for{
			UDPoutChan <- message
			time.Sleep(time.Millisecond*10)
		}
	}
	*/