package udp

import (
	."../definitions"
	"./localip"
	"./bcast"
	"./peers"
	"time"
	"fmt"
	"reflect"
)


func UDPInit(UDPoutChan chan Message, UDPinChan chan Message, peerChan chan PeerUpdate) (localIP string) {
	fmt.Println("UDPinit")
	internetConnection := make(chan bool)
	go LocalMode(internetConnection)
	for{
		localIP, err := localip.LocalIP()
		if err != nil {
			continue
		} else {
			internetConnection <- true
			sendStatus(localIP)
			recieveStatus(peerChan)

			go transmitMessage(UDPoutChan, localIP)
			go recieveMessage(UDPinChan, localIP)
		
			return localIP
		}
		time.Sleep(time.Second)
	}	
}

func MasterInit(peerChan 		chan PeerUpdate, 	isMaster chan bool, 
				peerMasterChan 	chan PeerUpdate, 	localIP string, 
				UDPoutChan 		chan Message, 		masterIDChan chan string) (masterID string){
	fmt.Println("masterInit")
	select{
	case peerInfo := <- peerChan:
		companions := peerInfo.Peers
		if (companions[0] == localIP ){
			masterID = localIP
			isMaster <- true
			asterIDChan <- masterID
			peerMasterChan <- peerInfo
		} 
	}
	if(masterID == ""){
		go askPeersAboutMaster(peerChan, localIP, UDPoutChan)
		select{
		case masterID = <- masterIDChan:
			fmt.Println("masteridchan" + masterID)
			masterIDChan <- masterID
		}
	}
	return masterID
}

func UDPUpkeep(	peerChan 	chan PeerUpdate,	peerMasterChan 	chan PeerUpdate, 
				isMaster 	chan bool, 			masterIDChan 	chan string, 
				UDPoutChan 	chan Message,		masterID 		string, 
				localIP 	string){
	for{
		Upkeep:
			select {
			case peerInfo := <-peerChan:
				fmt.Println("Peerupdate")

				companions := peerInfo.Peers
				lostCompanions := peerInfo.Lost
				newCompanion := peerInfo.New
				for _, lostCompanion := range lostCompanions{
					if(lostCompanion == masterID){
						masterID = ""
						for _,companion := range companions{
							if (masterID == ""){
								masterID = companion
							}
							if (companion < masterID){
								masterID = companion
							}
						}
						if( masterID == localIP){
							isMaster <- true
							peerMasterChan <- peerInfo
						}
						masterIDChan <- masterID
					}
				}
				if(masterID == localIP && newCompanion != ""){
					go askAboutMaster(newCompanion, localIP, UDPoutChan)
					ticker := time.NewTicker(time.Second * 1).C 
					i := 0
					select{
					case otherMasterID := <- masterIDChan:
						if(otherMasterID < masterID){
							isMaster <- false
							masterID = otherMasterID
						}
					case <- ticker:
						go askAboutMaster(newCompanion, localIP, UDPoutChan)						//rebroadcasting if there is no reply
						i+=1
						if(i > 5){
							break Upkeep
						}
					}
				}
			}
	}
}

func askPeersAboutMaster(peerChan chan PeerUpdate, localIP string, UDPoutChan chan Message){
	select{
	case peerInfo := <- peerChan:
		companions := peerInfo.Peers
		for _, companion := range companions{
			askAboutMaster(companion, localIP, UDPoutChan)
		}
	}
}

func askAboutMaster(companion string, localIP string, UDPoutChan chan Message) {
	fmt.Println("I am asking if "  + companion + " is a master")
	msg := Message{}
	msg.MsgType = 3
	msg.SenderID = localIP
	msg.RecieverID = companion
	UDPoutChan <- msg
	fmt.Println("I am finished asking")
	return
}

func transmitMessage(UDPoutChan chan Message, localIP string){
	transmitChan := make(chan Message)
	echoChan := make(chan Message)
	go bcast.Transmitter(MESSAGEPORT, transmitChan)
	go bcast.Receiver(ECHOPORT, echoChan)
	for{
		select{
		case message := <- UDPoutChan:
			fmt.Println("Transmitting")
			message.SenderID = localIP 										//adding the localIP as senderID												
			transmitChan <- message 										//transmitting the mssage
			waitForEcho(transmitChan, echoChan, message)					//start new goroutine who waits for echo
		}
	}
}

func recieveMessage(UDPinChan chan Message, localIP string){
	recieveChan := make(chan Message)
	echoChan := make(chan Message)
	go bcast.Receiver(MESSAGEPORT, recieveChan)								//starting a receiver to recieve messages
	go bcast.Transmitter(ECHOPORT, echoChan)								//starting a transmitter to transmit echo
	for{
		select{
		case  message := <- recieveChan:
			fmt.Println("Recieved")
			if(message.RecieverID == localIP){								//checking to see if the message was ment for you
				echoChan <- message 										//putting out an echo on the echoport
				go func(){
					UDPinChan <- message 									//transmitting the message back to main and further
				}()
			}		
		}
	}
}

func waitForEcho(transmitChan chan Message, echoChan chan Message, message Message){

	ticker := time.NewTicker(time.Millisecond * 1000).C 					//waiting one second between resends
	i := 0
	for{
		select{								
		case <- ticker:
			fmt.Println("Echo")
			transmitChan <- message 										//rebroadcasting if there is no reply
			i+=1
			if(i > 5){
				return											//if no reply in five seconds assume peer lost and stop the  goroutine
				//HER MÅ OGSÅ MELDINGENE SENDES TILBAKE SÅ DE KAN BEHANDLS PÅ NYTT OG SENDES TIL NY RIKTIG PEER
			}
		case echo := <-echoChan:
			if(reflect.DeepEqual(echo.Elevators, message.Elevators) && echo.MsgType == message.MsgType){ 	//checking to see if you recieved the right echo
				fmt.Println("Right echo!")
				return																//when the right echo were recieved, stop the echo
			}
		}
	}
}

func sendStatus(localIP string){
	transmitStatus := make (chan bool)
	go peers.Transmitter(STATUSPORT, localIP, transmitStatus)
}

func recieveStatus(peerChan chan PeerUpdate){
	//i need another line
	go peers.Receiver(STATUSPORT, peerChan)
}