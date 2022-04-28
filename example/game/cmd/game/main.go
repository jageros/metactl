/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    main
 * @Date:    2022/3/10 1:58 下午
 * @package: game
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package main

import (
	"context"
	"fmt"
	"github.com/jageros/metactl/example/game/internal/service/player"
	"github.com/jageros/metactl/example/game/internal/session"
	"github.com/jageros/metactl/example/game/protos/meta"
	"github.com/jageros/metactl/example/game/protos/pb"
	"gopkg.in/olahol/melody.v1"
	"log"
	"net"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	Port   = 8078
	GateID = 1000
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	m := melody.New()
	m.HandleConnect(func(s *melody.Session) {
		log.Printf("new connect uid=%d cid=%d", s.MustGet("uid"), s.MustGet("cid"))
	})
	m.HandleMessageBinary(handleMsg)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", Port),
		Handler: &httpServer{m: m},
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}
	player.RegisterRpcHandle()

	go func() {
		<-ctx.Done()
		ctx2, _ := context.WithTimeout(context.Background(), time.Second*5)
		s.Shutdown(ctx2)
	}()

	log.Println("Listening port: ", Port)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

type httpServer struct {
	m *melody.Melody
}

func ParesHeader(h http.Header) (int64, int64, error) {
	token := h.Get("X-Token")
	uid, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	var cid int64
	cliType := h.Get("X-ClientType")
	switch cliType {
	case "ios", "Ios", "IOS":
		cid = 1
	case "Android", "android", "ANDROID":
		cid = 2
	case "wechat", "Wechat", "WECHAT", "WeChat":
		cid = 3
	case "bytedance":
		cid = 4
	}
	return uid, cid, nil
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid, cid, err := ParesHeader(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("请求头部参数错误！")
		return
	}
	err = s.m.HandleRequestWithKeys(w, r, map[string]interface{}{"uid": uid, "cid": cid})
	if err != nil {
		log.Fatal(err)
	}
}

func handleMsg(ss *melody.Session, bytes []byte) {
	uid := ss.MustGet("uid")
	cid := ss.MustGet("cid")
	arg := &pb.PkgMsg{}
	err := arg.Unmarshal(bytes)
	if err != nil {
		log.Println(err)
		return
	}
	sess := session.New(uid.(int64), cid.(int64), GateID)

	resp, pbErr := onClientMsg(sess, arg)
	var reply = &pb.PkgMsg{
		Msgid: arg.Msgid,
	}
	if pbErr != nil {
		data, _ := pbErr.Marshal()
		reply.Type = pb.MsgType_Err
		reply.Payload = data
	} else if resp != nil {
		reply.Type = pb.MsgType_Reply
		reply.Payload = resp
	}

	if reply.Type != pb.MsgType_Unknown {
		data, _ := reply.Marshal()
		err = ss.WriteBinary(data)
		if err != nil {
			log.Println(err)
		}
	}
}

func onClientMsg(ss *session.Session, arg *pb.PkgMsg) ([]byte, *pb.ErrMsg) {
	var er = &pb.ErrMsg{
		Code: 200,
		Msg:  "successful",
	}

	resp, err := meta.Call(ss, arg.Msgid, arg.Payload)
	if err != nil {
		er.Code = -1
		er.Msg = err.Error()
		return nil, er
	}

	if resp != nil {
		return resp, nil
	} else if arg.Type == pb.MsgType_Req {
		return nil, er
	}
	return nil, nil
}
