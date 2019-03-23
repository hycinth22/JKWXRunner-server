// 提供对Session的智能管理，屏蔽底层细节，比直接操作token更为便捷。
//
// 实际并不存在Session表，通过对下层token、device等数个表的组合访问虚拟而来。
package sessionSrv

import (
	"github.com/inkedawn/JKWXFucker-server/database"
	"github.com/inkedawn/JKWXFucker-server/database/model"
	"github.com/inkedawn/JKWXFucker-server/service"
	"github.com/inkedawn/JKWXFucker-server/service/userIDRelationSrv"
	"time"

	"github.com/inkedawn/JKWXFucker-server/service/deviceSrv"
	"github.com/inkedawn/JKWXFucker-server/service/userCacheSrv"
	"github.com/inkedawn/go-sunshinemotion"
)

const PhoneNum = "123"

// 如果Session库中没有该帐号的Session， 则创建Session并保存后返回。
// 如果Session库已有该帐号的Session， 则检查Token的过期时间，
// 未过期则直接返回保存的Session， 已过期则更新Session
//
// 注意，该函数只检查Token的过期时间，并不会实际发送请求来验证Token有效性。
// 如果返回的Session包含失效Token，需要手动调用NewSession来完成更新。
func SmartGetSession(db *database.DB, acc model.Account) (s *ssmt.Session, err error) {
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	device, err := deviceSrv.GetDevice(tx, acc.DeviceID)
	if err != nil {
		return nil, err
	}
	SSMTDevice := deviceSrv.ToSSMTDevice(device)

	userToken, err := getTokenByUID(tx, acc.ID)
	if err != nil {
		if err == ErrNoToken {
			// 新用户登录
			return newSession(tx, acc, SSMTDevice)
		} else {
			return nil, err
		}
	}
	SSMTToken := toSSMTToken(userToken)
	if !tokenNotExpired(userToken) {
		// 过期更新
		return newSession(tx, acc, SSMTDevice)
	}

	// Resume Session
	s = ssmt.CreateSession()
	s.Device, s.Token = &SSMTDevice, &SSMTToken
	s.User = &ssmt.UserIdentify{
		UserID:   userToken.RemoteUserID,
		SchoolID: acc.SchoolID,
		Username: acc.StuNum,
	}
	return s, nil
}

// 创建一个Session并保存到Session库。不管是否已有该账号的Session
// 自动从service/device获取该Account的Device
//
// 返回的error可以直接与SSMT提供error比较
func NewSession(db *database.DB, acc model.Account) (s *ssmt.Session, err error) {
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	device, err := deviceSrv.GetDevice(tx, acc.DeviceID)
	if err != nil {
		return nil, err
	}
	SSMTDevice := deviceSrv.ToSSMTDevice(device)
	s, err = newSession(tx, acc, SSMTDevice)
	if err != nil {
		return nil, err
	}
	return s, err
}

func newSession(db *database.DB, acc model.Account, SSMTDevice ssmt.Device) (*ssmt.Session, error) {
	s := ssmt.CreateSession()
	s.Device = &SSMTDevice
	info, err := s.Login(acc.SchoolID, acc.StuNum, PhoneNum, ssmt.PasswordHash(acc.Password))
	if err != nil {
		return nil, err
	}

	// save into session storage
	err = saveToken(db, fromSSMTToken(s.User.UserID, *s.Token))
	if err != nil {
		return nil, service.WrapAsInternalError(err)
	}
	// update userInfo cache
	err = userCacheSrv.SaveCacheUserInfo(db, userCacheSrv.FromSSMTUserInfo(info, s.User.UserID, time.Now()))
	if err != nil {
		return nil, err
	}
	// update userID relation
	err = userIDRelationSrv.SaveRelation(db, acc.ID, s.User.UserID)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// 注意，该函数只检查Token的过期时间，并不会实际发送请求来验证Token有效性。
func tokenNotExpired(userToken token) bool {
	return time.Now().Before(userToken.ExpirationTime)
}
