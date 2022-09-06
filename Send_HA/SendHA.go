package main
import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"strconv"
	"time"
)
// Khai bao
// IP
var ip string = ""
// Temp
var temp string = ""
// Mem
var mem string = ""
// Disk
var disk string = ""
// Cpu
var user int = 0
var system int = 0 
var idle int = 0
var avg float64 = 0

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 3, 64)
}

func sleep_ms(ms time.Duration) {
	time.Sleep(time.Millisecond * ms)   // sleep(nano second)
}

//get disk of PI
func get_disk(){
	out, err := exec.Command("sh","-c","df | grep /dev/mmcblk0p1").Output() // Lấy thông số Disk
	if err != nil {
		log.Fatal(err)
	}
	var out1 = string(out)
	var list = strings.Split(out1, " ")

//	fmt.Printf("list : %q\n", list)

	total,_ := strconv.Atoi(list[3])
	used,_ := strconv.Atoi(list[4])

//	fmt.Printf(" %T\t%v\n",total, total)
//	fmt.Printf(" %T\t%v\n",used, used)

	var percent float64 = (float64(used)*100)/float64(total)
//	fmt.Printf("Phan tram used : %T\t%v\n",percent, percent)

	disk = FloatToString(percent)
}

//get teamperture of PI
func get_temp(){
	out, err := exec.Command("sh","-c", "cat /sys/devices/virtual/thermal/thermal_zone0/temp").Output() // Lấy thông số Cpu temp
	if err != nil {
			log.Fatal(err)
	}
	out1,_ := strconv.Atoi(string(out[0:(len(out)-1)]))
	out3 := float64(out1)/1000
	// fmt.Printf("%f\n", out3)
	temp = FloatToString(out3)
}

//get IP of PI
func getIP(){
	out, err := exec.Command("sh","-c","ifconfig eth0| grep inet").Output() // Lấy thông số ip_pi
	if err != nil {
					log.Fatal(err)
	}
	var out1 = string(out)
	var list = strings.Split(out1, " ")

	// fmt.Printf("list : %q\n", list)

	// fmt.Printf(" %T\t%v\n",list[9], list[9])
	ip = list[9]
}

//get Memory of PI
func Get_Mem(){
	out, err := exec.Command("sh","-c","free -m | grep Mem").Output() // Lấy thông số Memory
	if err != nil {
			log.Fatal(err)
	}
	var out1 = string(out)
	var list = strings.Split(out1, " ")

	// fmt.Printf("list : %q\n", list)

	total,_ := strconv.Atoi(list[13])
	used,_ := strconv.Atoi(list[22])

	// fmt.Printf(" %T\t%v\n",total, total)
	// fmt.Printf(" %T\t%v\n",used, used)

	var percent float64 = (float64(used)*100)/float64(total)
	// fmt.Printf("Phan tram used : %T\t%v\n",percent, percent)

	mem = FloatToString(percent)
}

// Get cpu Pi
func Get_cpu(){
	out, err := exec.Command("sh","-c"," cat /proc/stat | grep 'cpu '").Output()
	if err != nil {
		log.Fatal(err)
	}
	var out1 = string(out)
	var cpu = strings.Split(out1, " ")
	
	// fmt.Printf("Old :  %v\t%v\t%v\n",user, system,idle)
	
	x,_ := strconv.Atoi(cpu[2])
	user = x - user
	y,_ := strconv.Atoi(cpu[4])
	system = y -system
	z,_ := strconv.Atoi(cpu[5])
	idle =  z- idle
	avg = ((float64(user +system))*100)/(float64(user + system + idle))
	
	// fmt.Printf("New :  %v\t%v\t%v\n",x, y,z)
	// fmt.Printf("Sub :  %v\t%v\t%v\n",user, system,idle)
	// fmt.Printf("percent : %v\n\n", avg)
		
	user = x
	system = y
	idle = z	
}

// Push data Thingsboard
func Publish_data(){
	var payload string = "{" + "\"TEMP\":" + "\"" + temp + "\"" + "," + "\"IP Pi\":" + "\"" + ip + "\"" + "," + "\"Memory Pi\":" + "\"" + mem + "\"" + "," + "\"CPU Pi\":" + "\"" + FloatToString(avg) + "\"" + "}" // Ghep các bản tin từ code cũ thành 1 bản tin
	fmt.Printf(payload) 
	fmt.Printf("\n")
	//mqtt_demo_Thingsboard_Io.Publish(THINGSBOARD_TOPIC_IN, 0, false, payload)
	//mqtt_amazon.Publish(THINGSBOARD_TOPIC_IN, 0, false, payload)
	
}



func main() {

	// mqtt_begin()
	// mqtt_domoticz.Publish(DOMITICZ_TOPIC_IN, 0, false,reset)	// reset buzzer ve off
	// mqtt_domoticz.Subscribe(DOMITICZ_TOPIC_OUT, 0, mqtt_domoticz_messageHandler)  // QoS = 0
	// mqtt_amazon.Subscribe(THINGSBOARD_TOPIC_OUT, 0,mqtt_messageHandler)
	
	// Loop:
	for {
		sleep_ms(1000) // delay 1s gửi dữu liệu 1 lần
		getIP()
		get_temp()
		get_disk()
		Get_Mem()
		Get_cpu()
		Publish_data()	
	}
}