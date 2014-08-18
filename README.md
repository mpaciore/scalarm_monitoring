monitoring_deamon
=================
building guide:

for development environment (http):
  go install -a monitoring_deamon/manager
  
for production environment (https):
  go install -a -tags prod monitoring_deamon/manager

for production environment (https) without certificates checking (UNSAFE):
  go install -a -tags "prod certOff" monitoring_deamon/manager
