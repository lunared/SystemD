package bot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"regexp"

	"../helpers"
)

//Bot contains the basic information that the bot needs to perform it's bot-like function
type Bot struct {
	config Config
	conn   net.Conn
}

//Botable is a basic interface for a bot
//Requires connection functions and handler functions
type Botable interface {
	connect()
	listen()
	joinChannels()

	pingHandler([]string)
	errorHandler([]string)
	noticeHandler([]string)
	comdHandler([]string)
	privHandler([]string)
	joinHandler([]string)
	quitHandler([]string)
}

//Config contains the configuration information for the bot
type Config struct {
	Host     string   `json:"host"`
	Port     string   `json:"port"`
	Password string   `json:"password"`
	Channels []string `json:"channels"`
	Nickname string   `json:"nickname"`
	Debug    bool     `json:"debug"`

	WeatherAPIKey string `json:"weather_api_key,omitempty"`
}

//These regular expressions are used to
//1) determine which handler to use and
//2) what peices of information should be pulled out of each message
var (
	rePing   = regexp.MustCompile(`PING :(.+)$`)
	reError  = regexp.MustCompile(`ERROR :(.+)$`)
	reNotice = regexp.MustCompile(`:[^ ]+ NOTICE [^ ]+ :(.+)$`)
	reComd   = regexp.MustCompile(`:[^ ]+ (\d\d\d) ([^:]+)(?: :(.+))?$`)
	rePriv   = regexp.MustCompile(`:([^!]+)!([^@]+)@([^ ]+) PRIVMSG ([^ ]+) :(.+)$`)
	reJoin   = regexp.MustCompile(`:([^!]+)!([^@]+)@([^ ]+) JOIN :([^ ]+)$`)
	reQuit   = regexp.MustCompile(`:([^!]+)[^ ]+ QUIT :(.+)$`)
)

//Run starts the bot
//This tells the bot to attempt to create a connection and to listen on that conenction
func Run() (b Bot) {
	b.connect()
	defer b.CloseConn()
	go b.listen()
	return
}

//CloseConn closes the conenction that this bot currently has
func (b Bot) CloseConn() {
	b.conn.Close()
}

//readConfig reads the config and populates the bot's config variable
func (b Bot) readConfig() {
	rawConfig, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(rawConfig, &b.config)
	if err != nil {
		panic(err)
	}
}

//connect creates the inital connection to the irc server
//Additionally it tries to join the channels in the config
func (b Bot) connect() net.Conn {
	connString := fmt.Sprintf("%s:%s", b.config.Host, b.config.Port)

	connection, err := net.Dial("tcp", connString)

	if err != nil {
		panic(err)
	}

	b.joinChannels()

	return connection
}

//output ...
//Output anything the server sends us and process it
func (b Bot) listen() {

	fmt.Println("Starting Listener")
	scanner := bufio.NewScanner(b.conn)

	find := func(reg *regexp.Regexp, msg string) []string {
		return reg.FindStringSubmatch(msg)
	}

	b.debug("Starting scanner")

	for scanner.Scan() {
		msg := scanner.Text()

		b.debug(string(msg) + "~")

		switch {
		case rePing.MatchString(msg):
			go b.pingHandler(find(rePing, msg))

		case reError.MatchString(msg):
			go b.errorHandler(find(reError, msg))

		case reNotice.MatchString(msg):
			go b.noticeHandler(find(reNotice, msg))

		case reComd.MatchString(msg):
			go b.comdHandler(find(reComd, msg))

		case rePriv.MatchString(msg):
			go b.privHandler(find(rePriv, msg))

		case reJoin.MatchString(msg):
			go b.joinHandler(find(reJoin, msg))

		case reQuit.MatchString(msg):
			go b.quitHandler(find(reQuit, msg))

		}
	}
}

//joinChannels loops through the channels listed in the config and sends a join message for each one
func (b Bot) joinChannels() {
	fmt.Printf("Joining channels %s", b.config.Channels)
	for _, c := range b.config.Channels {
		b.sendCommand("JOIN", c)
		fmt.Println("Joining channel " + c)
	}
}

//sendCommand is a convenience function that sends a command to the IRC server
func (b Bot) sendCommand(command string, text string) {
	fmt.Fprintf(b.conn, "%s %s\r\n", command, text)
	fmt.Printf("%s %s\n", command, text)
}

//sendPrivMsg is a convenience function that acts as a wrapper on top of sendCommand for PRIVMSG commands
func (b Bot) sendPrivMsg(channel string, text string) {
	b.sendCommand("PRIVMSG", fmt.Sprintf("%s :%s\r\n", channel, text))
}

//debug is a convenience function that wraps helpers.Debug passing in the bot's debug flag additionally
func (b Bot) debug(message string) {
	helpers.Debug(message, b.config.Debug)
}
