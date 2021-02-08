package socks

import (
	"context"
	"fmt"
	"sync"
	"xg/log"

	"github.com/gorilla/websocket"
)

type SocksManager struct {
	sync.Mutex
	clientMap map[int]*websocket.Conn
}

func (s *SocksManager) AddClient(ctx context.Context, userID int, conn *websocket.Conn) error {
	s.Lock()
	defer s.Unlock()
	//如果已连接，则关闭
	oldConn, ok := s.clientMap[userID]
	if ok {
		oldConn.Close()
	}

	s.clientMap[userID] = conn
	return nil
}

func (s *SocksManager) SendMessage(ctx context.Context, userID int, title string, msg string) error {
	conn, ok := s.clientMap[userID]
	if !ok {
		log.Info.Printf("No such client: %v", userID)
		return nil
	}
	data := fmt.Sprintf("[%v]%v", title, msg)
	return conn.WriteMessage(websocket.TextMessage, []byte(data))
}

var (
	_socketManagerOnce sync.Once
	_socketManager     *SocksManager
)

func GetSocksManager() *SocksManager {
	_socketManagerOnce.Do(func() {
		_socketManager = new(SocksManager)
	})
	return _socketManager
}
