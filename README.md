<div align="center">
    <!-- <img src="" width="200"/> -->
    <h1>~ MYRUNES ~</h1>
    <strong>Save your League of Legends rune pages without wasting money.</strong><br><br>
    <img src="https://forthebadge.com/images/badges/made-with-go.svg" height="30" />&nbsp;
    <img src="https://forthebadge.com/images/badges/made-with-vue.svg" height="30" />&nbsp;
    <img src="https://forthebadge.com/images/badges/fuck-it-ship-it.svg" height="30" />&nbsp;
    <a href="https://zekro.de/discord"><img src="https://img.shields.io/discord/307084334198816769.svg?logo=discord&style=for-the-badge" height="30"></a>
</div>

---

# Introduction

MYRUNES is a little web tool where you can simply create and store League of Legends rune pages without spending ingame (or even real) money for rune pages. Just visit [myrunes.com](https://myrunes.com), create an account and save your runes to be ready for the nex pick and ban.  
Of course, if you don't trust us, you can download the source code and build the binaries and front end to be hosted on your own server envoirement.

---

# To Do & Future Goals

Here you can see a list of stuff which needs to be done and ideas we have to be implemented later *(order by potential priority)*.

- [ ] favorite champions on the front page
- [ ] export and import rune pages
- [ ] sort champions by lane / position
- [ ] share pages by public link
- [ ] 'forgot password' option by using a mail server
- [ ] public profiles and pages list

---

# Self Hosting

First of all, if you want to self host the MYRUNES system, your envoirement should pass certain requirements:

- [**MongoDB**](https://www.mongodb.com/)  
  The server application uses MongoDB as database and storage system.

- **[PM2](https://pm2.io/)** or **[screen](https://linux.die.net/man/1/screen)**  
  ...or something else to deamonize an application which is highly recommendet for running the server component.

Also, you need the following toolchains for building the backend and frontend components:

- **[go compiler toolchain](https://golang.org/)**
- **[dep package manager](https://github.com/golang/dep)**
- **[nodejs](https://nodejs.org/en/)** and **[npm](https://www.npmjs.com/)** *(npm will be automatically installed with nodejs)*
- **[Vue CLI](https://cli.vuejs.org/)**

Also, it is recommendet to install **[GNU make](https://www.gnu.org/software/make/)** to simplyfy the build process. If you are using windows, you can install **[make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)**.

