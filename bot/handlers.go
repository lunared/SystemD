package bot

import (
	"fmt"
	"strings"
)

func (b Bot) pingHandler(args []string) {
	code := args[1]

	b.sendCommand("PONG", code)
	b.debug("Pinged and Ponged, code is " + code)
}

func (b Bot) errorHandler(args []string) {
	err := args[1]
	b.debug("~~Error occured: " + err)
}

func (b Bot) noticeHandler(args []string) {
	//text := args[1]
	//fmt.Println("~~~Notice Handler")
}

func (b Bot) privHandler(args []string) {
	nick := args[1]
	//user := args[2]
	//host := args[3]
	dest := args[4]
	mesg := args[5]

	b.debug(fmt.Sprintf("%v", args))

	responseDest := dest

	if dest == b.config.Nickname {
		responseDest = nick
	}

	switch {
	case strings.HasPrefix(mesg, "!echo"):
		mesg = strings.Replace(mesg, "!echo", "", 1)
		mesg = strings.TrimSpace(mesg)
		b.sendPrivMsg(responseDest, mesg)
	case strings.HasPrefix(mesg, "!stop"):
		b.CloseConn()

	case checkWeatherCommand(mesg):
		printWeatherData(b.conn, responseDest, mesg, b.config.WeatherAPIKey)

		//Removing database calls temporairly
		/*case strings.Contains(mesg, "!count"):
		count, err := getUpdateCount()
		if err != nil {
			sendPrivMsg(conn, responseDest, "Error connecting to database")
			fmt.Println(err)
		} else {
			sendPrivMsg(conn, responseDest, "Current count is: "+fmt.Sprint(count))
		}*/
	}

}

func (b Bot) comdHandler(args []string) {
	number := args[1]
	//who := args[2]
	//text := args[3]

	b.debug("~~~COMD Handler number is " + number)

	switch number {
	case "443":
		b.debug("~~Nick in use trying again~~")
	case "001":
		b.joinChannels()
		b.debug("~~~Join Handler")
	}
}

func (b Bot) joinHandler(args []string) {
	nick := args[1]
	//user := args[2]
	//host := args[3]
	channel := args[4]

	b.debug("~~~Join Handler")

	b.sendPrivMsg(channel, "Hi "+nick+"!")
}

func (b Bot) quitHandler(args []string) {
	user := args[1]
	reason := args[2]

	b.debug("User has left: " + user)

	if !strings.HasSuffix(user, "[m]") {
		if user == "irisu" {
			b.sendPrivMsg("#squad", "ciiiiiiiiiiim irisu died again: "+reason)
		} else {
			b.sendPrivMsg("#squad", "Bye "+user)
		}

	}

	b.debug("~~~Quit Handler")
}
