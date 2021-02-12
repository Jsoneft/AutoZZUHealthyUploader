# AutoZZUHealthyUploader

### Summary 

ZZU  健康上报自动上报脚本

利用 http框架  实现对一些接口的POST请求

利用gomail 框架 实现打卡邮件提醒服务

邮件截图如下

<img src="https://flysky-tencent-1302120781.cos.ap-chengdu.myqcloud.com/markdownImg/image-20210212143357204.png" alt="image-20210212143357204" style="zoom:50%;" />

## Features

- [x] 打卡群发邮件提醒服务

- [x] 日志

- [x] 多用户同时使用

- [x] 可自行定制 二次开发

## Requirements

详见 `go.mod`


## Usage

```sh
git clone https://github.com/Jsoneft/AutoZZUHealthyUploader.git
go run main.go
#  You can run the service in the background with this command
#  nohup go run main.go &
```



## Development

一台STMP服务器，用于自动收发邮件 (如果没有可以联系我 一条龙)

一台运行服务的服务器

## FAQ

*联系开发小哥哥  QQ: 1018437256*

## Support

走过路过给个星傲~

## Contact

1018437256@qq.com

## License

Mozilla Public License Version 2.0

## Attention

请保持，低调谨慎的态度使用脚本，本脚本隶属于Mozilla Public License，开发者免责。