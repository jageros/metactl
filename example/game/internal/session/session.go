/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    session
 * @Date:    2022/3/10 5:46 下午
 * @package: session
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package session

type Session struct {
	Uid      int64
	ClientId int64
	GateId   int64
}

func New(uid, clid, gid int64) *Session {
	return &Session{
		Uid:      uid,
		ClientId: clid,
		GateId:   gid,
	}
}

func (s *Session) GetUid() int64 {
	return s.Uid
}

func (s *Session) GetClientId() int64 {
	return s.ClientId
}

func (s *Session) GetGid() int64 {
	return s.GateId
}
