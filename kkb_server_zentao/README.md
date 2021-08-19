# 1.kkb_server_zentao部署及使用

## 1.1 部署zentao (docker)

```bash
mkdir -p /data/zbox && \
docker run -d -p 10030:80 -p 10031:3306 -p 10032:22 -p 10033-10039:10033-10039 \
        -e ADMINER_USER="root" -e ADMINER_PASSWD="010227" -e BIND_ADDRESS="false" \
        -v /data/zbox/:/opt/zbox/ \
        --add-host smtp.exmail.qq.com:163.177.90.125 \
        --name kkb-zentao-server \
        helenfrank/kkb_server_zentao:1.1
```

- `ADMINER_USER` : 设置web登录数据库Adminer账号。
- `ADMINER_PASSWD` : 设置web登录数据库Adminer密码。
- `BIND_ADDRESS`: 如果设置值为`false`，MySQL 服务器将不会绑定地址。
- `SMTP_HOST`: 设置smtp服务器IP和主机。（如果不能发送邮件，会有帮助。）也可以`extra_host`在docker-compose.yaml中使用，或者`--add-host`在使用`dokcer run`命令时使用param 。

注：禅道管理员账号为**admin**，默认初始化密码为**123456**。而MySQL root账号密码是**010227** （部署时要修改）

### 1.1.1 修改项

进入mysql , 开启远程访问

```bash
use mysql;
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY '010227' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```



## 1.2 kkb_server_zentao

### 1.2.1 项目结构图

```bash
kkb_server_zentao/
├── etc
│   ├── Account.json
│   ├── config_zentao.json
│   └── my.cnf
├── kkb_server_zentao
└── README.md
```



### 1.2.2 kkb_server_zentao

运行程序，挂在后台运行即可



## 1.2.3 etc下的配置文件

`Account.json`

用于生成`Token`，可在运行时添加

```json
[
    {
        "apiKey":"helen",
        "secretKey":"010227"
    },{
        "apiKey":"test1",
        "secretKey":"010227"
    }
]

```



`config_zentao.json`

连接zentao数据库的配置文件，修改后需要重启程序

```json
{
    "account": "root",
    "password": "010227",
    "ip": "127.0.0.1",
    "port": ":10031",
    "database": "zentao",
    "maxIdleConns": 10
}

```



`my.cnf`

暂时用来设置端口，更改后需要重启程序

```shell
[server]
port = 10227

```





# 2.API使用文档

## 2.1 获取token

### 2.1.1 请求方式

| 请求方式 | 地址                        |
| -------- | --------------------------- |
| GET      | http://127.0.0.1:10227/auth |

### 2.1.2 参数

| 参数名    | 解释     | 示例 |
| --------- | -------- | ---- |
| apiKey    | 接口密钥 |      |
| secretKey | 通用密钥 |      |

### 2.1.3 返回参数

`token` , `err`

主要根据返回的err来判断是否出错

err为空，表示无错

err携带错误信息，表示对应错误



## 2.2 添加zentao用户

### 2.2.1 请求方式

| 请求方式 | 地址                                         |
| -------- | -------------------------------------------- |
| POST     | http://127.0.0.1:10227/user/ZenTaoInsertUser |

### 2.2.2 参数

**json传递**

| 参数名       | 解释                                                    | 示例 |
| ------------ | ------------------------------------------------------- | ---- |
| realname     | 真实姓名                                                |      |
| mobileNumber | 电话号码                                                |      |
| email        | 邮箱                                                    |      |
| gender       | 性别（男，man为男，其他任意为女），可忽略不写，默认为女 |      |
| token        | 在http://127.0.0.1:10227/auth获取，失效时间2h           |      |
| account      | 账户 （暂定为电话号码，需要自行写参）                   |      |
| password     | 密码 （暂定为电话号码，需要自行写参）                   |      |

### 2.2.3 返回参数

`err` , `msg`(辅助说明)

主要根据返回的err来判断是否出错

err为空，表示无错

err携带错误信息，表示对应错误



## 2.3 添加zentao用户到项目

### 2.3.1 请求方式

| 请求方式 | 地址                                                   |
| -------- | ------------------------------------------------------ |
| POST     | http://127.0.0.1:10227/project/ZenTaoInsertUserProject |

### 2.3.2 参数

| 参数名  | 解释                                          | 示例                                                         |
| ------- | --------------------------------------------- | ------------------------------------------------------------ |
| root    | 项目id，在zentao后台获得                      | 2                                                            |
| days    | 任务天数                                      | 7                                                            |
| account | 用户账户，字符串数组形式                      | {"account":["18246340407","18612493872","13114550552","13579914479","13501152081"],"days":7,"root":2} |
| token   | 在http://127.0.0.1:10227/auth获取，失效时间2h |                                                              |

### 2.3.3 返回状态

`err` , `msg`(辅助说明)

主要根据返回的err来判断是否出错

err为空，表示无错

err携带错误信息，表示对应错误
