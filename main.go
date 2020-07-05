package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetSecID(url string) (secID string, error error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("Cookie", "_ga=GA1.2.685263550.1587277283; _gid=GA1.2.143250871.1587911549; tt_webid=6820028204934923790; _ba=BA0.2-20200301-5199e-c7q9NP0laGm7KfaPfGcH")
	req.Header.Add("status", "302")
	res, err := client.Do(req)
	if err == nil {
		var secIDRegexp = regexp.MustCompile(`sec_uid=(.*?)&`)
		secIDs := secIDRegexp.FindStringSubmatch(res.Request.URL.String())
		if len(secIDs) > 0 {
			secID = secIDs[1]
		}
	}
	return
}

type Info struct {
	UserInfo struct {
		UniqueId string `json:"unique_id"`
		TotalFavorited string `json:"total_favorited"`
		Nickname string `json:"nickname"`
		FollowerCount int `json:"follower_count"`
		AwemeCount int `json:"aweme_count"`
	} `json:"user_info"`
}

func GetDouyinID(secID string) (info Info, error error) {
	url := "https://www.iesdouyin.com/web/api/v2/user/info/?sec_uid=" + secID
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return info, err
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("Cookie", "_ga=GA1.2.685263550.1587277283; _gid=GA1.2.143250871.1587911549; tt_webid=6820028204934923790; _ba=BA0.2-20200301-5199e-c7q9NP0laGm7KfaPfGcH")
	req.Header.Add("status", "302")
	res, err := client.Do(req)
	if err == nil {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return info, err
		}

		err = json.Unmarshal(b, &info)
		if err != nil {
			return info, err
		}
	}
	return
}

func read3(path string) (string, error) {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return "", err
	}
	return string(fd), nil
}


func main()  {
	file, err := os.OpenFile("log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	info := log.New(io.MultiWriter(file, os.Stderr), "INFO: ", log.Ldate|log.Ltime)

	//读取账号
	user, err := read3("./user.txt")
	uids := strings.Split(user, "\r\n")
	for i := 0; i < len(uids); i++ {
		secID, err := GetSecID(uids[i])
		if err != nil {
			info.Println(err, user)
		} else {
			userInfo, err := GetDouyinID(secID)
			if err != nil {
				info.Println(err, user)
			} else {
				info.Println("昵称:" + userInfo.UserInfo.Nickname + " 作品数:" + strconv.Itoa(userInfo.UserInfo.AwemeCount) + " SecID:" + secID + " 抖音ID:" + userInfo.UserInfo.UniqueId)
			}
		}
	}
}
