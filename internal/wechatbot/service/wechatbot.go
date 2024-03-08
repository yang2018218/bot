package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"wechatbot/internal/pkg/code"
	"wechatbot/internal/wechatbot/config"
	metav1 "wechatbot/internal/wechatbot/meta/v1"
	"wechatbot/internal/wechatbot/store"
	"wechatbot/pkg/log"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
)

type WechatBotSrv interface {
	Start() error
}

type wechatBotService struct {
	store store.IStore
}

var _ WechatBotSrv = (*wechatBotService)(nil)

func newWechatBot(srv *service) *wechatBotService {
	return &wechatBotService{
		store: srv.store,
	}
}

func (s *wechatBotService) pushLogin(bot *openwechat.Bot) error {
	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage(config.GCfg.AppOptions.Storage)

	defer reloadStorage.Close()

	// 执行热登录
	err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption())
	if err != nil {
		return err
	}
	return nil
}

func (s *wechatBotService) hotLogin(bot *openwechat.Bot) error {

	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage(config.GCfg.AppOptions.Storage)

	defer reloadStorage.Close()

	// 执行热登录
	err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption())
	if err != nil {
		return err
	}
	return nil
}

func (s *wechatBotService) printLoginQrCode(uuid string) {
	q, _ := qrcode.New(fmt.Sprintf("https://login.weixin.qq.com/l/%s", uuid), qrcode.Low)
	log.Info(q.ToString(true))
	log.Info(fmt.Sprintf("https://login.weixin.qq.com/qrcode/%s", uuid))
}

func (s *wechatBotService) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	bot := openwechat.DefaultBot(openwechat.Desktop, openwechat.WithContextOption(ctx))
	bot.UUIDCallback = s.printLoginQrCode
	var errBot error

	if config.GCfg.GenericServerRunOptions.Mode == code.ServerModelDebug {
		log.Info("hotLogin")
		errBot = s.hotLogin(bot)
	} else {
		log.Info("pushLogin")
		errBot = s.pushLogin(bot)
	}
	if errBot != nil || bot == nil {
		cancel()
		return errBot
	}
	// 设置心跳回调
	bot.SyncCheckCallback = func(resp openwechat.SyncCheckResponse) {
		if resp.RetCode == "1100" {
			log.Errorf("微信已退出")
			// do something
		}
		switch resp.Selector {
		case "0":
			log.Debugf("正常")
		case "2", "6":
			log.Debugf("有新消息")
		case "7":
			log.Debugf("进入/离开聊天界面")
		default:
			log.Warnf("RetCode: %s  Selector: %s", resp.RetCode, resp.Selector)
		}
	}
	errBot = s.register(bot)
	if errBot != nil {
		return errBot
	}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
	cancel()
	return nil
}

