package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/eatmoreapple/openwechat"

    "github.com/longbai/wechatbot/handlers"
)

//graceful shutdown
func WaitSignal(bot *openwechat.Bot) {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
    <-sigCh
    log.Println("exit")
    bot.Exit()
}

func main() {
    bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

    // 注册消息处理函数
    bot.MessageHandler = handlers.Handler
    // 注册登陆二维码回调
    bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

    // 创建热存储容器对象
    reloadStorage := openwechat.NewFileHotReloadStorage("storage.bak")
    // 执行热登录
    err := bot.HotLogin(reloadStorage)
    if err != nil {
        if err = bot.Login(); err != nil {
            log.Printf("login error: %v \n", err)
            return
        }
    }

    // 获取登陆的用户
    self, err := bot.GetCurrentUser()
    if err != nil {
        log.Println(err)
        return
    }

    // 获取所有的好友
    friends, err := self.Friends()
    log.Println(friends, err)

    // 获取所有的群组
    groups, err := self.Groups()
    log.Println(groups, err)

    go WaitSignal(bot)
    // 阻塞主goroutine, 直到发生异常或者用户主动退出
    bot.Block()
}
