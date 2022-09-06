package main
import (
        "fmt"
        "log"
        "os/exec"
        "strings"
        "strconv"
)

var mem string = ""

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 3, 64)
}
func main() {
	out, err := exec.Command("sh","-c","free -m | grep Mem").Output() // Lấy thông số Memory
	if err != nil {
			log.Fatal(err)
	}
	var out1 = string(out)
	var list = strings.Split(out1, " ")

	fmt.Printf("list : %q\n", list)

	total,_ := strconv.Atoi(list[13])
	used,_ := strconv.Atoi(list[22])

	fmt.Printf(" %T\t%v\n",total, total)
	fmt.Printf(" %T\t%v\n",used, used)

	var percent float64 = (float64(used)*100)/float64(total)
	fmt.Printf("Phan tram used : %T\t%v\n",percent, percent)

	mem = FloatToString(percent)
	payload := "{" + "\"Memory Pi\":" + "\"" + mem + "\"" + "}"
	fmt.Printf(payload)
	fmt.Printf("\n")
}