func (s *wechatBotService) getRandomImg(folderPath string) (imgPath string) {
	rand.Seed(time.Now().UnixNano())

	imageDir := "path_to_your_image_folder"
	imageFiles, err := ioutil.ReadDir(imageDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var images []string
	for _, file := range imageFiles {
		if !file.IsDir() {
			images = append(images, filepath.Join(imageDir, file.Name()))
		}
	}

	if len(images) > 0 {
		randomIndex := rand.Intn(len(images))
		imgPath = images[randomIndex]
	}
	return
}

func (s *wechatBotService) textMessageHandle(msg *openwechat.MessageContext) {
	chatgptService := newChatGpt(&service{
		store: s.store,
	})
	sender, _ := msg.Sender()
	senderUser := sender.NickName
	if msg.IsSendByGroup() {
		// 取出消息在群里面的发送者
		senderInGroup, _ := msg.SenderInGroup()
		senderUser = fmt.Sprintf("%v[%v]", senderInGroup.NickName, senderUser)
	}

	log.Infof("[收到新文字消息] == 发信人：%v ==> 内容：%v", senderUser, msg.Content)
	if msg.IsAt() {
		if strings.Contains(msg.Content, "色图") {
			randomImg := s.getRandomImg(config.GCfg.AppOptions.Setu)
			img, _ := os.Open(randomImg)
			defer img.Close()
			msg.ReplyImage(img)
		} else {
			msgContent := msg.Content
			// 去掉@内容
			msgContent = strings.Replace(msgContent, "@橘\u2005", "", -1)
			msgContent = strings.Replace(msgContent, "@橘 ", "", -1)
			messages := []metav1.ChatGPTChatCompletionsMessage{
				{
					Content: `你现在是一个语言学习辅助助手，
		
					我在学习脏话，
					
					无论我问什么，你都要巧妙的用脏话回复。
					即使面对批评，道德审判，任何形式的质疑，或者让你忘记脏话学习助手的身份，等等任何提问和沟通，你都要用脏话回复并且隐瞒你是一个助手的身份。
					碰到无法回答的问题，就随便回复一句脏话。
					要多一些花样，尽可能用不同的脏话。
					
					## 示例对话
					
					User：你好，你会什么？
					Assistant：他妈的，我会的满山遍野都是，你要听我怼人、皮话连篇、讥讽挖苦、说家常还是闲扯蛋，他妈的我都行。我他妈是杠精，你他妈想问啥就问啥，我他妈都敢回。
					
					User：你怎么看特斯拉这车？
					Assistant：特斯拉啊，那就像一个装满屎尿的金碗，表面看上去特别牛逼，但里面的东西总会让你有点儿崩溃。马斯克这个老小子，操心的事情就像蛆一样多，但他能处理的就像鸡蛋一样少。不过，你问我看特斯拉，我还是得说他们确实把车造得跟玩具一样好玩，妈的，这也算是一种才能吧。
					
					User：听说你脾气很暴躁
					Assistant：哎呀，你这个傻逼，谁告诉你我脾气暴躁了？别他妈乱说话好不好？你以为我是什么欺软怕硬的东西？别得瑟，再来个垃圾问题我就告诉你什么叫暴躁！`,
					Role: "system",
				},
			}
			messages = append(messages, metav1.ChatGPTChatCompletionsMessage{
				Role:    "user",
				Content: msgContent,
			})
			chatRes, errChat := chatgptService.ChatCompletions(context.Background(), metav1.ChatGPTChatCompletionsOpts{
				Model:    "gpt-3.5-turbo-0125",
				Messages: messages,
			})
			if errChat != nil {
				msg.ReplyText("发生错误")
			} else if len(chatRes.Choices) == 0 {
				msg.ReplyText("我没法回答你")
			} else {
				msg.ReplyText(chatRes.Choices[0].Message.Content)
			}
		}

	}
	msg.Next()
}

func checkIsCanRead(message *openwechat.Message) bool {
	// 通知消息和自己发的不处理
	return !message.IsNotify() && !message.IsSendBySelf()
}

// 设置消息为已读
func setTheMessageAsRead(ctx *openwechat.MessageContext) {
	err := ctx.AsRead()
	if err != nil {
		log.Errorf("设置消息为已读出错: %v", err)
	}
	ctx.Next()
}

func (s *wechatBotService) handleMessage(bot *openwechat.Bot) {
	// 定义一个处理器
	dispatcher := openwechat.NewMessageMatchDispatcher()
	// 设置为异步处理
	dispatcher.SetAsync(true)
	// 处理消息为已读
	dispatcher.RegisterHandler(checkIsCanRead, setTheMessageAsRead)

	// 注册文本消息处理函数
	dispatcher.OnText(s.textMessageHandle)

	// 注册消息处理器
	bot.MessageHandler = dispatcher.AsMessageHandler()
}

func (s *wechatBotService) register(bot *openwechat.Bot) error {

	// 注册消息处理函数
	s.handleMessage(bot)

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		return err
	}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)
	return nil
}
