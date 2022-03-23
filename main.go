package main

import (
	"fmt"
	"net/url"

	"github.com/getlantern/systray"
	"github.com/kardianos/service"
	"github.com/skratchdot/open-golang/open"
	"github.com/yaling888/soulx/common"
	"github.com/yaling888/soulx/constant"
	"github.com/yaling888/soulx/icon"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTemplateIcon(icon.Icon, icon.Icon128)
	systray.SetTitle("SoulX")
	systray.SetTooltip("SoulX")

	go func() {
		var (
			hasInstalled = true
			hasStarted   = false
		)

		sStatus, err := soulService.Status()
		if err == service.ErrNotInstalled {
			hasInstalled = false
		}

		if sStatus == service.StatusRunning {
			hasStarted = true
		}

		mSwitch := systray.AddMenuItemCheckbox("开启Soul", "启动/停止Soul", hasStarted)
		mController := systray.AddMenuItem("控制台", "打开控制台")

		mFunc := systray.AddMenuItem("功能设置", "功能设置")
		mReload := mFunc.AddSubMenuItem("重载配置", "重新加载配置文件")
		mStartOnBoot := mFunc.AddSubMenuItemCheckbox("加入系统服务", "安装/卸载系统服务", hasInstalled)

		systray.AddSeparator()
		mAbout := systray.AddMenuItem("关于", "关于")
		mAbout.AddSubMenuItem("         SoulX", "SoulX").Disable()
		mAbout.AddSubMenuItem(fmt.Sprintf("Version：%s", constant.Version), "").Disable()
		mAbout.AddSubMenuItem(fmt.Sprintf("Date：%s", constant.BuildTime), "").Disable()

		systray.AddSeparator()
		mQuit := systray.AddMenuItem("退出", "退出SoulX")

		if !hasStarted {
			mController.Hide()
			mReload.Hide()
		}

		for {
			select {
			case <-mSwitch.ClickedCh:
				if !hasStarted {
					if !hasInstalled {
						err = installService()
						if err == nil {
							hasInstalled = true
							mStartOnBoot.Check()
						}
					}
					err = startService()
					if err == nil {
						hasStarted = true
						mSwitch.Check()
						mController.Show()
						mReload.Show()
					}
				} else {
					err = stopService()
					if err == nil {
						hasStarted = false
						mSwitch.Uncheck()
						mController.Hide()
						mReload.Hide()
					}
				}
			case <-mReload.ClickedCh:
				_ = reloadConfig()
			case <-mController.ClickedCh:
				c, err := common.Parse()
				if err != nil {
					continue
				}
				ctrUrl := &url.URL{
					Scheme: "http",
					Host:   c.ExternalController,
					Path:   "/ui/",
				}
				_ = open.Run(ctrUrl.String())
			case <-mStartOnBoot.ClickedCh:
				if !hasInstalled {
					err = installService()
					if err == nil {
						hasInstalled = true
						mStartOnBoot.Check()
					}
				} else {
					if hasInstalled {
						if hasStarted {
							err = stopService()
							if err == nil {
								hasStarted = false
								mSwitch.Uncheck()
								mController.Hide()
								mReload.Hide()
							}
						}

						err = uninstallService()
						if err == nil {
							hasInstalled = false
							mStartOnBoot.Uncheck()
						}
					}
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// clean up here
}
