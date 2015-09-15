# StatusOK

Monitor Your Website and Apis from your computer.Get Notified through slack or email when your server is down or response time more than expected.


## Simple Setup to monitor your website and recieve a notitification to Gmail when your website is down

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

Step 2: Download bin file from here and run the belwo command from your terminal
```
$ statusok --config config.json
```
Thats it !!!! you are done

To run as background process add & at the end

```
$ statusok --config config.json &	
```
to stop the process 
```
$ jobs
$ kill %jobnumber
```

## Complete Setup with InfluxDb and Grafanna :

![alt text](https://github.com/sanathp/StatusOK/raw/master/screenshots/graphana.png "Graphana Screenshot")


Install Infulxdb - url

Install Grafana - url

write config file with influx db deatails as below

run statusok --config config.json

open Grafana db and create dashboard from influxdb data as mentioned here -- url 
--write url own detailed explanation ?


## Database :

currently only influ db is supported .
if you want to monitor using your own database justw write a file.

### Write your own database client
	

## Notifications:

currenlty the below 4 types are supported .click them for more information
slack
mail
mailgun
http endpoint

### Write your own notification client






## Contribution

feel free to pull requests.

if you have written extension for some othe database . or written notification for any other client . you can generate a pull request .


## TODO





## License






