package localLift


import (
	."../definitions"
	"../driver"
	"../master"
	"fmt"
	"time"
)


func LocalMode(	internetConnection chan bool, currentStateChan chan Elevator, 
				elevIn chan Elevator, elevOut chan Elevator, currentState Elevator) {
	elevators 		:= make(map[string]Elevator)
	localIP  		:= ""
	change1 		:= false
	change2			:= false

	elevators[localIP] = currentState
	elevIn <- currentState
	for{
		select{
		case elevator := <- elevOut:
			elevators[localIP] = elevator
			elevators, change1 = master.IsTheElevatorFinished(elevators, localIP)
			elevators, change2 = master.CalculateOptimalElevator(elevators, localIP)
			currentState = elevators[localIP]
			if(change1 || change2){
				go func(){
					fmt.Println(elevators)
					elevIn <- currentState
					}()
				change1, change2 = false, false
			}
		case <- internetConnection:
			currentStateChan <- elevators[localIP] 
			return
		}
	}
}

func Elev_driver(incm_elev_update chan Elevator, out_elev_update chan Elevator) int {
	//---Create channels------------------------------
	target 		:= make(chan int)
	lights 		:= make(chan Buttons)
	statusIn 	:= make(chan Elevator)
	statusOut 	:= make(chan Elevator)


	//---Init of driver-------------------------------
	init_result := driver.Elev_init(target,lights,statusIn,statusOut)
	if init_result == 0 {
		fmt.Println("Init failed")
		return 0 //The elevator failed to initialize
	}

	time.Sleep(time.Second)
	//---Normal operation-----------------------------
	for {
		select {
		case local_lift := <-incm_elev_update:
			target <- local_lift.Queue[0]
			lights <- local_lift.Light
			go func(){
				statusIn <- local_lift
			}()
		case lift_status := <-statusOut:
			out_elev_update <- lift_status        									
		}
	}
}

