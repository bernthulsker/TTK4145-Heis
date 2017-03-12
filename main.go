package main

import (
	. "./definitions"
	"./localLift"
	"./udp"
	"./master"
	//"time"
	"fmt"
)
/*
At ordre plasserers i kø							<---- tight
Lys													<---- tight
kjører forbi fjerde									<---- tight
local mode											<---- tight (Kanskej nå? MÅ teste på datamaskin!)
andre etasje wut?									<---- tight
processing pairs?									
spre ordre ved master død
genrelle feilmeldinger her og der

concrurent map read write lol wut?					<----- dette kommer og går


*/





func main(){
	go stateMachine()

	select{}
}


func stateMachine(){
	//Make channels
	UDPoutChan 		:= make(chan Message)
	UDPinChan 		:= make(chan Message)
	masterMessage 	:= make(chan Message)
	peerChan 		:= make(chan PeerUpdate)
	peerMasterChan 	:= make(chan PeerUpdate)
	isMaster 		:= make(chan bool)
	masterIDChan 	:= make(chan string)
	elevOut 		:= make(chan Elevator)
	elevIn 			:= make(chan Elevator)
	internetConnect := make(chan bool)
	masterID 		:= ""
	localIP			:= ""
	state 			:= "Initializing"
	currentState 	:= Elevator{}
	initialized 	:= false
	go udp.CheckInternetConnection(internetConnect)

	for{
		StateMachine:
		switch state {

		case "Initializing":

			localIP = udp.UDPInit(UDPoutChan, UDPinChan, peerChan)
			if( localIP == ""){ 
				state = "No internet"
				break 
			}

			go master.MasterLoop(isMaster, masterMessage, peerMasterChan, UDPoutChan)
			go treatMessages(UDPinChan, UDPoutChan, masterMessage, masterIDChan, elevIn, elevOut, localIP)
			masterID = udp.MasterInit(peerChan, isMaster, peerMasterChan, localIP, UDPoutChan, masterIDChan)
			
			go localLift.Elev_driver(elevIn, elevOut)
			go udp.UDPUpkeep(peerChan, peerMasterChan, isMaster, masterIDChan, UDPoutChan, masterID, localIP)

			initialized = true

			state = "Normal operation"

		case "Normal operation":
			elevOut <- currentState
			messageBackup := Message{}
			messageBackup.Elevators = make(map[string]Elevator)
			for{
				select{
				case internet := <- internetConnect:
					if(!internet){
						state = "No internet"
						fmt.Println(state)
						break StateMachine
					}
				}
			}

		case "No internet":

			internetConnection := make(chan bool)
			currentStateChan := make(chan Elevator)
			currentState.Order = currentState.Light 

			go localLift.LocalMode(internetConnection, currentStateChan, currentState)
			for{
				select{
				case internet := <- internetConnect:
					if(internet){
						if(initialized) {state = "Normal operation"
						} else{ state = "Initializing" }
						internetConnection <- true
						select{
						case currentState := <- currentStateChan:
							currentState.Order = currentState.Light 
							break StateMachine
						}
					}
				}		
			}
		}
	}
}


func treatMessages(	UDPinChan 		chan Message, 	UDPoutChan 		chan Message, 
					masterMessage 	chan Message, 	masterIDChan 	chan string, 
					elevIn 			chan Elevator, 	elevOut 		chan Elevator, 
					localIP 		string){

	fmt.Println("Treat Messages")
	messageBackup := Message{}
	messageBackup.Elevators = make(map[string]Elevator)
	masterID := ""
	for{
		select{
		case messageBackup = <- UDPinChan:
			if (messageBackup.MsgType == 1 && localIP == masterID){
				fmt.Println("I got an order and my ID is " + localIP)
				masterMessage <- messageBackup
			} else if (messageBackup.MsgType == 2){
				elevIn <- messageBackup.Elevators[localIP]
			} else if (messageBackup.MsgType == 3){
				fmt.Println("Someone asked if " + localIP + " is master")
				master.AmIMaster(messageBackup, masterID, UDPoutChan, localIP)
			} else if (messageBackup.MsgType == 4){
				fmt.Println("I was told that " + messageBackup.SenderID + " is the master")
				masterIDChan <- messageBackup.SenderID
				masterID = messageBackup.SenderID
			}
		case masterID = <- masterIDChan:
			fmt.Println("I got a masterID")

		case elev_status := <- elevOut:
			messageBackup.Elevators[localIP] = elev_status
			messageBackup.MsgType = 1
			messageBackup.RecieverID = masterID
			UDPoutChan <- messageBackup                	
		}
	}
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