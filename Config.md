## Writing a Config File

Config file should be in JSON format (Support for other formats will be added in future).

## Pattern
```
{
"notifications":{
	//notification clients to send notifications
	"slack":{
		"channel":"#general",
		"username":"statusok",
		"channelWebhookURL":"https://hooks.slack.com/services/T09SF8/E5Tl7"
	}
},
"database":{
	//database client details to save data
	"influxDb":{
		"host":"localhost",
		"port":8086,
		"databaseName":"statusok",
		"username":"",
		"password":""
	}
},
"requests":[
		//an array of request objects	
		{
			"url":"https://google.com",
			"requestType":"GET",
			"checkEvery":30,
			"responseCode":200,		
			"responseTime":800
		}
		.....
	]
},
"notifyWhen":{
	"meanResponseCount":10 //A notification will be triggered if mean response time of last 10 requests is less than given response time. Default value is 5
},
"port":3215 //By default the server runs on port 7321.You can define your custom port number as below
"concurrency":2 //Max Number of requests that can be performed concurrently.Default value is 1.

}

```
[Click here](https://github.com/sanathp/StatusOK/blob/master/sample_config.json) to view the Sample config file. Scroll down for more details on requests,notifications and database setup.

## Requests

You can monitor all types of REST APIs or websites as below.

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

### Request  parameters 

Description for each request parameter.

| Parameter      | Description   
| ------------- |------------- 
| url     | Http Url 
| requestType     | Http Request Type in all capital letters  e.g. GET,PUT,POST,DELETE 
| headers     | A list of key value pairs which will be added to header of a request
| formParams     | A list of key value pairs which will be added to body of the request.By deafult content type is "application/x-www-form-urlencoded".For application/json content type add "Content-Type":"application/json" to headers
| urlParams     | A list of key value pairs which will be appended to url e.g: http://google.com?name=statusok
|checkEvery| Time interval in seconds.If the value is 120,the request will be performed every 2 minutes
|responseCode|Expected response code when a request is performed.Default values is 200.If response code is not equal then an error notification is triggered.
|responseTime|Expected response time in milliseconds,when mean response time is below this value a notification is triggered


## Notifications 

Notifications will be triggered when mean response time is below given response time for a request or when an error is occured.Currently the below clients are supported to receive notifications.

```
1)Slack
2)Smtp Server
3)Mailgun
4)Http EndPoint
5)Dingding
```
[Write your own client](https://github.com/sanathp/statusok/blob/master/Config.md#write-your-own-notification-client)

### Slack

To recieve notifications to your Slack Channel,add below block to your config file with your slack details

```
"slack":{
	"channel":"#ChannelName",
	"username":"your user name",
	"channelWebhookURL":"slack webhook Url"
}

```

### E-Mail
To recieve notifications to your email using smtp server,add below block to your config file with your smtp server details.

```
"mail":{
	"smtpHost":"smtp host name",
	"port":port-no,
	"username":"your mail username",
	"password":"your mail passwrd",
	"from":"from email id",
	"to":"to email id"
}

```
### Mailgun

To recieve notifications to your email Using Mailgun add below block to your config file with your mailgun details.

```
"mailGun":{
	"email":"your email id",
	"apiKey":"your api key",
	"domain":"domain name",
	"publicApiKey":"your publick api key"
}
```
### Http EndPoint
To recieve notifications to any http Endpoint add below block to your config file with request details. Notification will be sent as a value to parameter "message".

```
"httpEndPoint":{
	"url":"http://mywebsite.com",
	"requestType":"POST",
	"headers":{
		"Authorization":"Bearer ac2168444f4de69c27d6384ea2ccf61a49669be5a2fb037ccc1f",
		"Content-Type":"application/json"
	}
}
```	
### Dingding
To recieve notifications to any Dingding add below block to your config file with request details.

```
"dingding":{
        "url": "https://oapi.dingtalk.com/robot/send?access_token=3d21b0b12499ab54e74805",
        "requestType":"POST",
        "headers":{
            "Content-Type":"application/json"
        }
}
```
[Dingding Dev Document](https://open-doc.dingtalk.com/docs/doc.htm?spm=a219a.7629140.0.0.Tvbh61&treeId=257&articleId=105735&docType=1)

### Write Your own Notification Client

If you want to recieve Notifications to any other clients. Write a struct with below methods and add the Struct to NotificationTypes in [notify.go](https://github.com/sanathp/statusok/blob/master/notify/notify.go) file.

```
GetClientName() string
Initialize() error
SendResponseTimeNotification(notification ResponseTimeNotification) error
SendErrorNotification(notification ErrorNotification) error
```
If you have written a new notification client which is useful to others, feel free to create a pull request.

## Database
 
Save Requests response time information and error information to your database by adding database details to config file. Currently only Influxdb 0.9.3+ is supported.[Add support to your database](https://github.com/sanathp/statusok/blob/master/Config.md#save-data-to-any-other-database)

### Influx Db 0.9.3+

Install Influx db using the below commands.

```
wget http://influxdb.s3.amazonaws.com/influxdb_0.9.3_amd64.deb
dpkg -i influxdb_0.9.3_amd64.deb
/etc/init.d/influxdb start

More Details : https://influxdb.com/docs/v0.9/introduction/installation.html
```
Default username,password is empty and port number is 8086.Add influxDb details as below inside database parameter to your config file.

```
"influxDb":{
	"host":"localhost",
	"port":8086,
	"databaseName":"statusok",
	"username":"",
	"password":""
}
```

To visualize data in influxdb you need to install grafana.

Run below commands to install grafana.

```
 wget https://grafanarel.s3.amazonaws.com/builds/grafana_2.1.3_amd64.deb
 apt-get update
 apt-get install -y adduser libfontconfig
 dpkg -i grafana_2.1.3_amd64.deb
 service grafana-server start

```
Graphana will be running on port 3000 (http://localhost:3000)

Create a new Dahsboard to view graphs as mentioned here http://docs.grafana.org/datasources/influxdb .

![alt text](https://github.com/sanathp/StatusOK/raw/master/screenshots/graphana.png "Graphana Screenshot")

### Save Data to any other Database

Write a struct with below methods and add the Struct to DatabaseTypes in [database.go](https://github.com/sanathp/statusok/blob/master/database/database.go) file.

```
Initialize() error
GetDatabaseName() string
AddRequestInfo(requestInfo RequestInfo) error
AddErrorInfo(errorInfo ErrorInfo) error
```

If you have written structs to support any new database, feel free to create a pull request.

## Logs

By defualt logs are written to stdout in json format.If you want logs to be written to a file,give file path as mentioned below


```
$ statusok --config config.json --log logfilepath.log
```
