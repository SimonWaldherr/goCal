# goCal

an experimental CalDav Server in Golang  
based on <https://github.com/samedi/caldav-go>

## install

install goCal with go get:

```bash
go get -u github.com/SimonWaldherr/goCal
```

## start

Before using goCal you need to check its configuration.
You can manage all user-accounts in the ```user.csv```-file.
To start goCal run this in your Terminal:

```bash
goCal -port=80 -sport=443 -sslcrt="path/to/your/sslcert.crt" -sslkey="path/to/your/sslkey.key" -user="path/to/your/user.csv" -storage="path/to/the/folder/where/the/ics-files/will/be/stored"
```

## use

now you can start using goCal, connect to your CalDAV-Server via https://localhost:443/ or subscribe to your ics-feed at https://localhost:443/icsfeed/ (or via webcal://localhost/).  

if you use a wordpress blog, you can install ```ICS Calendar``` and show your events at your blog.
