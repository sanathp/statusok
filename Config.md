# Writing a Config File

A config file should be written in json format (Other formats will be added in future).

## Pattern
```
{
"notification":{
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
	"meanResponseCount":5 
	//A notification will be triggered if mean response time of last 5 requests is less than given response time. Default value is 10
},
	//By default the server runs on port 7321.You can define your custom port number as below
"port":3215

}

```
Sample config file here --url

## Requests

You can monitor all types of REST apis or websites as below

```
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

url 		 : url of the api
requestType  : Http Request Type in all capital letters  eg. GET PUT POST DELETE 
headers 	 : a list of key value pairs which will be added to header of a request
formParams   : a list of key value pairs which will be added to body of the request . By deafult 
			  content type is application/x-www-form-urlencoded . for apllication/json content type add 	"Content-Type":"application/json" to headers
urlParams    : a list of key value pais which will be appended to url .
				 eg: http://google.com?name= statusok

checkEvery   : time in seconds .
responseCode : expected response code when a request is performed .Default values is 200.If response code is not equal then an error notification is triggered.
responseTime : expected response time in milliseconds. when mean response time is below this value a notification is triggered.


## Notifications 

Notifications will be triggered when mean response time is below given response time for a request or when an error is occured

To recieve notifications to your Slack Channel add below block to your config file with your mailgun details

```
"slack":{
	"channel":"#ChannelName",
	"username":"your user name",
	"channelWebhookURL":"slack webhook Url"
}

```

To recieve notifications to your email using smtp server add below block to your config file with your mailgun details

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

To recieve notifications to your email Using Mailgun add below block to your config file with your mailgun details

```
"mailGun":{
	"email":"your email id",
	"apiKey":"your api key",
	"domain":"domain name",
	"publicApiKey":"your publick api key"
}
```

To recieve notifications to any http Endpoint add below block to your config file with your mailgun details. Notification will be sent as a value to parameter "message"

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

To receive notifications to your custom client .write a .go file with the below methods to notify package.

```
	GetClientName() string
	Initialize() error
	SendResponseTimeNotification(notification ResponseTimeNotification) error
	SendErrorNotification(notification ErrorNotification) error
```

Add your object to NotificationTypes struct in notify.go file --url