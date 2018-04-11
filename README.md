**v0.1**

# multiversal-clipboard
Android and golang apps to share clipboard between Mac and smartphone using firebase

![demo gif](https://github.com/Salvatore-Giordano/multiversal-clipboard/blob/master/demo.gif)

## Why?

Ios already has this feature with Mac os, but Android doesn't.
Using these apps we're able to have a shared clipboard between mac os and Android devices.

## Installation

- To install the Android app download the apk and install it normally. Then open the app and click on register and copy the userId.

- To install the app on Mac os just download the binary file and run it with `./multiversal-clipboard <userId>`

### Development

If you want to develop or to install with your own firebase project you need to:

- create a firebase project
- add the android app to firebase and download the google-services.json and add it in the project
- install and configure golang 
- install dep (https://github.com/golang/dep)
- cd in project/go and run `dep ensure`
- create a service account key on firebase and download the json
- fill the struct in main.gop according to the json
- set the database url in main.go

## How it works?

The system uses the firebase realtime database. 
On android there is a service running in background listening for changes on the userId node and on clipboard events.
In golang a cron monitors both the clipboard and the firebase database.

## Contributions

Contributions of any kind are more than welcome! 
Feel free to fork and improve multiversal-clipboard in any way you want, make a pull request, or open an issue.

## Acknowledgements

Portions of this project utilize the following copyrighted material, the use of which is hereby acknowledged:

- https://github.com/robfig/cron
- https://github.com/getlantern/systray

## Upcoming

- [ ] Better Android app design
- [ ] Deploy Android app on play store
- [ ] Ability to set custom userId/account
- [ ] Add some ux on Mac os using systray
- [ ] Port on windows/linux
- [ ] Share the clipboard with more than one device
- [ ] Encrypt clipboards
