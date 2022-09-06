package main
import (
        "fmt"
        "log"
        "os/exec"
        "strings"
        "strconv"
        "time"
)
// Khai báo các biến .

// CPU PI
var user int = 0
var system int = 0
var idle int = 0
var cpu string = ""

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 3, 64) // Ham chuyen tu float sang string - Lấy 3 số sau dấu ","
}

func sleep_ms(ms time.Duration) {
        time.Sleep(time.Millisecond * ms)   // sleep(nano second)
}

func main() {
	for{
	sleep_ms(2000)
	out, err := exec.Command("sh","-c"," cat /proc/stat | grep 'cpu '").Output() // lấy thông số CPU của PI.
	if err != nil {
			log.Fatal(err)
	}
	var out1 = string(out)
	var cpu = strings.Split(out1, " ")
	fmt.Printf("Old :  %v\t%v\t%v\n",user, system,idle)
	x,_ := strconv.Atoi(cpu[2])
	user = x - user
	y,_ := strconv.Atoi(cpu[4])
	system = y -system
	z,_ := strconv.Atoi(cpu[5])
	idle =  z- idle
	var avg float64 = ((float64(user +system))*100)/(float64(user + system + idle)) // tính tỷ lệ % CPU hoạt động.
	fmt.Printf("New :  %v\t%v\t%v\n",x, y,z)
	fmt.Printf("Sub :  %v\t%v\t%v\n",user, system,idle)
	fmt.Printf("Persen : %v\n\n", avg)
	payload := "{" + "\"CPU Pi\":" + "\"" + FloatToString(avg) + "\"" + "}" // Tạo bản tin chuẩn
	fmt.Printf(" %T\t%v\n",payload, payload)
	user = x
	system = y
	idle = z
	}
}
