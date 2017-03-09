package driver

import (
	. "../definitions"
	//"bufio"
	"fmt"
	//"os"
	"time"
)

func elev_init() int { //Initilizes the elevator and the IO. Returns 0 if init fails. Returns the position the elevator finishes in otherwise.
	if Io_init() == 0 {
		return 0
	}

	elev_go(-1)
	for elev_check_floor_sensor() == 0 {
	}
	elev_go(0)

	return (elev_check_floor_sensor())
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

func elev_check_floor_sensor() int { //Returns the floor if the elevator is there, otherwise returns 0
	for i := 0; i < FLOORS; i++ {
		if Io_read_bit(SENSOR_FLOOR1+i) == 1 {
			elev_set_floor_light(i + 1)
			return i + 1
		}
	}
	return 0
}

func elev_stop_at_floor() { //1 sets the light, 0 clears it
	Io_set_bit(LIGHT_DOOR_OPEN)
	time.Sleep(time.Second * 1)
	Io_clear_bit(LIGHT_DOOR_OPEN)
}

func elev_go_to_floor(target chan int) { //Returns if the requested floor is out of range. Stops the elevator if something is written to kill
floor 			:= elev_check_floor_sensor()
current_target 	:= floor

	for{
		dummy := elev_check_floor_sensor()
		if dummy != 0{
			floor = dummy

		}
		select{
			case current_target = <- target:

			default:
				if current_target > FLOORS || current_target <1{
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
/*func elev_go_to_floor(floor int, kill chan bool, stopped_at_floor chan int) { //Returns if the requested floor is out of range. Stops the elevator if something is written to kill
	if floor > FLOORS || floor < 1 {
		return
	}
	var last_floor int
	for {
		temp := elev_check_floor_sensor()
		if temp != 0 {
			last_floor = temp
		}
		select {
		case <-kill:
			elev_go(0)
			return
		default:
			if last_floor == floor {
				elev_go(0)
				elev_stop_at_floor()
				stopped_at_floor <- floor
			}
			if last_floor > floor {
				elev_go(-1)
			}
			if last_floor < floor {
				elev_go(1)
			}

		}
	}
}*/

func elev_check_buttons(button_presses chan Orders) {
	button_inputs := Orders{}
	for {
		//Reads the internal orders
		button_inputs.IntOrders[0] = Io_read_bit(BUTTON_COMMAND1)
		button_inputs.IntOrders[1] = Io_read_bit(BUTTON_COMMAND2)
		button_inputs.IntOrders[2] = Io_read_bit(BUTTON_COMMAND3)
		button_inputs.IntOrders[3] = Io_read_bit(BUTTON_COMMAND4)

		//Reads external up orders
		button_inputs.ExtUpOrders[0] = Io_read_bit(BUTTON_UP1)
		button_inputs.ExtUpOrders[1] = Io_read_bit(BUTTON_UP2)
		button_inputs.ExtUpOrders[2] = Io_read_bit(BUTTON_UP3)

		//Reads external down orders
		button_inputs.ExtDwnOrders[1] = Io_read_bit(BUTTON_DOWN2)
		button_inputs.ExtDwnOrders[2] = Io_read_bit(BUTTON_DOWN3)
		button_inputs.ExtDwnOrders[3] = Io_read_bit(BUTTON_DOWN4)

		if button_inputs.IntOrders[0] == 1 || button_inputs.IntOrders[1] == 1 || button_inputs.IntOrders[2] == 1 || button_inputs.IntOrders[3] == 1 || button_inputs.ExtUpOrders[0] == 1 || button_inputs.ExtUpOrders[1] == 1 || button_inputs.ExtUpOrders[2] == 1 || button_inputs.ExtDwnOrders[1] == 1 || button_inputs.ExtDwnOrders[2] == 1 || button_inputs.ExtDwnOrders[3] == 1 {
			button_presses <- button_inputs
		}
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

func elev_check_motordir() int {

	return Io_read_bit(MOTORDIR)
}

func elev_status_checker(status chan Elevator) {
	status_elev 	:= Elevator{}
	status_change 	:= false

	for {
		floor := elev_check_floor_sensor()
		if floor != status_elev.Floor && floor != 0 {
			status_elev.Floor = floor
			status_change = true
		}

		dir := elev_check_motordir()
		if dir != status_elev.Direction {
			status_elev.Direction = dir
			status_change = true
		}

		if status_change {
			status <- status_elev
			status_change = false
		}
	}
}

func Elev_driver(incm_elev_update chan Elevator, out_elev_update chan Elevator) int {
	//---Create channels------------------------------
	target 			:= make(chan int)
	lights 			:= make(chan Orders)
	status 			:= make(chan Elevator)

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
			lights <- local_lift.Requests
		case lift_status := <- status:
			out_elev_update <- lift_status
		}
	}
}

/*
func Elev_driver(incm_elev_update chan Elevator, out_elev_update chan Elevator) int {
	//--------Init of driver-------------
	init_result := elev_init()
	if init_result == 0 {
		fmt.Println("Init failed")
		return 0 //The elevator failed to initialize
	}
	local_lift := Elevator{}
	last_target := 0

	kill := make(chan bool, 1)
	find_next_floor := make(chan int)
	lights := make(chan Orders)
	i := 0

	go elev_light_controller(lights)

	//----------Normal operation-----------
	for {
		select {
		case local_lift = <-incm_elev_update:
			i = 0
			kill <- true
			go elev_go_to_floor(local_lift.Queue[i], kill, find_next_floor)
			lights <- local_lift.Requests
			continue
		case <-find_next_floor:
			i++
		default:
			if i >= FLOORS {
				continue
			}
			if last_target != local_lift.Queue[i] {
				last_target = local_lift.Queue[i]
				go elev_go_to_floor(local_lift.Queue[i], kill, find_next_floor)
			}

		}
	}
}
*/

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

func Elev_test() {
	into_elev 	:= make(chan Elevator)
	outof_elev 	:= make(chan Elevator)


	go Elev_driver(into_elev, outof_elev)



	foo := Elevator{}

	foo.Queue= [4]int{4,1,1,1}	
	into_elev <- foo

	time.Sleep(time.Second*3)
	

	foo.Queue= [4]int{1,1,1,1}	
	into_elev <- foo

	for {}
	}
