package route

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"xg/socks"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	allowOrigin := os.Getenv("allow_origin")
	allowOriginParts := strings.Split(allowOrigin, ",")
	fmt.Println("Allow Origin:", allowOriginParts)
	if len(allowOriginParts) < 1 {
		panic("invalid allow origin")
	}
	for i := range allowOriginParts {
		if origin == allowOriginParts[i] {
			return true
		}
	}
	return false
}

func (s *Server) registerSocks(c *gin.Context) {
	w := c.Writer
	r := c.Request
	user := s.getJWTUser(c)
	client, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	socks.GetSocksManager().AddClient(c.Request.Context(), user.UserId, client)

	// defer client.Close()
	// for {
	// mt, message, err := client.ReadMessage()
	// 	if err != nil {
	// 		log.Println("read:", err)
	// 		break
	// 	}
	// 	log.Printf("recv: %s", message)
	// 	err = client.WriteMessage(mt, message)
	// 	if err != nil {
	// 		log.Println("write:", err)
	// 		break
	// 	}
	// }
}

// func OpenSockets() {
// 	http.HandleFunc("/socks/register", register)
// }
