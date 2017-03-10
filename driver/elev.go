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

func elev_check_floor_sensor1() int{
	return 1
}

func elev_go_to_floor(target chan int) { //Returns if the requested floor is out of range. Stops the elevator if something is written to kill
	

	floor := elev_check_floor_sensor1()
	current_target := floor

	for {
		dummy := elev_check_floor_sensor1()
		if dummy != 0 {
			floor = dummy

		}
		select {
		case current_target = <-target:

		default:
			time.Sleep(time.Millisecond * 10)
			if current_target > FLOORS || current_target < 1 {
				continue
			}
			if current_target == floor {
				elev_go(0)
				elev_stop_at_floor()
			}
			if current_target > floor {
				elev_go(1)
			}
			if current_target < floor {
				elev_go(-1)
			}

		}
	}
}

func elev_stop_at_floor() { //1 sets the light, 0 clears it
	Io_set_bit(LIGHT_DOOR_OPEN)
	//time.Sleep(time.Second * 3)
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
			
			time.Sleep(time.Millisecond*75)
		}
	}
}

func elev_check_buttons(button_presses chan Orders) {
	button_inputs := Orders{}
	dummy_inputs  := Orders{}

	for {
		//Reads the internal orders
		dummy_inputs.IntOrders[0] 		= Io_read_bit(BUTTON_COMMAND1)
		dummy_inputs.IntOrders[1] 		= Io_read_bit(BUTTON_COMMAND2)
		dummy_inputs.IntOrders[2] 		= Io_read_bit(BUTTON_COMMAND3)
		dummy_inputs.IntOrders[3] 		= Io_read_bit(BUTTON_COMMAND4)

		//Reads external up orders
		dummy_inputs.ExtUpOrders[0] 	= Io_read_bit(BUTTON_UP1)
		dummy_inputs.ExtUpOrders[1] 	= Io_read_bit(BUTTON_UP2)
		dummy_inputs.ExtUpOrders[2] 	= Io_read_bit(BUTTON_UP3)

		//Reads external down orders
		dummy_inputs.ExtDwnOrders[1] 	= Io_read_bit(BUTTON_DOWN2)
		dummy_inputs.ExtDwnOrders[2] 	= Io_read_bit(BUTTON_DOWN3)
		dummy_inputs.ExtDwnOrders[3] 	= Io_read_bit(BUTTON_DOWN4)

		if button_inputs != dummy_inputs && (dummy_inputs.IntOrders[0] == 1 || dummy_inputs.IntOrders[1] == 1 || dummy_inputs.IntOrders[2] == 1 || dummy_inputs.IntOrders[3] == 1 || dummy_inputs.ExtUpOrders[0] == 1 || dummy_inputs.ExtUpOrders[1] == 1 || dummy_inputs.ExtUpOrders[2] == 1 || dummy_inputs.ExtDwnOrders[1] == 1 || dummy_inputs.ExtDwnOrders[2] == 1 || dummy_inputs.ExtDwnOrders[3] == 1) {
			button_inputs = dummy_inputs
			fmt.Println(button_inputs)
			button_presses <- button_inputs	
		}
		dummy_inputs2 := Orders{}
		if dummy_inputs == dummy_inputs2{
			button_inputs = dummy_inputs
		}
		time.Sleep(time.Millisecond*10)
	}
}

func elev_light_controller(orders chan Orders) {
	for {
		select {
		case lights := <-orders:
			Io_write_bit(LIGHT_COMMAND1, lights.IntOrders[0])
			Io_write_bit(LIGHT_COMMAND3, lights.IntOrders[2])
			Io_write_bit(LIGHT_COMMAND2, lights.IntOrders[1])
			Io_write_bit(LIGHT_COMMAND4, lights.IntOrders[3])

			Io_write_bit(LIGHT_UP1, lights.ExtUpOrders[0])
			Io_write_bit(LIGHT_UP2, lights.ExtUpOrders[1])
			Io_write_bit(LIGHT_UP3, lights.ExtUpOrders[2])

			Io_write_bit(LIGHT_DOWN2, lights.ExtDwnOrders[1])
			Io_write_bit(LIGHT_DOWN3, lights.ExtDwnOrders[2])
			Io_write_bit(LIGHT_DOWN4, lights.ExtDwnOrders[3])
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
	buttons 		:= make(chan Orders)
	direction 		:= make(chan int)
	floor_sense 	:= make(chan int)

	go elev_poll_floor_sensor(floor_sense)
	go elev_check_motordir(direction)
	go elev_check_buttons(buttons)

	for {
		select {
			case presses 	:= <-buttons:
				status_elev.Order = presses
			case dir 		:= <- direction:
				status_elev.Direction = dir
			case 			   <-ticker:
			case floor   	:= <- floor_sense:
				status_elev.Floor = floor
			}
		status <- status_elev
		status_elev.Order = Orders{}
	}
}

func Elev_driver(incm_elev_update chan Elevator, out_elev_update chan Elevator) int {
	//---Create channels------------------------------
	target := make(chan int)
	lights := make(chan Orders)
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
			target <- local_lift.Queue[0]
			lights <- local_lift.Light
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
		break
		select{
		case dummy_elev:= <- outof_elev:
			dummy_elev2 := Elevator{}
			target := 1
			blip := true
			for i:= 0;i<FLOORS;i++{
				if dummy_elev.Order.IntOrders[i] == 1{
					target = i+1
					blip = false
				}
				if blip{
					break
				}
			}
			dummy_elev2.Queue[0] = target
			into_elev <- dummy_elev2
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
