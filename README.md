# About

This CTF bot was created to emulate a victim client and test client-side vulnerabilities. 

There is a yaml configuration file with which you can configure cookie settings for your bot, depending on which you can prepare the bot for your task. (e.g. xss or cors)


# Requirements:
* Go 1.21+

        $ go install github.com/cotsom/ctf-client-bot@latest

# Configure
    domain: localhost
    cookie:
      session: example.bot.token
    httpOnly: true
    timeout: 200ms

domain -> domain in cookie

cookie -> cookies list (key: value)

httponly -> httponly in cookie

timeout -> the time for which the bot will stay on the page (to wait for its full loading)

# Usage

### Basic
    CTF client bot

    Usage:
        bot -f [config.yml]

### Docker
    docker compose up -d

After launching, the bot will work on port `5555`. 

To get him to visit your site, you just need to send him an http request: 

    $ curl <botip>:5555?url=https://google.com
