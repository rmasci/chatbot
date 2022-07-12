package main

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/aichaos/rivescript-go"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
	"time"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebSocketConnection struct {
	*websocket.Conn
}

// WS json response defines the response sent back from websocket
type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json "-"`
}

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConnection]string)

func (conf *Conf) Home(w http.ResponseWriter, r *http.Request) {
	vm := jet.VarMap{}
	host := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	vm.Set("host", host)
	vm.Set("botname", conf.Botname)
	err := renderPage(w, "home.jet", vm)
	if err != nil {
		fmt.Fprintln(w, err)
	}
}

// Upgrades connection to websocket
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	log.Println("Client", r.RemoteAddr, "Connected.")

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
	}
	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error, restarting: %v\n", r)
		}
	}()
	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err == nil {
			payload.Conn = *conn
			wsChan <- payload
		} else {
			log.Println("Error in loop", err)
			break
		}
		fmt.Printf(".")
	}
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("websocket err")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func ListenToWsChannel(conf Conf) error {
	var response WsJsonResponse
	bot := rivescript.New(rivescript.WithUTF8())
	bot.SetUnicodePunctuation(`[!?;]`)
	err := bot.LoadDirectory(conf.Rivedir)
	if err != nil {
		return err
	}
	bot.SortReplies()
	log.Println("rivescript initialized")
	for {
		e := <-wsChan
		switch e.Action {
		case "username":
			clients[e.Conn] = e.Username
			users := getUserList(conf.Botname)
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "left":
			response.Action = "list_users"
			delete(clients, e.Conn)
			response.ConnectedUsers = getUserList(conf.Botname)
			broadcastToAll(response)
		case "listusers":
			log.Println("List users")
			response.ConnectedUsers = getUserList(conf.Botname)
			broadcastToAll(response)
		case "load":
			response.Action = "list_users"
			response.ConnectedUsers = getUserList(conf.Botname)
			broadcastToAll(response)
		case "broadcast":
			// This is where messages are sent....
			response.Action = "broadcast"
			tn := time.Now().Format(time.Kitchen)
			//response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)
			response.Message = fmt.Sprintf(`<div class="list-group list-group-flush border-bottom scrollarea">
                    <div class="d-flex w-100">
                      <strong class="mb-1">%s %s:</strong>&nbsp;%s
                    </div>
                  </div>`, tn, e.Username, e.Message)
			msg := isChatbot(conf.Botname, e.Message)
			if msg != "" {
				var reply string
				var err error
				reply, err = getReply(bot, e.Username, msg)
				if err != nil {
					reply = fmt.Sprintf("%s", err)
				}
				//reply = strings.ReplaceAll(reply, "\n", "<br>")
				fmt.Println(reply)
				tn = time.Now().Format(time.Kitchen)
				response.Message = fmt.Sprintf(`%s<div class="list-group list-group-flush border-bottom scrollarea">
                  <div class="d-flex w-100">
                      <strong class="mb-1">%s %s:</strong>&nbsp;@%s,&nbsp;%s
                    `, response.Message, tn, conf.Botname, e.Username, reply) // response.Message,tn, chatbot, e.Username, reply)
			}
			broadcastToAll(response)
		}
	}
}

func getUserList(chatbot string) []string {
	var hasChatbot bool
	var userList []string
	for _, x := range clients {
		if x == "cloudgenie" {
			hasChatbot = true
		}
		if x != "" {
			userList = append(userList, x)
		}
	}
	if !hasChatbot {
		userList = append(userList, chatbot)
	}
	sort.Strings(userList)
	return userList
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}
	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
