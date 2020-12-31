

# StatusOK

statusOk는 http Request를 주기적으로 전송하고, 전송 결과를 데이터베이스에 저장하는 기능을 수행하는 서비스입니다.

# 간단한 예제
여기에 Request를 전송하고, 그 결과를 InfluxDb(시계열 데이터베이스)에 저장하는 간단한 예제가 있습니다. (config.json 파일)

```json

{
"database":{
    "influxDb":{
        "host":"influxdb",
        "port":8086,
        "databaseName":"statusok",
        "username":"user",
        "password":""
    }
},
    "requests":[
        {
            "url":"http://m.lalavla.com",
            "requestType":"GET",
            "checkEvery":10,    
            "responseTime":800
        }
    ]
}

```

이 예제는 m.lalavla.com에 단순한 GET 요청을 매 10초마다 전송하며, responseTime의 상한선을 800ms 로 설정합니다.
만약 m.lalavla.com 으로 부터 response가 800ms 이내로 오지 않는다면, statusOk는 이 요청을 실패로 간주합니다.

요청이 실패가 되면 설정에 따라 사용자에게 알림을 보낼 수 있습니다. 
또한, 요청에 헤더를 싣거나, formParameter, urlParameter를 실어 전송할 수도 있습니다.
(이제까지 언급한 기능들은 Original version 기능이며, Original version의 기능은 README_origin.md를 참조하세요.)


## 추가된 기능

운영환경에서 이 기능만으로 모니터링을 하기엔 충분하지 않습니다.
그렇기 때문에 httpResponse 코드나 responseTime 가 정상적 이더라도 ,
ResponseBody에 올바른 응답값이 전송되었는지 확인 할 수 있도록 새로운 기능을 추가하였습니다. 아래 예제를 보시죠.

```json

    {
        "database":{
            "influxDb":{
                "host":"influxdb",
                "port":8086,
                "databaseName":"statusok",
                "username":"user",
                "password":""
            }
        },
        "requests":[
            {
                "url":"https://m.lalavla.com/service/main/main.html",
                "requestType":"GET",
                "checkEvery":10,    
                "responseTime":800,
                "headers":{
                    "User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"
                }
                ,"FailedRequestAlertLevel" : -1
                ,"SaveBodyAlways" : "true"
                ,"Advanced":[
                    {
                        "checkType" : "contains",
                        "matchExpression" : "<script src=\"//jscdn.appier.net/aa.js?id=gsretail.com\" defer></script>",
                        "alertLevelRanges" : [
                            {
                                "from" : "0",
                                "to"   : "31",
                                "alertLevel" : "0"
                            },
                            {
                                "from" : "32",
                                "to" : "60",
                                "alertLevel" : "1"
                            }
                        ]
                    },
                    {
                        "checkType" : "contains",
                        "matchExpression" : "<script src=\"//jscdn.appier.net/aa.js?id=gsretail.com\"  22 defer></script>"
                    }
                ]
            }
        ]
    }


```

위 config 파일을 살펴보기 전에, 새롭게 추가된 개념을 알아둬야 할 필요가 있습니다.
그것은 바로 **"alertLevel"** 이라는 기능이며, 정수로 구현되는 이 level은 응답의 상태에 따라 높게, 또는 낮게 설정하여,
심각한 정도를 표현 할 수 있습니다. 이 **"alertLevel"** 을 사용하면, Grafana에서 임계값을 설정하기가 유용할 것 입니다.

위에서부터 차례대로 보겠습니다.
(첫 줄부터 "headers" 까지는 기존 original version 기능이므로 생략합니다.)

```json
 "FailedRequestAlertLevel" : -1
```
"FailedRequestAlertLevel" 은 request가 실패 했을때(설정된 responseTime을 초과하거나, responseCode가 실패이거나... 등) 설정될 alertLevel을 명시합니다.
생략할수 있으며, 생략하면 유사시 0이 세팅됩니다.
```json
 "SaveBodyAlways" : "true"
```
"SaveBodyAlways" 는 값이 "true"일 경우 request의 실패와 성공에 상관하지 않고 항상 responseBody값을 데이터베이스에 저장할지 결정합니다. 

