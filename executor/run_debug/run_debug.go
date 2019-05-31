package main

import (
	"fmt"
	"github.com/inkedawn/JKWXFucker-server/utils"
	"github.com/inkedawn/go-sunshinemotion"
	"time"
)

var (
	s = ssmt.CreateSession()
	r = ssmt.Record{
		UserID:    6418,
		SchoolID:  60,
		Distance:  4.871,
		BeginTime: time.Date(2019, 5, 16, 7, 53, 33, 795044373, utils.TimeZoneBeijing), // 2019-05-16 07:53:33.795044373
		EndTime:   time.Date(2019, 5, 16, 9, 03, 26, 863371645, utils.TimeZoneBeijing), // 2019-05-16 09:03:26.863371645 +0800 CST m=+10403.853299344,
		IsValid:   true,
	}
)

func main() {
	_, err := s.Login(60, "041740431", "123", ssmt.PasswordHash("cc990814."))
	if err != nil {
		fmt.Printf("登录失败")
		panic(err)
	}
	s.Device.IMEI = "866841816672493"
	s.Device.IMSI = "866841816672493"
	err = s.UploadRecord(r)
	if err != nil {
		fmt.Printf("上传记录失败：%#v。 RecordDump: %s", err, utils.DumpStruct(r))
		panic(err)
	}
	fmt.Printf("上传记录成功。 RecordDump: %s", utils.DumpStruct(r))
}
