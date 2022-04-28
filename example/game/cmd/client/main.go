/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    main
 * @Date:    2022/3/11 10:07
 * @package: client
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jageros/metactl/example/game/protos/pb"
	"log"
	"math/rand"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	h := http.Header{}
	h.Set("X-Token", "800831")
	h.Set("X-ClientType", "ios")
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, "ws://127.0.0.1:8078", h)
	if err != nil {
		log.Fatal(1-1, err)
	}

	conn.SetPongHandler(func(appData string) error {
		log.Println("pong", appData)
		return nil
	})

	go func() {
		<-ctx.Done()
		err := conn.Close()
		if err != nil {
			log.Println(1-2, err)
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go ReadMsg(conn, wg)

	pkg := pb.PkgMsg{
		Type: pb.MsgType_Req,
	}
	tk := time.NewTicker(time.Second)
loop:
	for {
		select {
		case <-tk.C:
			i := rand.Intn(5)

			pkg.Msgid = pb.MsgID(100 + i)
			byt, _ := pkg.Marshal()
			err = conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
			if err != nil {
				log.Println(1-3, err)
			}
			err = conn.WriteMessage(websocket.BinaryMessage, byt)
			if err != nil {
				log.Println(1-4, err)
			}

		case <-ctx.Done():
			break loop
		}
	}
	wg.Wait()
}

func ReadMsg(conn *websocket.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	err := conn.SetReadDeadline(time.Now().Add(time.Second * 60))
	if err != nil {
		log.Println(2-1, err)
	}

	conn.SetPingHandler(func(appData string) error {
		log.Println("recv server ping msg.")
		err := conn.SetReadDeadline(time.Now().Add(time.Second * 60))
		if err != nil {
			log.Println(2-2, err)
		}
		return conn.WriteMessage(websocket.PongMessage, []byte{})
	})

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(2-4, err)
			return
		}
		pkg := &pb.PkgMsg{}
		err = pkg.Unmarshal(data)
		if err != nil {
			log.Println(2-5, err)
			continue
		}
		ParseRespMsg(pkg)
	}
}

func ParseRespMsg(pkg *pb.PkgMsg) {
	fmt.Printf("PkgType=%s MsgId=%s ", pkg.Type.String(), pkg.Msgid.String())
	if pkg.Type == pb.MsgType_Err {
		errMsg := &pb.ErrMsg{}
		err := errMsg.Unmarshal(pkg.Payload)
		if err != nil {
			fmt.Println()
			log.Println(3-1, err)
		} else {
			fmt.Printf("Err = {Code=%d Msg=%s}\n", errMsg.Code, errMsg.Msg)
		}
		return
	}
	switch pkg.Msgid {
	case pb.MsgID_C2S_FETCH_CONFIG:
		cfg := &pb.Config{}
		err := cfg.Unmarshal(pkg.Payload)
		if err != nil {
			fmt.Println()
			log.Println(3-1, err)
		} else {
			fmt.Printf("Resp = {version=%s channel=%s time=%s}\n", cfg.Version, cfg.Channel, time.Unix(cfg.ServerTime, 0).String())
		}
	case pb.MsgID_C2S_PLAYER_LOGIN:
		resp := &pb.LoginResp{}
		err := resp.Unmarshal(pkg.Payload)
		if err != nil {
			fmt.Println()
			log.Println(6, err)
		} else {
			fmt.Printf("Resp = {Uid=%d Avatar=%s Nickname=%s Gold=%d Level=%d}\n", resp.Uid, resp.Avatar, resp.Nickname, resp.Gold, resp.Level)
		}
	}
}
