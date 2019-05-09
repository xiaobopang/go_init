package helper

/*
 * Created Date: Friday December 7th 2018
 * Author: Pangxiaobo
 * Last Modified: Friday December 7th 2018 6:27:55 pm
 * Modified By: the developer formerly known as Pangxiaobo at <10846295@qq.com>
 * Copyright (c) 2018 Pangxiaobo
 */

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"qiniupkg.com/x/log.v7"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numLetterBytes = "0123456789"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func Md5(str string) string {
	if str == "" {
		return ""
	}
	init := md5.New()
	init.Write([]byte(str))
	return fmt.Sprintf("%x", init.Sum(nil))
}

func Sha256(s string) string {
	if s == "" {
		return ""
	}
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Base64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// 字符串转换整型
func IpString2Int(ipstring string) int64 {
	ipSegs := strings.Split(ipstring, ".")
	ipInt := 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.Atoi(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return int64(ipInt)
}

// 整型转换成字符串
func IpInt2String(ipInt int) string {
	ipSegs := make([]string, 4)
	var len int = len(ipSegs)
	buffer := bytes.NewBufferString("")
	for i := 0; i < len; i++ {
		tempInt := ipInt & 0xFF
		ipSegs[len-i-1] = strconv.Itoa(tempInt)
		ipInt = ipInt >> 8
	}
	for i := 0; i < len; i++ {
		buffer.WriteString(ipSegs[i])
		if i < len-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}

func randString(n int, LetterBytes string) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(LetterBytes) {
			b[i] = LetterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func RandString(n int) string {
	return randString(n, letterBytes)
}

func RandNumString(n int) string {
	return randString(n, numLetterBytes)
}

//生成以日期随机的id
func GenDateRandId() string {
	return time.Now().Format("20060102150405") + RandNumString(6)
}

var Rander = rand.New(src)

// 随机数生成
// @Param	min 	int	最小值
// @Param 	max		int	最大值
// @return  int
func RandInt(min int, max int) int {
	num := Rander.Intn((max - min)) + min
	return num
}

// 获取服务器IP
func GetLocalIp() string {
	addrSlice, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Get local IP addr failed!")
		return "127.0.0.1"
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

// 是否是email
func IsEmail(email string) bool {
	if email == "" {
		return false
	}
	ok, _ := regexp.MatchString(`^([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+\.[0-9a-zA-Z]{2,3}$`, email)
	return ok
}

//时间转时间戳
func DateTimeToTimestamp(datetime string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai") //设置时区
	tt, err := time.ParseInLocation("2006-01-02 15:04:05", datetime, loc)
	if err != nil {
		log.Error("datetime to timestamp error: ", err)
		return 0
	}
	return tt.Unix()
}

//时间转时间戳
func DateToTimestamp(datetime string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai") //设置时区
	tt, err := time.ParseInLocation("2006-01-02", datetime, loc)
	if err != nil {
		log.Error("date to timestamp error: ", err)
		return 0
	}
	return tt.Unix()
}

//时间戳 to 时间
func TimestampToDate(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

//当前时间
func NowDatetime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 当前时间戳
func NowTimestamp() int64 {
	return time.Now().Unix()
}

// implode 类似php的 implode函数,将数组转为字符串
func Implode(originArr []string, symbol string) string {
	newStr := strings.Join(originArr, symbol)
	return newStr
}

//explode 类似php的 explode函数，将字符串转为数组
func Explode(originStr string, symbol string) []string {
	strArr := strings.Split(originStr, symbol)
	return strArr
}
