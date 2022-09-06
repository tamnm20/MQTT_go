package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"strconv"
)

var disk string = ""

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 3, 64)
}
func main() {
	out, err := exec.Command("sh","-c","df | grep /dev/mmcblk0p1").Output() // Lấy thông số Disk
	if err != nil {
		log.Fatal(err)
	}
	var out1 = string(out)
	var list = strings.Split(out1, " ")

	fmt.Printf("list : %q\n", list)

	total,_ := strconv.Atoi(list[3])
	used,_ := strconv.Atoi(list[4])

	fmt.Printf(" %T\t%v\n",total, total)
	fmt.Printf(" %T\t%v\n",used, used)

	var percent float64 = (float64(used)*100)/float64(total)
	fmt.Printf("Phan tram used : %T\t%v\n",percent, percent)

	disk = FloatToString(percent)
	payload := "{" + "\"Disk Pi\":" + "\"" + disk + "\"" + "}"
	fmt.Printf(payload)
	fmt.Printf("\n")
}