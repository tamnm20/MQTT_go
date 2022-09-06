// Khai bao
// IP
var ip string = ""
// Temp
var temp string = ""
// Mem
var mem string = ""
// Cpu
var user int = 0
var system int = 0 
var idle int = 0
var avg float64 = 0

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 3, 64)
}

//get teamperture of PI
func Test_temp(){
	out, err := exec.Command("/opt/vc/bin/vcgencmd", "measure_temp").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
	
	out2 := out[0:(len(out)-1)]
	fmt.Printf("%s\n", out2)
	temp = string(out2)
	
}

func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)
    ip = localAddr.IP.String()
    return localAddr.IP
}


//get Memory of PI
func Get_Mem(){
	out, err := exec.Command("sh","-c","free -b | grep Mem").Output()
	if err != nil {
		log.Fatal(err)
	}
	var out1 = string(out)
	var list = strings.Split(out1, " ")

	total,_ := strconv.Atoi(list[6])
	 used,_ := strconv.Atoi(list[9])
	 var percent float64 = (float64(used)*100)/float64(total)

	 fmt.Printf("Phan tram used : %T\t%v\n",percent, percent)
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
	
	fmt.Printf("Old :  %v\t%v\t%v\n",user, system,idle)
	
	x,_ := strconv.Atoi(cpu[2])
	user = x - user
	y,_ := strconv.Atoi(cpu[4])
	system = y -system
	z,_ := strconv.Atoi(cpu[5])
	idle =  z- idle
	avg = ((float64(user +system))*100)/(float64(user + system + idle))
	
	fmt.Printf("New :  %v\t%v\t%v\n",x, y,z)
	fmt.Printf("Sub :  %v\t%v\t%v\n",user, system,idle)
	fmt.Printf("percent : %v\n\n", avg)
		
	user = x
	system = y
	idle = z
	
	}

// Push data Thingsboard
func Publish_data(){
	var payload string = "{" + "\"TEMP\":" + "\"" + temp + "\"" + "," + "IP" + ":" + ip + "," + "\"Memory Pi\":" + "\"" + mem + "\"" + "," + "\"CPU Pi\":" + "\"" + FloatToString(avg) + "\"" + "}" // Ghep các bản tin từ code cũ thành 1 bản tin 
	mqtt_demo_Thingsboard_Io.Publish(THINGSBOARD_TOPIC_IN, 0, false, payload)
	mqtt_amazon.Publish(THINGSBOARD_TOPIC_IN, 0, false, payload)
	
}



func main() {

	mqtt_begin()
	mqtt_domoticz.Publish(DOMITICZ_TOPIC_IN, 0, false,reset)	// reset buzzer ve off
	mqtt_domoticz.Subscribe(DOMITICZ_TOPIC_OUT, 0, mqtt_domoticz_messageHandler)  // QoS = 0
	mqtt_amazon.Subscribe(THINGSBOARD_TOPIC_OUT, 0,mqtt_messageHandler)
	
	// Loop:
	for {
		sleep_ms(5000) // delay 5s gửi dữu liệu 1 lần
		GetOutboundIP()
		Test_temp()
		Get_Mem()
		Get_cpu()
		Publish_data()
		
	}
}