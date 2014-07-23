monitoring-deamon
=================
building guide:

for development environment (http):
  go build manager
  
for production environment (https):
  go build -tags prod manager

for production environment (https) without certificates checking:
  go build -tags 'prod certOff' manager
