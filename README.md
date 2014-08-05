monitoring-deamon
=================
building guide:

for development environment (http):
  go install -a manager
  
for production environment (https):
  go install -a -tags prod manager

for production environment (https) without certificates checking (UNSAFE):
  go install -a -tags "prod certOff" manager
