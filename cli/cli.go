package cli

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func GetUserInput(q string, stdin io.Reader) string {
	fmt.Println(q)
	reader := bufio.NewReader(os.Stdin)
	userId, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(userId)
}

func Clr() {
	fmt.Print("\033[H\033[2J")
}
