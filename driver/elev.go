package driver

import (
	. "../definitions"
	"fmt"
	"time"
)

func elev_init() int { //Initilizes the elevator and the IO. Returns 0 if init fails. Returns 1 otherwise.
	if Io_init() == 0 {
		return 0
	}

	elev_go(-1)
	for Io_read_bit(SENSOR_FLOOR1) == 0 && Io_read_bit(SENSOR_FLOOR2)== 0 && Io_read_bit(SENSOR_FLOOR3)== 0 && Io_read_bit(SENSOR_FLOOR4) == 0 {
	}
	elev_go(0)

	return 1
}

func elev_go(dir int) { //Dir =1, elevator goes up. Dir = 0, elevator stops, Dir = -1 elevator goes down.
	if dir == -1 {
		Io_set_bit(MOTORDIR)
		Io_write_analog(MOTOR, MOTOR_SPEED)
	}
	if dir == 1 {
		Io_clear_bit(MOTORDIR)
		Io_write_analog(MOTOR, MOTOR_SPEED)
	}
	if dir == 0 {
		Io_write_analog(MOTOR, 0)
	}
}


func elev_go_to_floor(target chan int) { //Returns if the requested floor is out of range. Stops the elevator if something is written to kill
	floor := make(chan int)
	go elev_poll_floor_sensor(floor)
	current_target 	:= 	<- floor
	current_floor 	:= current_target
	has_stopped		:= true
	for{
		select{
		case current_target = <- target:
			
			if current_target == current_floor {
				elev_go(0)
				elev_stop_at_floor()
			}
			if current_target > current_floor {
				elev_go(1)
				has_stopped = false
			}
			if current_target < current_floor {
				elev_go(-1)
				has_stopped = false
			}
		case current_floor = <- floor:
				if current_target == current_floor {
				elev_go(0)
				if (!has_stopped){
					has_stopped = true
					elev_stop_at_floor()
				}
			}
		}
	}
}

func elev_stop_at_floor() { 
	fmt.Println("stopping")
	Io_set_bit(LIGHT_DOOR_OPEN)
	time.Sleep(time.Second * 3)
	Io_clear_bit(LIGHT_DOOR_OPEN)
}

func elev_set_floor_light(floor int) {
	if floor == 1 {
		Io_clear_bit(LIGHT_FLOOR_IND1)
		Io_clear_bit(LIGHT_FLOOR_IND2)
		return
	}
	if floor == 2 {
		Io_clear_bit(LIGHT_FLOOR_IND1)
		Io_set_bit(LIGHT_FLOOR_IND2)
		return
	}
	if floor == 3 {
		Io_set_bit(LIGHT_FLOOR_IND1)
		Io_clear_bit(LIGHT_FLOOR_IND2)
		return
	}
	if floor == 4 {
		Io_set_bit(LIGHT_FLOOR_IND1)
		Io_set_bit(LIGHT_FLOOR_IND2)
		return
	}
	return
}

func elev_poll_floor_sensor(floor_sense chan int){ //Returns the floor if the elevator is there, otherwise returns 0
	floor := -1

	for{
		time.Sleep(time.Millisecond*50)
		for i :=0 ; i< FLOORS; i++{
			dummy := Io_read_bit(SENSOR_FLOOR1 +i)
			if dummy == 1{
				floor = i+1
				floor_sense <- floor
				break
			}else if i == FLOORS -1{
				floor = 0
				floor_sense <- floor
				
			}
			
		}
	}
}

func elev_check_buttons(button_presses chan Buttons) {
	button_inputs := Buttons{}
	dummy_inputs  := Buttons{}

	for {
		//Reads the internal orders
		dummy_inputs.IntButtons[0] 		= Io_read_bit(BUTTON_COMMAND1)
		dummy_inputs.IntButtons[1] 		= Io_read_bit(BUTTON_COMMAND2)
		dummy_inputs.IntButtons[2] 		= Io_read_bit(BUTTON_COMMAND3)
		dummy_inputs.IntButtons[3] 		= Io_read_bit(BUTTON_COMMAND4)

		//Reads external up orders
		dummy_inputs.ExtUpButtons[0] 	= Io_read_bit(BUTTON_UP1)
		dummy_inputs.ExtUpButtons[1] 	= Io_read_bit(BUTTON_UP2)
		dummy_inputs.ExtUpButtons[2] 	= Io_read_bit(BUTTON_UP3)

		//Reads external down orders
		dummy_inputs.ExtDwnButtons[1] 	= Io_read_bit(BUTTON_DOWN2)
		dummy_inputs.ExtDwnButtons[2] 	= Io_read_bit(BUTTON_DOWN3)
		dummy_inputs.ExtDwnButtons[3] 	= Io_read_bit(BUTTON_DOWN4)
		
		

		if button_inputs != dummy_inputs {
		//if button_inputs != dummy_inputs && (dummy_inputs.IntOrders[0] == 1 || dummy_inputs.IntOrders[1] == 1 || dummy_inputs.IntOrders[2] == 1 || dummy_inputs.IntOrders[3] == 1 || dummy_inputs.ExtUpOrders[0] == 1 || dummy_inputs.ExtUpOrders[1] == 1 || dummy_inputs.ExtUpOrders[2] == 1 || dummy_inputs.ExtDwnOrders[1] == 1 || dummy_inputs.ExtDwnOrders[2] == 1 || dummy_inputs.ExtDwnOrders[3] == 1) {
			button_inputs = dummy_inputs
			button_presses <- button_inputs	
		}
		time.Sleep(time.Millisecond*50)
	}
}

