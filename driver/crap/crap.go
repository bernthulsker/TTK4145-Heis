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