```json
"Advanced" : [ ... ]
```
"Advanced" 는 배열로 구성되며, 여러 구성값들이 들어 갈 수 있습니다.
각 구성값들은 responseBody에 어떤 값이 들어가 있는가에 따라 alertLevel을 조정 할 수 있도록 구성되어있으며, 아래 예제와 같은 형태를 가져야합니다.
# ※참고 : "의 유무에 주의하여야 합니다.

```json
"Advanced" : [ 
                    {
                        "checkType" : "contains",
                        "matchExpression" : "매치할 문자열",
                        "alertLevelRanges" : [
                            {
                                "from" : "0",
                                "to"   : "31",
                                "alertLevel" : "0"
                            },
                            {
                                "from" : "32",
                                "to" : "60",
                                "alertLevel" : "1"
                            }
                        ]
                    },
                    {
                        "checkType" : "contains",
                        "matchExpression" : "매치할 문자열2",
                        "alertLevelRanges" : [
                            {
                                "from" : "0",
                                "to"   : "31",
                                "alertLevel" : "0"
                            },
                            {
                                "from" : "32",
                                "to" : "60",
                                "alertLevel" : "1"
                            }
                        ]
                    }
]
```
checkType은 반드시 "contains" 거나 "regex" 여야 하지만, "regex"는 현재 구현되지 않았으며, 향후 추가될 예정입니다.

위 예제는 , 응답 받은 responseBody를 체크하여, '매치할 문자열' 이 몇개나 들어가 있는지 카운트 한 다음,
이 카운트가 0~31이면 alertLevel을 0으로 설정,
이 카운트가 32~60이면 alertLevel을 1으로 설정합니다.

요청내에 여러 구성값들이 서로 같은 alertLevel을 설정하려고 할 수 있는데, 이러한 경우는 가장 높은 alertLevel이 최종적으로 설정됩니다.

## 사용법

1. 우선 docker와 docker compose를 설치합니다.

2. 원하는 디렉토리로 이동하여 statusok-master를 clone합니다. clone이 완료되면 clone을 시도한 디렉토리에 "ecommerceteam"이라는 디렉토리가 생성됩니다.
3. 다음 명령어를 통해 statusOk root 디렉토리로 이동합니다.

```bash
   cd ecommerceteam/statusok/statusok-master/
```

1. 이동된 디렉토리에서, docker-compose up -d 를 입력하면 자동으로 influxdb, grafana, statusok가 시작되며 연동됩니다. 만약 config.json을 수정하였다면, docker-compose up -d --build 와 같이 커맨드 마지막에 다시 빌드해야함을 명시해야합니다.

   모든 컨테이너가 작동되면 아래와 같은 메시지를 보게 될 것입니다.

```bash
Creating statusok-master_influxdb_1 ... done
Creating statusok-master_grafana_1  ... done
Creating statusok-master_statusok_1 ... done
```

5. 이제 브라우저에서 http:{ip주소}:3000 을치고, 그라파나 관리 콘솔로 접속할 수 있습니다.

docker 설치 : 
sudo yum install docker

docker-compose 설치 : https://docs.docker.com/compose/install/

## Grafana DataSource 등록 
Grafana DataSource등록 시, URL을 입력하는 부분에는 다음과 같이 입력합니다.

```
http://influxdb:8086
```

이는 docker-compose를 통해 influxdb 컨테이너의 주소를 "influxdb"로 정했기 때문입니다.(docker-compose.yml의 services 참조)
또한 database 명은 statusok로 입력하면 됩니다.(docker-compose.yml의 services/influxdb/environment : PRE_CREATE_DB 참조)

# Original StatusOK

# StatusOK

Monitor your Website and APIs from your computer.Get notified through Slack or E-mail when your server is down or response time is more than expected.


## Simple Version

Simple Setup to monitor your website and recieve a notification to your Gmail when your website is down.

Step 1: Write a config.json with the url information 
```
{
	"notifications":{
		"mail":{
			"smtpHost":"smtp.gmail.com",
			"port":587,
			"username":"yourmailid@gmail.com",
			"password":"your gmail password",
			"from":"yourmailid@gmail.com",
			"to":"destemailid@gmail.com"
		}
	},
	"requests":[
		{
			"url":"http://mywebsite.com",
			"requestType":"GET",
			"checkEvery":30,	
			"responseTime":800
		}
	]
}
```
Turn on access for your gmail https://www.google.com/settings/security/lesssecureapps .

