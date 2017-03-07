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


func UDPInit(UDPoutChan chan Message, UDPinChan chan Message, isMaster chan bool, masterIDChan chan string, peerChan chan peers.PeerUpdate) (localIP string) {
	fmt.Println("UDPinit")

	localIP, err := localip.LocalIP()
	if err != nil {
		return							//FIX HVA SOM SKAL SKJE VED FEIL
	}
	localIP = "Alice"

	sendStatus(localIP)
	recieveStatus(peerChan)

	go transmitMessage(UDPoutChan, localIP)
	go recieveMessage(UDPinChan, localIP)
	
	return localIP	
}

func MasterInit(peerChan chan peers.PeerUpdate, isMaster chan bool, localIP string, UDPoutChan chan Message) (masterID string){
	fmt.Println("masterInit")
	select{
	case peerInfo := <- peerChan:
		companions := peerInfo.Peers
		if (companions[0] == localIP ){
			masterID = localIP
			isMaster <- true
		} 
	}
	if(masterID == ""){
		go askPeersAboutMaster(peerChan, localIP, UDPoutChan)
	}
	return masterID
}

func UDPUpkeep(peerChan chan peers.PeerUpdate, peerMasterChan chan peers.PeerUpdate, isMaster chan bool, localIP string, masterIDChan chan string, masterID string, UDPoutChan chan Message){
	for{
		select {
		case peerInfo := <-peerChan:
			fmt.Println("Peerupdate")

			companions := peerInfo.Peers
			lostCompanions := peerInfo.Lost
			newCompanion := peerInfo.New
			for _, lostCompanion := range lostCompanions{
				fmt.Println(lostCompanion + masterID)
				if(lostCompanion == masterID){
					masterID = ""
					for _,companion := range companions{
					fmt.Println(companion)
						if (masterID == ""){
							masterID = companion
						}
						if (companion < masterID){
							masterID = companion
						}
					}
					if( masterID == localIP){
						isMaster <- true
					}
					masterIDChan <- masterID
				}
			}
			if(masterID == localIP && newCompanion != ""){
				go askAboutMaster(newCompanion, localIP, UDPoutChan)
				select{
				case otherMasterID := <- masterIDChan:
					if(otherMasterID < masterID){
						isMaster <- false
						masterID = otherMasterID
					}
				}
			}
		}
	}
}

func askPeersAboutMaster(peerChan chan peers.PeerUpdate, localIP string, UDPoutChan chan Message){
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
}

func transmitMessage(UDPoutChan chan Message, localIP string){
	transmitChan := make(chan Message)
	go bcast.Transmitter(MESSAGEPORT, transmitChan)
	for{
		select{
		case message := <- UDPoutChan:
			message.SenderID = localIP 										//adding the localIP as senderID			
			transmitChan <- message 										//transmitting the mssage
			go waitForEcho(transmitChan, message)							//start new goroutine who waits for echo
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
			if(message.RecieverID == localIP){								//checking to see if the message was ment for you
				echoChan <- message 										// putting out an echo on the echoport								
				UDPinChan <- message 										//transmitting the message back to main and further

			}
		}
	}
}

func waitForEcho(transmitChan chan Message, message Message){
	ticker := time.NewTicker(time.Millisecond * 1000).C 					//waiting one second between resends
	echoChan := make(chan Message)
	doneChan := make(chan bool)
	i := 0
	go bcast.Receiver(ECHOPORT, echoChan)
	for{
		select{
		case <- ticker:
			transmitChan <- message 										//rebroadcasting if there is no reply
			i+=1
			if(i > 5){
				doneChan <- true											//if no reply in five seconds assume peer lost and stop the  goroutine
				//HER MÅ OGSÅ MELDINGENE SENDES TILBAKE SÅ DE KAN BEHANDLS PÅ NYTT OG SENDES TIL NY RIKTIG PEER
			}
		case echo := <-echoChan:
			if(reflect.DeepEqual(echo.Elevators, message.Elevators) && echo.Order == message.Order && echo.MsgType == message.MsgType){ 											//checking to see if you recieved the right echo
				doneChan <- true 											
			}
		case <- doneChan:
			return															//when the right echo were recieved, stop the echo
		}
	}
}

func sendStatus(localIP string){
	transmitStatus := make (chan bool)
	go peers.Transmitter(STATUSPORT, localIP, transmitStatus)
}

func recieveStatus(peerChan chan peers.PeerUpdate){
	//i need another line
	go peers.Receiver(STATUSPORT, peerChan)
}