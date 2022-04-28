/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    rpc
 * @Date:    2022/3/10 5:42 下午
 * @package: player
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package player

import (
	"fmt"
	"github.com/jageros/metactl/example/game/protos/meta"
	"github.com/jageros/metactl/example/game/protos/meta/sess"
	"github.com/jageros/metactl/example/game/protos/pb"
	"time"
)

func fetchConfig(ss sess.ISession) (resp *pb.Config, err error) {
	resp = &pb.Config{
		ServerTime: time.Now().Unix(),
		Version:    "v3.0",
		Channel:    "LoLo",
	}
	return
}

func login(ss sess.ISession, arg *pb.LoginArg) (resp *pb.LoginResp, err error) {
	fmt.Println(arg.String())
	resp = &pb.LoginResp{
		Uid:      1001,
		Avatar:   "https://www.xxx.com/avatar.jpg",
		Nickname: "Jager",
		Gold:     1000,
		Level:    99,
	}
	return
}

func playing(ss sess.ISession, arg *pb.PlayingArg) (err error) {
	fmt.Println(arg.String())
	return
}

func logout(ss sess.ISession) (err error) {
	return nil
}

func RegisterRpcHandle() {
	meta.C2S_FETCH_CONFIG.RegistryHandle(fetchConfig)
	meta.C2S_PLAYER_LOGIN.RegistryHandle(login)
	meta.C2S_PLAYER_PLAYING.RegistryHandle(playing)
	meta.C2S_PLAYER_LOGOUT.RegistryHandle(logout)
}