Step 2: Download bin file from [here](https://github.com/sanathp/statusok/releases/) and run the below command from your terminal
```
$ ./statusok --config config.json
```
Thats it !!!! You will receive a mail when your website is down or response time is more.

To run as background process add & at the end

```
$ ./statusok --config config.json &	
```
to stop the process 
```
$ jobs
$ kill %jobnumber
```

## Complete Version using InfluxDb

![alt text](https://github.com/sanathp/StatusOK/raw/master/screenshots/graphana.png "Graphana Screenshot")

You can save data to influx db and view response times over a period of time as above using graphana.

[Guide to install influxdb and grafana](https://github.com/sanathp/statusok/blob/master/Config.md#database) 

With StatusOk you can monitor all your REST APIs by adding api details to config file as below.A Notification will be triggered when you api is down or response time is more than expected.

```json
{
	"url":"http://mywebsite.com/v1/data",
	"requestType":"POST",
	"headers":{
		"Authorization":"Bearer ac2168444f4de69c27d6384ea2ccf61a49669be5a2fb037ccc1f",
		"Content-Type":"application/json"
	},
	"formParams":{
		"description":"sanath test",
		"url":"http://google.com"
	},
	"checkEvery":30,
	"responseCode":200,		
	"responseTime":800
},

{
	"url":"http://mywebsite.com/v1/data",
	"requestType":"GET",
	"headers":{
		"Authorization":"Bearer ac2168444f4de69c27d6384ea2ccf61a49669be5a2fb037ccc1f",		
	},
	"urlParams":{
		"name":"statusok"
	},
	"checkEvery":300,
	"responseCode":200,		
	"responseTime":800
},

{
	"url":"http://something.com/v1/data",
	"requestType":"DELETE",
	"formParams":{
		"name":"statusok"
	},
	"checkEvery":300,
	"responseCode":200,		
	"responseTime":800
}

```
[Guide to write config.json file](https://github.com/sanathp/statusok/blob/master/Config.md#writing-a-config-file)

[Sample config.json file](https://github.com/sanathp/StatusOK/blob/master/sample_config.json)

To run the app

```
$ ./statusok --config config.json &
```

## Database

Save Requests response time information and error information to your database by adding database details to config file. Currently only Influxdb 0.9.3+ is supported.

You can also add data to your own database.[view details](https://github.com/sanathp/statusok/blob/master/Config.md#save-data-to-any-other-database)

## Notifications

Notifications will be triggered when mean response time is below given response time for a request or when an error is occured . Currently the below clients are supported to receive notifications.For more information on setup [click here](https://github.com/sanathp/statusok/blob/master/Config.md#notifications)

1. [Slack](https://github.com/sanathp/statusok/blob/master/Config.md#slack)
2. [Smtp Email](https://github.com/sanathp/statusok/blob/master/Config.md#e-mail)
3. [Mailgun](https://github.com/sanathp/statusok/blob/master/Config.md#mailgun)
4. [Http EndPoint](https://github.com/sanathp/statusok/blob/master/Config.md#http-endpoint)
5. [Dingding](https://github.com/sanathp/statusok/blob/master/Config.md#dingding)

Adding support to other clients is simple.[view details](https://github.com/sanathp/statusok/blob/master/Config.md#write-your-own-notification-client)

## Running with plain Docker

```
docker run -d -v /path/to/config/folder:/config sanathp/statusok
```

*Note*: Config folder should contain config file with name `config.json`

## Running with Docker Compose

Prepare docker-compose.yml config like this:

```
version: '2'
services:
  statusok:
    build: sanathp/statusok
    volumes:
      - /path/to/config/folder:/config
    depends_on:
      - influxdb
  influxdb:
    image: tutum/influxdb:0.9
    environment:
      - PRE_CREATE_DB="statusok" 
    ports:
      - 8083:8083 
      - 8086:8086
  grafana:
    image: grafana/grafana
    ports:
      - 3000:3000
```

Now run it:

```
docker-compose up
```

## Contribution

Contributions are welcomed and greatly appreciated. Create an issue if you find bugs.
Send a pull request if you have written a new feature or fixed an issue .Please make sure to write test cases.

## License
```
Copyright 2015 Sanath Kumar Pasumarthy

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License
```
