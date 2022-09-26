package main
import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"strconv"
	mqtt "github.com/eclipse/paho.mqtt.golang" 
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

var  mqtt_pione     mqtt.Client

const  MOSQUITTO_HOST string = "tcp://localhost:1883"
const  HA_TOPIC_IN  string = "xuong/server/specification"

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 3, 64)
}

func sleep_ms(ms time.Duration) {
	time.Sleep(time.Millisecond * ms)   // sleep(nano second)
}

func mqtt_begin() {
    opts_pione := mqtt.NewClientOptions()
    opts_pione.AddBroker(MOSQUITTO_HOST)
    opts_pione.SetUsername("nmtam")
    opts_pione.SetPassword("221220")
    opts_pione.SetCleanSession(true)
	opts_pione.SetOnConnectHandler(onConnectHandler)
    opts_pione.SetConnectionLostHandler(onConnectionLostHandler)
    mqtt_pione = mqtt.NewClient(opts_pione)
    if token := mqtt_pione.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    } else {
        fmt.Printf("MQTT pione Connected\n")
    }
}

//get disk of PI
func get_disk(){
	out, err := exec.Command("sh","-c","df | grep /dev/mmcblk0p1").Output() // Lấy thông số Disk
	if err != nil {
		log.Fatal(err)
	}
	var out1 = string(out)
	var list = strings.Split(out1, " ")
	total,_ := strconv.Atoi(list[3])
	used,_ := strconv.Atoi(list[4])
	var percent float64 = (float64(used)*100)/float64(total)
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
	total,_ := strconv.Atoi(list[13])
	used,_ := strconv.Atoi(list[22])
	var percent float64 = (float64(used)*100)/float64(total)
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
	x,_ := strconv.Atoi(cpu[2])
	user = x - user
	y,_ := strconv.Atoi(cpu[4])
	system = y -system
	z,_ := strconv.Atoi(cpu[5])
	idle =  z- idle
	avg = ((float64(user +system))*100)/(float64(user + system + idle))
	user = x
	system = y
	idle = z	
}

// Push data Thingsboard
func Publish_data(){
	var payload string = "{" + "\"CpuPiTemp\":" + "\"" + temp + "\"" + "," + "\"IpPi\":" + "\"" + ip + "\"" + "," + "\"DiskPi\":" + "\"" + disk + "\"" + "," + "\"MemoryPi\":" + "\"" + mem + "\"" + "," + "\"CpuPi\":" + "\"" + FloatToString(avg) + "\"" + "}" // Ghep các bản tin từ code cũ thành 1 bản tin
	fmt.Printf(payload) 
	fmt.Printf("\n")
	mqtt_pione.Publish(HA_TOPIC_IN, 0, false, payload)	// QoS = 0
}

func onConnectionLostHandler(mqtt_pione mqtt.Client, reason error) {
    fmt.Printf("Disconnected\n");
}

func onConnectHandler(mqtt_pione mqtt.Client) {
    fmt.Printf("Connected\n");
}

func main() {
	sleep_ms(180000) // Thời gian chờ HA và MQTT broker khởi động
	mqtt_begin()
	// Loop:
	for {
		sleep_ms(3000) // delay 3s gửi dữu liệu 1 lần
		getIP()
		get_temp()
		get_disk()
		Get_Mem()
		Get_cpu()
		Publish_data()	
	}
}