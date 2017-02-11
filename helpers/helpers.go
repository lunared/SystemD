package helpers

import (
	"fmt"
	"net"
)

//Debug will print the message if debug is true, otherwise it's a no-op
func Debug(message string, debug bool) {
	if debug {
		fmt.Println(message)
	}
}

//SendCommand is a convenience function that sends a command to the IRC server
func SendCommand(conn net.Conn, command string, text string) {
	fmt.Fprintf(conn, "%s %s\r\n", command, text)
}

//SendPrivMsg is a convenience function that acts as a wrapper on top of SendCommand for PRIVMSG commands
func SendPrivMsg(conn net.Conn, channel string, text string) {
	SendCommand(conn, "PRIVMSG", fmt.Sprintf("%s :%s\r\n", channel, text))
}
