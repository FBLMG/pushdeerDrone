package main

//引入包
import (
	"strings"
	"strconv"
	"net/http"
	"fmt"
	"time"
	"os"
)

//初始化参数
var pushKey string
var text string
var apiUrl string
var messageType string
var desp string
var repoName string
var buildNumber string
var buildStartTime string
var commitBranch string
var authorName string
var commitMessage string
var apiKey string

/**
程序初始化
 */
func init() {
	repoName = os.Getenv("DRONE_REPO_NAME")           //仓库-系统参数
	buildNumber = os.Getenv("DRONE_BUILD_NUMBER")     //构建次数-系统参数
	buildStartTime = os.Getenv("DRONE_STAGE_STARTED") //开始构建时间-系统参数
	commitBranch = os.Getenv("DRONE_COMMIT_BRANCH")   //提交分支-系统参数
	authorName = os.Getenv("DRONE_COMMIT_AUTHOR")     //提交者-系统参数
	commitMessage = os.Getenv("DRONE_COMMIT_MESSAGE") //提交信息-系统参数
	pushKey = os.Getenv("PLUGIN_PUSHKEY")             //推送Key-用户参数
	text = os.Getenv("PLUGIN_CONTENT")                //推送内容-用户参数
	apiUrl = os.Getenv("PLUGIN_URL")                  //推送链接-用户参数
	messageType = os.Getenv("PLUGIN_TYPE")            //推送类型-用户参数
	desp = os.Getenv("PLUGIN_DESP")                   //推送内容(markdown)-用户参数
	apiKey = os.Getenv("PLUGIN_API_KEY")              //推送内容-API密钥
	//处理数据
	buildStartTime = dealSystemTime() //时间转秒
	commitMessage = dealCommit()      //处理提交信息
	text = dealContent()              //初始化通知文案
}

/**
主函数
 */
func main() {
	//处理时间函数
	text = dealTime(text)
	//获取发送体
	payloadParams := dealMessageType(messageType, text, pushKey, desp)
	//组合参数
	payload := strings.NewReader(payloadParams)
	//调用
	_, err := http.Post(apiUrl, "application/x-www-form-urlencoded", payload)
	//判断是否调用失败
	if err != nil {
		fmt.Print(err)
		return
	}
	//输出成功
	fmt.Print("发送成功!!")
}

////////////////////////////////////////////处理数据////////////////////////////////////////////

/**
初始化通知文案
 */
func dealContent() string {
	//初始化默认文本
	content := "自动化部署成功\n\n" +
		"🏠仓库：" + repoName + "\n\n" +
		"⭕版本：" + buildNumber + "\n\n" +
		"🎅提交者：" + authorName + "\n\n" +
		"🕙耗时：" + buildStartTime + "\n\n" +
		"📖提交分支：" + commitBranch + "\n\n" +
		"📃提交信息：" + commitMessage + "\n\n"
	//判断是否需要自带文本
	if text == "" {
		text = content
	}
	//返回内容
	return text
}

/**
处理发送体
 */
func dealMessageType(messageType, text, pushKey, desp string) string {
	//初始化发送请求体
	sendMessage := ""
	//逻辑判断
	switch {
	case messageType == "text": //发送文字
		sendMessage = "pushkey=" + pushKey + "&text=" + text + "&type=text"
	case messageType == "image": //发送图片
		sendMessage = "pushkey=" + pushKey + "&text=" + text + "&type=image"
	case messageType == "markdown": //发送markdown
		sendMessage = "pushkey=" + pushKey + "&text=" + text + "&desp=" + desp + "&type=markdown"
	default: //发送文字
		sendMessage = "pushkey=" + pushKey + "&text=" + text + "&type=text"
	}
	//返回文本
	return sendMessage
}

/**
处理时间【处理成秒数】
 */
func dealTime(text string) string {
	//扫描结束时间特殊标
	startTimeIndex := strings.Index(text, "%0B")
	if startTimeIndex >= 0 {
		//提取结束时间
		startTime, _ := strconv.ParseInt(text[startTimeIndex+3:startTimeIndex+13], 10, 64)
		//获取当前时间戳
		endTime := time.Now().Unix()
		//获取秒
		seconds := endTime - startTime
		//拼接替换
		replaceString := "%0B" + text[startTimeIndex+3:startTimeIndex+13]
		//秒数转换
		secondsString := dealSeconds(seconds)
		//替换
		text = strings.Replace(text, replaceString, secondsString, -1)
	}
	//返回
	return text
}

/**
处理时间
 */
func dealSeconds(seconds int64) string {
	//秒数转换
	minute := seconds / 60
	second := seconds - minute*60
	//判断文案
	if second > 0 {
		return strconv.FormatInt(minute, 10) + "分" + strconv.FormatInt(second, 10) + "秒"
	} else {
		return strconv.FormatInt(minute, 10) + "分"
	}
}

/**
处理提交信息
 */
func dealCommit() string {
	return strings.Replace(commitMessage, "'", "", -1)
}

/**
处理时间【处理成秒数】
 */
func dealSystemTime() string {
	//构建时间首选系统参数-构建时间
	dealTimeValue := buildStartTime
	//字符串转int64
	startTime, _ := strconv.ParseInt(dealTimeValue, 10, 64)
	//获取当前时间戳
	endTime := time.Now().Unix()
	//获取秒
	seconds := endTime - startTime
	//秒数转换
	secondsString := dealSeconds(seconds)
	//返回
	return secondsString
}
