## Writing a Config File

Config file shouuld be in JSON format (Support for other formats will be added in future).

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
[Click here](https://github.com/sanathp/StatusOK/blob/master/sample_config.json) to view the Sample config file.

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


| Paramter      | Descriition   | 
| ------------- |:-------------:| 
| url     | url of the api| 
| requestType     | Http Request Type in all capital letters  eg. GET PUT POST DELETE | 
| headers     | a list of key value pairs which will be added to header of a request| 
| formParams     | a list of key value pairs which will be added to body of the request . By deafult content type is application/x-www-form-urlencoded . for apllication/json content type add 	"Content-Type":"application/json" to headers| 
| urlParams     |  a list of key value pais which will be appended to url .eg: http://google.com?name= statusok| 
|checkEvery|time interval in seconds.If the value is 120,a request will be sent to given url every 2 minutes|
|responseCode|expected response code when a request is performed .Default values is 200.If response code is not equal then an error notification is triggered.|
|responseTime| expected response time in milliseconds. when mean response time is below this value a notification is triggered|


## Notifications 

Notifications will be triggered when mean response time is below given response time for a request or when an error is occured . Currently the below clients are supported to receive notifications.

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
To recieve notifications to any http Endpoint add below block to your config file with request details. Notification will be sent as a value to parameter "message"

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
### Write Your own Notification Client

If you want to recieve Notifications to any other clients. Write a struct with below methods and add the Struct to NotificationTypes in notify.go file .

```
	GetClientName() string
	Initialize() error
	SendResponseTimeNotification(notification ResponseTimeNotification) error
	SendErrorNotification(notification ErrorNotification) error
```
If you have written a new notification client which is useful to others, feel free to create a pull request.

## Database
 
Save Requests response time information and error information to your database by adding db details to config file. Currently only Influxdb 0.9.3+ is supported,adding support for any other database is very easy.

### Influx Db 0.9.3+

Install Influx db using the below commands.

```
wget http://influxdb.s3.amazonaws.com/influxdb_0.9.3_amd64.deb
dpkg -i influxdb_0.9.3_amd64.deb
/etc/init.d/influxdb start
```
Default username,password is empty and port is 8086.Add influxDb details as below inside database parameter to your config file.

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

![alt text](https://github.com/sanathp/StatusOK/raw/master/screenshots/graphana.png "Graphana Screenshot")

Run below commands to install grafana.

```
 wget https://grafanarel.s3.amazonaws.com/builds/grafana_2.1.3_amd64.deb
 apt-get update
 apt-get install -y adduser libfontconfig
 dpkg -i grafana_2.1.3_amd64.deb
 service grafana-server start

```
Create a new Dahsboard to view graphs as mentioned here http://docs.grafana.org/datasources/influxdb .

### Save Data to any other Database

Write a struct with below methods and add the Struct to DatabaseTypes in database.go file.

```
	Initialize() error
	GetDatabaseName() string
	AddRequestInfo(requestInfo RequestInfo) error
	AddErrorInfo(errorInfo ErrorInfo) error
```

If you have written structs to support any new database, feel free to create a pull request.