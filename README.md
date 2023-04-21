# wechatbot
> 最近chatGPT异常火爆，本项目可以将个人微信化身GPT机器人，
> 项目基于[openwechat](https://github.com/eatmoreapple/openwechat) 开发。

[![Release](https://img.shields.io/github/v/release/869413421/wechatbot.svg?style=flat-square)](https://github.com/869413421/wechatbot/releases/tag/v1.0.1)
![Github stars](https://img.shields.io/github/stars/869413421/wechatbot.svg)
![Forks](https://img.shields.io/github/forks/869413421/wechatbot.svg?style=flat-square)

### 目前实现了以下功能
 * 提问增加上下文，更接近官网效果 
 * 机器人群聊@回复
 * 机器人私聊回复
 * 好友添加自动通过
 
# 使用前提
> * ~~目前只支持在windows上运行因为需要弹窗扫码登录微信，后续会支持linux~~   已支持
> * 有openai账号，并且创建好api_key，注册事项可以参考[此文章](https://juejin.cn/post/7173447848292253704) 。
> * 微信必须实名认证。

# 注意事项
> * 项目仅供娱乐，滥用可能有微信封禁的风险，请勿用于商业用途。
> * 请注意收发敏感信息，本项目不做信息过滤。

# 快速开始
> 非技术人员请直接下载release中的[压缩包](https://github.com/869413421/wechatbot/releases/tag/v1.1.1) ，解压运行。
````
# 获取项目
git clone https://github.com/l2063610646/wechatbot.git

# 进入项目目录
cd wechatbot

# 复制配置文件
copy config.dev.json config.json

# 启动项目
go run main.go
````

# 配置文件说明
````
{
    "api_key": "your api key",
    "auto_pass": true,
    "session_timeout": 60,
    "use_proxy": false,
    "proxy_url": "http://127.0.0.1:10809",
    "use_gpt": false
}

api_key：openai api_key
auto_pass:是否自动通过好友添加
session_timeout：会话超时时间，默认60秒，单位秒，在会话时间内所有发送给机器人的信息会作为上下文。
use_proxy：是否使用代理模式（如果你没有国外服务器的时候有用）
proxy_url：代理的地址（例：http://127.0.0.1:10809）
use_gpt：是否使用gpt-3.5-turbo模型，否则使用text-davinci-003
````

# 使用示例
### 向机器人发送`过`，清空会话信息。

### 私聊
直接将想要提问的问题通过私聊的方式发送给机器人

### 群聊
在群里需要@机器人，然后再加上需要说的问题，例如：@我的小机器人 你好

# Docker部署方式
使用前需要创建一个config.json配置文件，该配置文件需要放在容器内的/app目录下，例如
````
docker run -it --name wechatgpt -v /opt/wechatgpt/config.json:/app/config.json wechatgpt-liu
````
登录之后将在容器内的/app路径下，生成storage.json文件，该文件记录了登录信息

如果出现第二次登录登录不上的问题，就可以尝试将该文件删除，所以我们使用数据卷挂载的方式也是个不错的选择

````
docker run -it --name wechatgpt -v /opt/wechatgpt:/app wechatgpt-liu 
````
这样就可以在宿主机通过删除/opt/wechatgpt/storage.json来使用清除登录信息的效果了