func elev_light_controller(orders chan Buttons) {
	floor_light := make(chan int)
	go elev_poll_floor_sensor(floor_light)
	for {
		select {
		case lights := <-orders:
			Io_write_bit(LIGHT_COMMAND1, lights.IntButtons[0])
			Io_write_bit(LIGHT_COMMAND3, lights.IntButtons[2])
			Io_write_bit(LIGHT_COMMAND2, lights.IntButtons[1])
			Io_write_bit(LIGHT_COMMAND4, lights.IntButtons[3])

			Io_write_bit(LIGHT_UP1, lights.ExtUpButtons[0])
			Io_write_bit(LIGHT_UP2, lights.ExtUpButtons[1])
			Io_write_bit(LIGHT_UP3, lights.ExtUpButtons[2])

			Io_write_bit(LIGHT_DOWN2, lights.ExtDwnButtons[1])
			Io_write_bit(LIGHT_DOWN3, lights.ExtDwnButtons[2])
			Io_write_bit(LIGHT_DOWN4, lights.ExtDwnButtons[3])
		case floor := <- floor_light:
			if floor != 0{
				elev_set_floor_light(floor)
			}
		}
	}
}

func elev_check_motordir(dir chan int ){
	direction := -1 
	for{
		foo := Io_read_bit(MOTORDIR) 
		if foo != direction{
			direction = foo
			dir <-direction
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func elev_status_checker(status chan Elevator) {
	status_elev 	:= Elevator{}
	ticker 			:= time.NewTicker(time.Second).C
	buttons 		:= make(chan Buttons)
	direction 		:= make(chan int)
	floor_sense 	:= make(chan int)

	go elev_poll_floor_sensor(floor_sense)
	go elev_check_motordir(direction)
	go elev_check_buttons(buttons)

	for {
		select {
		case presses 	:= <-buttons:
			status_elev.Order = presses
			status <- status_elev
			status_elev.Order = Buttons{}
		case dir 		:= <- direction:
			status_elev.Direction = dir
			status <- status_elev
		case 			   <-ticker:
			status <- status_elev
		case floor   	:= <- floor_sense:
			status_elev.Floor = floor
			status <- status_elev
		case status_elev = <- status:
		}
	}
}

func Elev_driver(incm_elev_update chan Elevator, out_elev_update chan Elevator) int {
	//---Create channels------------------------------
	target := make(chan int)
	lights := make(chan Buttons)
	status := make(chan Elevator)


	//---Init of driver-------------------------------
	init_result := elev_init()
	if init_result == 0 {
		fmt.Println("Init failed")
		return 0 //The elevator failed to initialize
	}

	//---Start light controller and status checker----
	go elev_light_controller(lights)
	go elev_status_checker(status)
	go elev_go_to_floor(target)

	//---Normal operation-----------------------------
	for {
		select {
		case local_lift := <-incm_elev_update:
			fmt.Println("Recieving inncomming transmission")
			target <- local_lift.Queue[0]
			lights <- local_lift.Light
			status <- local_lift
		case lift_status := <-status:
			out_elev_update <- lift_status
		}
	}
}

func Elev_test() {
	into_elev := make(chan Elevator)
	outof_elev := make(chan Elevator)

	go Elev_driver(into_elev, outof_elev)

	for{
		select{
		case dummy_elev:= <- outof_elev:
			dummy_elev2 := Elevator{}

			for i,j := range dummy_elev.Order.IntOrders{
				if j == 1{
					dummy_elev2.Queue[0] = i+1
					break
				}
			}
			if (dummy_elev2.Queue[0] != 0){
				into_elev <- dummy_elev2	
			}

			}
		
		}
	select {}
	}

func pause() {
	t := time.Now()
	h, m, s := t.Clock()
	for {
		fmt.Println("Vi kommer snart tilbake. Vært borte siden:", h, m, s)
		sin := time.Since(t)
		fmt.Println("Som er så lenge siden:", sin)
		time.Sleep(time.Second * 1)

	}
}
