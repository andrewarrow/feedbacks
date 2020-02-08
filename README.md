ABOUT FEEDBACKS
==================

Feedbacks is a system for running N [feedback](https://github.com/andrewarrow/feedback)s

This main feedbacks executable is what you run, it will start up N other
feedbacks you specify in conf.toml.

The N other feedbacks are just simple http gin servers running on internal
ports only. This feedbacks executable runs on:

```
:443
:80
:2525
```

And handles https certs with [https://letsencrypt.org/](https://letsencrypt.org/)
and golang's autocert. i.e. you don't have to do anything but buy the domain.

It runs on 80 to just forward any non https request over to https.

It runs on 2525 to handle receiving email for all the domains you specify
in conf.toml. It sends each email thru [https://spamassassin.apache.org/](https://spamassassin.apache.org/) and records the email in the mysql database you specify
in conf.toml with the spam score.


EXAMPLE
==================

```
cp conf.toml.dist conf.toml
go build
brew install mysql
brew services start mysql
mysql -uroot

  CREATE USER 'dev'@'localhost' IDENTIFIED BY 'password'; 
  GRANT ALL ON *.* TO 'dev'@'localhost' WITH GRANT OPTION;

create database feedbacks;
quit

mysql -uroot feedbacks < migrations/first.sql
```

Edit your local copy of conf.toml:

```
[http]
hosts = [
  "cyborg.st",
  "many.pw",
  "jjaa.me"
]
```

In this example I'm using the 3 domains I own and telling feedbacks to
handle all the emails for all three, all the TLS certs, all the
hosting on 443 and 80. Each request that comes in will be handled by
the right feedback.

```
[paths]
sites = "./"
```

Change this path to where you have the code for these 3 feedbacks.

This means checkout a fresh copy of [https://github.com/andrewarrow/feedback](https://github.com/andrewarrow/feedback) and name it your domain. In my case I have:

```
./cyborg.st/
./many.pw/
./jjaa.me/
```

each one with a full copy of the code in https://github.com/andrewarrow/feedback.

Each one needs 

```
cp conf.toml.dist conf.toml
go build
```

and 

```
mysql -uroot
create database [name];
quit

mysql -uroot [name] < migrations/first.sql
```

Then you are off to the races making routes and controller and templates just
like you did in rails.

If you run feedbacks like:

```
LOCAL=jjaa.me ./feedbacks
```

It will run just on port 8080 to avoid https during local dev.

You can also just run 1 feedback on 1 port for local dev.
