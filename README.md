Scalarm Monitoring
============ 
Contents 
---------- 
* scalarm_monitoring - main monitoring program
* config - configuration for monitoring

Installation guide: 
---------------------- 
Go 
-- 
To build and install monitoring you need to install go programming language. 
You can install it from official binary distribution: 

https://golang.org/doc/install

or from source: 

https://golang.org/doc/install/source 

After that you have to specify your $GOPATH. Read more about it here: 

https://golang.org/doc/code.html#GOPATH 

Installation 
-------------- 
You can download it directly from GitHub. You have to download it into your $GOPATH/src folder 
``` 
git clone https://github.com/mpaciore/scalarm_monitoring
``` 
Now you can install monitoring: 
```` 
go install scalarm_monitoring 
```` 
This command will install monitoring in $GOPATH/bin. It's name will be scalarm_monitoring.

Config 
-------- 
The config folder contains single file config.json that contains required informations for monitor:

* InformationServiceAddress - address of working Information Service
* Login, Password - Scalarm credentials
* Infrastructures - list of infrastructures monitor has to check for records
* ScalarmCertificatePath - path to custom certificate (optional, by default looking in standard certificate directory)
* ScalarmScheme - http or https (default)

Example config:

```
{
	"InformationServiceAddress": "149.156.10.32:31034",
	"Login": "login",
	"Password": "password",
	"Infrastructures": 
	[
		"qsub",
		"qcg"
	],
	"ScalarmCertificatePath" : "cert.pem",
	"ScalarmScheme" : "https"
}
```
Example config can be found in config/config.json.

Run 
---- 
Before running program you have to copy contents of config folder to folder with executable of monitoring. By default it will be $GOPATH/bin 

