# go-paster

使用 GO 编写的粘贴器。用来解决部分窗口无法粘贴的问题。  
采用按键模拟的方式进行输入，直接调用 Windows SendInput API，并特别使用了 KEYEVENTF_SCANCODE 标志。这意味着它告诉操作系统“键盘上的第 X 号物理按键被按下了”，而不是“输入了字符 A”。这对于 KVM 和远程桌面至关重要。  

目前只支持 Windows

## 使用方法

> 若机器无 OpenGL 或在 RDP 下使用，双击无反应，可能需要下载 `opengl32.dll`，参考[下载](#下载)章节

1. 启动 GO-Paster
2. 粘贴所需要输入的内容入输入框
3. 调整 `Start Delay` 值为点击 `Start Typing` 按钮后延迟输入的时间。
4. 调整 `Interval` 值，为输入每个字符间的时间间隔
5. 点击 `Start Typing` 开始输入，此后程序进入倒计时，倒计时为 `Start Delay` 的时间，默认 5s。
6. 在倒计时内鼠标点击要输入的文本框，等待倒计时结束后自动输入。

## 下载

[Release](https://github.com/Rougnt/go-paster/releases)

> 由于 Fyne 依赖 OpenGL2.0 因此在部分场景，如无显卡机器可能需要软件模拟 OpenGL。  
> 可使用 Federico Dossena 构建的 [Mesa3D](https://fdossena.com/?p=mesa/index.frag)，解压其中的 `opengl32.dll` 到软件同级目录。
> 该 opengl32.dll 文件已在 Release 中一并提供 (x64 版本)。

## 编译

``` bash
# go build -o go-paster
# 在 mac 上编译 windows 文件
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o go-paster.exe 
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags -H=windowsgui -o go-paster.exe
```
