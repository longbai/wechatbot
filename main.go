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
    //bot := openwechat.DefaultBot()
    bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

    // 注册消息处理函数
    bot.MessageHandler = handlers.Handler
    // 注册登陆二维码回调
    bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

    // 创建热存储容器对象
    reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
    // 执行热登录
    err := bot.HotLogin(reloadStorage)
    if err != nil {
        if err = bot.Login(); err != nil {
            log.Printf("login error: %v \n", err)
            return
        }
    }
    go WaitSignal(bot)
    // 阻塞主goroutine, 直到发生异常或者用户主动退出
    bot.Block()
}
