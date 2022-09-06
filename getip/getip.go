package main
import (
        "fmt"
        "log"
        "os/exec"
        "strings"
        "strconv"
)

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 3, 64)
}
func main() {
        out, err := exec.Command("sh","-c","ifconfig eth0| grep inet").Output() // Lấy thông số ip_pi
        if err != nil {
                        log.Fatal(err)
        }
        var out1 = string(out)
        var list = strings.Split(out1, " ")

        fmt.Printf("list : %q\n", list)

        fmt.Printf(" %T\t%v\n",list[9], list[9])

        payload := "{" + "\"IP Pi\":" + "\"" + list[9] + "\"" + "}"
        fmt.Printf(payload)
        fmt.Printf("\n")
}
