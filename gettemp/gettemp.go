package main
import (
        "fmt"
        "log"
        "os/exec"
        "strconv"
)

var temp string = ""

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 3, 64)
}
func main() {
        out, err := exec.Command("sh","-c", "cat /sys/devices/virtual/thermal/thermal_zone0/temp").Output() // Lấy thông số Cpu temp
        if err != nil {
                log.Fatal(err)
        }
        out1,_ := strconv.Atoi(string(out[0:(len(out)-1)]))
        out3 := float64(out1)/1000
        fmt.Printf("%f\n", out3)
        temp = FloatToString(out3)
        payload := "{" + "\"CPU temp Pi\":" + "\"" + temp + "\"" + "}"
        fmt.Printf(payload)
        fmt.Printf("\n")
}