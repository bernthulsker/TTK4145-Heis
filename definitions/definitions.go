package definitions

const MOTOR_SPEED 	int = 2800

const ELEVATORS 	int = 3
const FLOORS 		int = 4
const MESSAGEPORT 	int = 20200
const ECHOPORT 		int = 20201
const STATUSPORT 	int = 20202

type Orders struct {
	IntOrders    	[FLOORS]int
	ExtUpOrders  	[FLOORS]int
	ExtDwnOrders 	[FLOORS]int //0 for external, 1 for internal
}

type Elevator struct {
	Alive     		bool
	Floor     		int //Last floor visited
	Position 		int
	Direction 		int
	Light	  		Orders
	Order     		Orders
	Queue     		[FLOORS]int //First element of list is current target of the elevator, 2nd element is next...
}

type Message struct {
	Elevators  		map[string]Elevator
	SenderID   		string
	RecieverID 		string
	MsgType    		int //Message identifier, 1 is input, 2 is queue,
}


type PeerUpdate struct {
	Peers 			[]string
	New   			string
	Lost  			[]string
}