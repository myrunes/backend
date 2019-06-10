<div align="center">
    <img src="assets/logo-1000-237.png" width="400"/>
    <br/>
    <strong>Save your League of Legends rune pages without wasting money.</strong><br><br>
    <img src="https://forthebadge.com/images/badges/made-with-go.svg" height="30" />&nbsp;
    <img src="https://forthebadge.com/images/badges/made-with-vue.svg" height="30" />&nbsp;
    <img src="https://forthebadge.com/images/badges/fuck-it-ship-it.svg" height="30" />&nbsp;
    <a href="https://zekro.de/discord"><img src="https://img.shields.io/discord/307084334198816769.svg?logo=discord&style=for-the-badge" height="30"></a>
    <br/><br/>
    <img alt="Docker Automated build" src="https://img.shields.io/docker/automated/zekro/myrunes.svg?color=cyan&logo=docker&style=for-the-badge">&nbsp;
    <img alt="Travis (.org)" src="https://img.shields.io/travis/zekrotja/myrunes.svg?logo=travis&style=for-the-badge">
</div>

---

# Introduction

MYRUNES is a little web tool where you can simply create and store League of Legends rune pages without spending ingame (or even real) money for rune pages. Just visit [myrunes.com](https://myrunes.com), create an account and save your runes to be ready for the nex pick and ban.  
Of course, if you don't trust us, you can download the source code and build the binaries and front end to be hosted on your own server envoirement.

---

# To Do & Future Goals


To see a list of current goals, ideas and bugs, please take a look to the MYRUNES [**Trello Board**](http://todo.myrunes.com).

---

# Self Hosting

## Requirements

First of all, if you want to self host the MYRUNES system, your envoirement should pass certain requirements:

- [**MongoDB**](https://www.mongodb.com/)  
  The server application uses MongoDB as database and storage system.

- **[PM2](https://pm2.io/)** or **[screen](https://linux.die.net/man/1/screen)**  
  ...or something else to deamonize an application which is highly recommendet for running the server component.

Also, you need the following toolchains for building the backend and frontend components:

- **[git](https://git-scm.com/)**
- **[go compiler toolchain](https://golang.org/)**
- **[dep package manager](https://github.com/golang/dep)**
- **[nodejs](https://nodejs.org/en/)** and **[npm](https://www.npmjs.com/)** *(npm will be automatically installed with nodejs)*
- **[Vue CLI](https://cli.vuejs.org/)**

Also, it is highly recommendet to install **[GNU make](https://www.gnu.org/software/make/)** to simplyfy the build process. If you are using windows, you can install **[make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)**.

## Compiling

1. Set up GOPATH, if not done yet. Read [here](https://golang.org/pkg/go/build/#hdr-Go_Path) how to do this.

2. Clone the repository into your GOPATH:  
   ```
   $ git clone https://github.com/zekroTJA/myrunes $GOPATH/src/github.com/zekroTJA/myrunes
   $ cd $GOPATH/src/github.com/zekroTJA/myrunes
   ```

3. Build binaries and assets using the `Makefile`:  
   ```
   $ make
   ```

Now, the server binary and the web assets are located in the `./bin` directory. You can move them whereever you want, just always keep the `web` folder in the same location where the server binary is located to ensure that all web assets can be found by the web server.

## Startup

Now, you just need to start the server binary passing the location of your prefered config location. A preset config file will be then automatically created. Now, enter your preferences and restart the server.

```
$ ./server -c /etc/myrunes/config.yml
```

--- 

