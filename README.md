Freehold Client
===================
Freehold Client is a Go specific client API for interacting with a given [freehold](https://bitbucket.org/tshannon/freehold) instance.

Usage is as follows:
```
	client, err := freeholdclient.New("https://freeholdinstance.org", "username", "passwordortoken", nil)
	if err != nil {
		panic(err)
	}

```

If you are running your freehold instance with a self-signed cert, you may want to specify your own cert handling using the optional tls config.
```
	client, err := freeholdclient.New("https://freeholdinstance.org", "username", "passwordortoken", 
		&tls.Config{InsecureSkipVerify: true}, //Ideally you'd setup proper cert authority validation
	)
	if err != nil {
		panic(err)
	}

```


It is recommended that when using the client that you generate a security token, and access the freehold instance, rather than storing the end user's password locally on the machine.

```
	client, err := freeholdclient.New(rootURL, username, password, tlsCfg)
	if err != nil {
		return err
	}

	token, err := client.NewToken("Token Name", "", "", time.Now().AddDate(0, 6, 0))
	if err != nil {
		return err
	}

	client, err = freeholdclient.New(rootURL, username, token.Token, tlsCfg)
	if err != nil {
		return err
	}

	//store token.Token for use later, and forget password

```
