package read_config



import (
    "bufio"
    "log"
    "os"
    "strings"
    "strconv"
)

func ReadConfig() (string, int){
    file, err := os.Open("config.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        text := scanner.Text()



        if err := scanner.Err(); err != nil {
            log.Fatal(err)
        }
        start_id := strings.Split(text, " ")[0]
        d := strings.Split(text, " ")[1]
        depth, _ := strconv.Atoi(d)
        return start_id, depth
    }
    return "", 0
}
