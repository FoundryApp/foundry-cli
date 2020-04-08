# Foundry - The fastest way to build Firebase Functions

<!-- [Foundry](https://foundryapp.co) is a CLI tool that creates the same environment as is the environment where your cloud functions run in the production. The environment is pre-configured, with the copy of your production data, and automatically runs your code with the near-instant feedback.

Foundry let's you to develop your Firebase Functions in the same environment as is your production environemnt with access to the production data from your Firestore database. Your code is evaluated in the cloud environmentEverything happens in your p -->

<!-- Foundry watches your functions' code and creates a copy of your Firebase Function production environment for your development.  -->

Foundry is a tool for building Firebase Functions. Foundry connects you to a pre-configured cloud environment that is identical to the production environment of your Firebase Functions. Together with the [config file](#Config), the cloud environment gives you an access to a copy of your production Firestore database and Firebase Auth users.


The key features of Foundry are:
- **Out of the box environment**: Foundry connects to a pre-configured cloud environment where you can interactively develop your Firebase Functions. No need to configure anything.

- **REPL for you Firebase Functions**: Foundry watches your functions' code for changes. With every change, it sends your code to the cloud environment, evaluates the code there, triggers your functions and sends back the results. Everything is automated, you can just keep coding.

- **Short deploy times and instant feedback**: Your code is always deployed by default. Every code change to your Firebase Functions triggers the CLI that pushes your code to the cloud environment. The output is sent back to you usually within 2 seconds. There isn't any waiting for your code to get deployed, it's always deployed.

- **Access to the production data**: The [config file](#Config) makes it easy to specify what part of your production Firestore and Auth users should be copied to the emulated Firestore in the cloud environment. You access this data the same way as you would in the production - with the official [Firebase Admin SDK](https://firebase.google.com/docs/admin/setup)

- **Automatic triggers for Firebase Functions**: Pre-define with what data should each Firebase Function be triggered in the [config file](#Config). The functions are then automatically triggered with every code change. This ensures that you always know whether your functions behave correctly.

- **Discover production bugs**: TODO

[https://www.foundryapp.co/](https://www.foundryapp.co/)

### Table of content
- **[Installation](#Installation)**
- **[Usage](#Usage)**
- **[Config](#Config)**
  - **[Firestore](#Firestore)**
  - **[Authentication](#Authentication)**
  - **[Functions](#Functions)**
- **[FAQ](#FAQ)**

## Installation

### Standalone Binary

The you can download it here

- **[macOS](github release link)**

- **[Windows](github release link)**

(Add to PATH?)

### Go Package (compile from source)

Make sure you have all the [requirements](#Compilation%20requirements) installed, then clone the repo

    git clone (this repo)

then go to the cloned folder

    cd

and build the binary with

[//]: #

    go build

you can then add the binary to your PATH with

    export PATH=$PATH:

(Add to path optional?)
(Some installing script?)

## Usage

### Run
In the project directory type

    foundry go

and just code. Your code is evaluated on each save and the results are streamed right back.

### Prompt Commands

you can watch only some functions from the config by using the `watch` command in prompt

    > watch <functionToWatch> <anotherFunctionToWatch>

to reset this setting type

    > (reset)

### Registration

The basic version can be used without registration, but

## Config

You can describe the emulated environment in the `foundry.yaml` file in your project root directory.

### Root Directory

    rootDir: .

### Service Account

    serviceAcc: <relativePathToServiceAccountJSON>

### Firestore

Before every run you can fill the emulated Firestore with values by adding this code into the config

    firestore:
      - collection: ''
        docs:
          - id: ''
            data: '{"":""}'

or if you added a [service account key](#How%20to%20get%20a%20service%20account%20JSON%20for%20your%20Firebase%20project) to config with

    serviceAcc: <path/to/serviceAcc.json>

you can fill the emulated Firestore directly from your production environment

    firestore:
      - collection: ''
        getFromProd: 5

      - collection: ''
        getFromProd: ['id1', 'id2']

### Authentication

### Functions

If you want a function to automatically run in our emulated environment

    functions:
      - name: <exportedFunctionName>
        type: <https/firestore/auth>
        trigger: <onCreate/onDelete/onUpdate>

      - name: <exportedFunctionName>
        type: <https/firestore/auth>
        trigger: <onCreate/onDelete/onUpdate>



## FAQ

### How to get a service account JSON for your Firebase project

Go to [Google Cloud Console](https://console.cloud.google.com/iam-admin/serviceaccounts) and choose your project

or from [Firebase Console](https://console.firebase.google.com/project)

### Are you storing my code or data from my production environment?

We don't store your code or data for duration that is longer than the lifetime of your session (specify until pod dies?).

## Compilation requirements

## License
[Mozilla Public License v2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE)
