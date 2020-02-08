FEEDBACKS
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




