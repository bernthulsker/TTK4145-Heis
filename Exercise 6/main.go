package main
import(
	"io/ioutil"
	"time"
	"os"
	"os/exec"
	"fmt"
	"strconv"
	"strings"
)
const FILENAME = "heart.txt"

func main(){

	counter := watchDog()
	go heartbeat(counter)
	go startBackup()

	select{}
}


func watchDog() (int) {
	for{
		counter := 0;
		info,err:= os.Stat(FILENAME)
		if err != nil{
			fmt.Println(err.Error())
		}

		if time.Since(info.ModTime()) > (time.Millisecond*500){
			counterByte,_ := ioutil.ReadFile(FILENAME)
			counterStr := string(counterByte)
			counter,_ = strconv.Atoi(counterStr)
			counter++
			return counter
		}
		time.Sleep(time.Millisecond*100)
	}
}

func heartbeat(counter int ){
	tick := time.NewTicker(time.Millisecond*100).C
	for{
		select{
		case <-tick:
			s := strconv.Itoa(counter)
			b := []byte(s)
			ioutil.WriteFile(FILENAME, b, 0644)
			fmt.Println(s)
			counter++
		}
	}
}

func startBackup(){
	cmd := exec.Command("gnome-terminal", "-x","go", "run", "main.go")
	cmd.Run()
}

func convert( b []byte ) string {
    s := make([]string,len(b))
    for i := range b {
        s[i] = strconv.Itoa(int(b[i]))
    }
    return strings.Join(s,",")
}