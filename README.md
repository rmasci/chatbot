# Chatbot Demo

This is a small program, written mainly with help from a Udemy class I am taking that shows how to use websockets. Fantastic course by Trevor Sawler:

[Working with WebSockets in Go (Golang) | Udemy](https://www.udemy.com/course/working-with-websockets-in-go/)

You should be able to clone this directory and run it. 

```bash
go run *.go -b chatbot -h 127.0.0.1 -p 8888 -r rive
```

> **Note:** Change **127.0.0.1** to your ip address if you want to test from other systems on your network.

When Started the console looks like this:

```bash
chatbot]$  go run *.go -b cloudgenie -h 10.1.1.130 -p 8888 -r rive
2022/07/12 17:43:39 Starting channel listener
2022/07/12 17:43:39 Starting Web Server On Port 8080
2022/07/12 17:43:39 rivescript initialized
```

This needs a lot of work, especially in the Javascript area. I don't speak Javascript.

When running you can ask the 'chatbot' questions. You'll need my tcpscan https://github.com/rmasci/tcpscan if you want to run the rivescript that executes a LocalCommand.

![](img/chatbotDemo.png)

Note: tcpscan can be downloaded from:

[GitHub - rmasci/tcpscan: Verify TCP connectivity.](https://github.com/rmasci/tcpscan)

## LocalCommand

In this version, if the first word of the rivescript response is LocalCommand that tells the chatbot it needs to exec the rest of the statement.   For example if you wanted to show the current date:

```rivescript
+ what is the current date
- LocalCommand date
```

When typed to the chatbot it looks like this:

![](img/chat.png)

For the tcpscan example above, the rivescript looks like this:

```rivescript
+ tcpscan *
- LocalCommand /usr/local/bin/tcpscan <star1>
```

When a user types: chatbot tcpscan https://www.google.com -s it uses a <star1>, so that any options to tcpscan can be passed along. You could also hard code those:

```rivescript
+ scan yahoo
- LocalCommand /usr/local/bin/tcpscan https://www.yahoo.com -s
```

Result:

```text
13:46 Rich: chatbot scan yahoo
13:46 chatbot: @Rich,
+------------------+---------+-----------+-------------+---------------------------+
|          Address |    Port |    Status |         TCP |                       SSL |
+==================+=========+===========+=============+===========================+
|    www.yahoo.com |     443 |      Open |    277.36ms |    TLS v1.2 / OK: 28 days |
+------------------+---------+-----------+-------------+---------------------------+
```
