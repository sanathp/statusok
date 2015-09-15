# StatusOK

Monitor Your Website and Apis from your computer.Get Notified through slack or email when your server is down or response time more than expected.


## Simple Setup

Simple Setup to monitor your website and recieve a notitification to Gmail when your website is down

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

Install Infulxdb 

write config file with influx db deatails as below

run statusok --config config.json


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

complete details on how to setup .


## Contribution

Feel free to Create pull requests.Write Test cases for the functionalities you have written.if you have written extension for some othe database . or written notification for any other client .

## License

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





