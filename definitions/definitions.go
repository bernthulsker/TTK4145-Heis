
package definitions

const MOTOR_SPEED 	int = 2800

const ELEVATORS 	int = 3
const FLOORS 		int = 4
const MESSAGEPORT 	int = 20200
const ECHOPORT 		int = 20201
const STATUSPORT 	int = 20202

type Buttons struct {
	IntButtons    	[FLOORS]int
	ExtUpButtons  	[FLOORS-1]int
	ExtDwnButtons 	[FLOORS-1]int 		
}

type Elevator struct {
	Floor     		int 				//Last floor visited
	Position  		int 				//SONDRE MÅ FÅ VITE HVOR HEISEN ER !!!!!!! "get some help"
	Direction 		int
	Light     		Buttons
	Order     		Buttons
	Queue     		[FLOORS]int 		//First element of list is current target of the elevator, 2nd element is next...
}

type Message struct {
	Elevators  		map[string]Elevator
	SenderID   		string
	RecieverID 		string
	MsgType    		int 				//Message identifier, 1 is input, 2 is queue,
}

type PeerUpdate struct {
	Peers 			[]string
	New   			string
	Lost  			[]string
}