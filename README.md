Monitoring Daemon 
============ 
Contents 
---------- 
* scalarm_monitoring_daemon - main monitoring daemon program
* config - configuration for monitoring daemon

Installation guide: 
---------------------- 
Go 
-- 
To build and install monitoring daemon you need to install go programming language. 
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
git clone https://github.com/mpaciore/scalarm_monitoring_daemon
``` 
Now you can install monitoring: 
```` 
go install scalarm_monitoring_daemon 
```` 
This command will install monitoring daemon in $GOPATH/bin. It's name will be monitoringDaemon 
Build Options 
---------------- 
With -tags option you can specify build options:  
* no parameter: http requests 
* prod : https requests 
* certOff: disabling certificate checking for https 

Paramters can be mixed. For example: 
``` 
go install -tags "prod certOff" scalarm_monitoring_daemon
``` 
Note: Use -a option in go install if you didn't change any files after previous install. 
Config 
-------- 
The config folder contains single file config.json that contains required informations for monitor:

InformationServiceAddress - address of working Information Service
Login, Password - Scalarm credentials
Infrastructures - list of infrastructures monitor has to check for records


Run 
---- 
Before running program you have to copy contents of config folder to folder with executable of monitoring daemon. By default it will be $GOPATH/bin 

