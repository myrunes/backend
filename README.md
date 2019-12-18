<div align="center">
    <img src="https://raw.githubusercontent.com/myrunes/myrunes/master/assets/logo-dark-1000-237.png" width="400"/>
    <br/>
    <strong>Save your League of Legends rune pages without wasting money.</strong><br><br>
    <img src="https://forthebadge.com/images/badges/made-with-go.svg" height="30" />&nbsp;
    <img src="https://forthebadge.com/images/badges/made-with-vue.svg" height="30" />&nbsp;
    <a href="https://stackshare.io/myrunes/myrunes"><img src="https://img.shields.io/badge/tech-stack-blue?style=for-the-badge" height="30"/></a>&nbsp;
    <a href="https://zekro.de/discord"><img src="https://img.shields.io/discord/307084334198816769.svg?logo=discord&style=for-the-badge" height="30"></a>
    <br/><br/>
    <a href="https://hub.docker.com/r/zekro/myrunes"><img src="https://img.shields.io/docker/cloud/automated/zekro/myrunes.svg?color=cyan&logo=docker&logoColor=cyan&style=for-the-badge" height="30"></a>&nbsp;
    <a href="https://travis-ci.org/myrunes/myrunes"><img src="https://img.shields.io/travis/myrunes/myrunes.svg?logo=travis&style=for-the-badge" height="30"></a>&nbsp;
    <a href="https://github.com/myrunes/myrunes/actions"><img src="https://img.shields.io/github/workflow/status/myrunes/myrunes/CI?label=Actions&logo=github&style=for-the-badge" height="30"/></a>&nbsp;
</div>

---

# Introduction

MYRUNES is a little web tool where you can simply create and store League of Legends rune pages without spending ingame (or even real) money for rune pages. Just visit [myrunes.com](https://myrunes.com), create an account and save your runes to be ready for the next pick and ban.  
Of course, if you don't trust us, you can download the source code and build the binaries and front end to be hosted on your own server environment.

---

# To Do & Future Goals


To see a list of current goals, ideas and bugs, please take a look to the MYRUNES [**Trello Board**](http://todo.myrunes.com).

---

# Self Hosting

## Docker

You can self-host this application by using the supplied [**docker images**](https://cloud.docker.com/u/zekro/repository/docker/zekro/myrunes).

Just use the following command to pull the latest stable image:  
```
# docker pull zekro/myrunes:latest
```

On startup, you need to bind the exposed web server port `8080` and the volume `/etc/myrunes` to your host system:

```
# docker run \
  -p 443:8080 \
  -v /etc/myrunes:/etc/myrunes \
  zekro/myrunes:latest
```

You can use following configuration with a MongoDB container using Docker Compose:

```yml
version: '3'

services:

  mongo:
    image: 'mongo:latest'
    expose:
      - '27017'
    volumes:
      - './mongodb/data/db:/data/db'
      - '/home/mgr/dump:/var/dump'
    command: '--auth'
    restart: always
 
  myrunes:
    image: "zekro/myrunes:latest"
    ports:
      - "443:8080"
    volumes:
      - "/etc/myrunes:/etc/myrunes"
    environment:
      # You dont need to define the configuration
      # with environment variables fi you prefer 
      # using the config file instead.
      - 'DB_HOST=mongo'
      - 'DB_PORT=27017'
      - 'DB_USERNAME=myrunes'
      - 'DB_PASSWORD=somepw'
      - 'DB_AUTHDB=myrunes'
      - 'DB_DATADB=myrunes'
      - 'TLS_ENABLE=true'
      - 'TLS_KEY=/etc/cert/key.pem'
      - 'TLS_CERT=/etc/cert/cert.pem'
    ports:
      - '443:8080'
    restart: always
```

## As daemon

First of all, if you want to self host the MYRUNES system, your environment should pass certain requirements:

- [**MongoDB**](https://www.mongodb.com/)  
  The server application uses MongoDB as database and storage system.

- **[PM2](https://pm2.io/)** or **[screen](https://linux.die.net/man/1/screen)**  
  ...or something else to deamonize an application which is highly recommended for running the server component.

Also, you need the following toolchains for building the backend and frontend components:

- **[git](https://git-scm.com/)**
- **[go compiler toolchain](https://golang.org/)**
- **[dep package manager](https://github.com/golang/dep)**
- **[nodejs](https://nodejs.org/en/)** and **[npm](https://www.npmjs.com/)** *(npm will be automatically installed with nodejs)*
- **[Vue CLI](https://cli.vuejs.org/)**

Also, it is highly recommended to install **[GNU make](https://www.gnu.org/software/make/)** to simplify the build process. If you are using windows, you can install **[make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)**.

## Compiling

1. Set up GOPATH, if not done yet. Read [here](https://golang.org/pkg/go/build/#hdr-Go_Path) how to do this.

2. Clone the repository into your GOPATH:  
   ```
   $ git clone https://github.com/myrunes/myrunes $GOPATH/src/github.com/myrunes/myrunes
   $ cd $GOPATH/src/github.com/myrunes/myrunes
   ```

3. Build binaries and assets using the `Makefile`:  
   ```
   $ make
   ```

Now, the server binary and the web assets are located in the `./bin` directory. You can move them wherever you want, just always keep the `web` folder in the same location where the server binary is located to ensure that all web assets can be found by the web server.

## Startup

Now, you just need to start the server binary passing the location of your preferred config location. A preset config file will be then automatically created. Now, enter your preferences and restart the server.

```
$ ./server -c /etc/myrunes/config.yml
```

--- 

