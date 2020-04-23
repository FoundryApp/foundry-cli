
# Foundry - The fastest way to build Firebase CloudFunctions

- Website: [https://foundryapp.co](https://foundryapp.co)
- Docs: [https://docs.foundryapp.co](https://docs.foundryapp.co)
- Community Slack: [Join Foundry Community Slack](https://join.slack.com/t/community-foundry/shared_invite/zt-dcpyblnb-JSSWviMFbRvjGnikMAWJeA)
- Youtube channel: [Foundry Youtube Channel](https://www.youtube.com/channel/UCvNVqSIXlW6nSPlAvW78TQg)

<img alt="Foundry" src="https://firebasestorage.googleapis.com/v0/b/foundryapp.appspot.com/o/foundry-logo.svg?alt=media&token=9625306d-3577-4aab-ab12-bbde0daae849" width="600px">

Foundry lets you build your Firebase Cloud Functions notably faster, with less configuration, and with easy access to your production data.
Foundry consists of an open-sourced command-line tool called Foundry CLI and a pre-configured cloud environment for your development.


## Watch the 5-minute video explaining Foundry
<!-- [![Watch the 5 min video explaining Foundry](https://img.youtube.com/vi/wYPbR8MnNfE/maxresdefault.jpg)](https://youtu.be/wYPbR8MnNfE) -->
[![Watch the 5-min video explaining Foundry](https://firebasestorage.googleapis.com/v0/b/foundryapp.appspot.com/o/video-thumbnail.png?alt=media&token=a0273107-e55c-42a6-b6d2-bb24a1da722c)](https://youtu.be/wYPbR8MnNfE)


The key features of Foundry are:
- **Develop with a copy of your production data:** Specify what data you want to copy from your production Firestore and production users. We copy the data and fill the emulated Firestore and Firebase Auth. No need to maintain any custom scripts. You access this data as you would normally in your Firebase functions code - with the official Admin SDK.

- **Real-time feedback:** You don't have to manually trigger your functions to run them, Foundry triggers them for you every time you make a change in your code and sends you back the output usually within 1-2 seconds. You just define your Cloud Functions and how you want to trigger them in the configuration file. It's like Read-Eval-Print-Loop for your Cloud Functions.

- **Develop in the environment identical to the production environment:** Your Firebase Cloud Functions will run in a cloud environment that is identical to the environment where your functions are deployed. This way, you won't have unexpected production bugs. You don't have to create a separate Firebase project as your staging environment. Foundry is your staging environment.

- **Zero environment configuration:** There isn't any configuration. Just run `$ foundry init` and then `$ foundry go` and you're ready.

- **Easily test integration of your Cloud Functions:** Foundry gives you an access the emulated Firestore database and emulated users. You can specify with what data they should be filled with and what parts of production Firestore and users data should be copied to the cloud development environment. Together with the specification of how your Cloud Functions should be triggered every time you save your code, Foundry can load and trigger your Cloud Functions in the same way as they would be triggered on the Firebase platform.


## Getting Started & Documentation
Documentation is available on the [Foundry website](https://docs.foundryapp.co)

### Quick start

macOS:
```bash
$ brew tap foundryapp/foundry-cli
$ brew install foundry

$ cd <directory where is a package.json for your Firebase Functions>
$ foundry init
$ foundry go
```

Linux:
```bash
# Download the pre-combiled binary
$ curl https://github.com/FoundryApp/foundry-cli/releases/download/0.1.0/foundry-linux-0.1.0 --output ./foundry

# Add Foundry to one of the directories in your PATH
$ mv ./foundry /usr/local/bin/foundry

$ cd <directory where is a package.json for your Firebase Functions>
$ foundry init
$ foundry go
```

## License
[Mozilla Public License v2.0](https://github.com/foundryapp/foundry-cli/blob/master/LICENSE)
