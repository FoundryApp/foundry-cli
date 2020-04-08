# Foundry - The fastest way to build Firebase Functions

<!-- [Foundry](https://foundryapp.co) is a CLI tool that creates the same environment as is the environment where your cloud functions run in the production. The environment is pre-configured, with the copy of your production data, and automatically runs your code with the near-instant feedback.

Foundry let's you to develop your Firebase Functions in the same environment as is your production environemnt with access to the production data from your Firestore database. Your code is evaluated in the cloud environmentEverything happens in your p -->

<!-- Foundry watches your functions' code and creates a copy of your Firebase Function production environment for your development.  -->

<img alt="Foundry" src="https://firebasestorage.googleapis.com/v0/b/foundryapp.appspot.com/o/foundry-logo.svg?alt=media&token=9625306d-3577-4aab-ab12-bbde0daae849" width="600px">


Foundry is a CLI tool for building Firebase Functions. Foundry connects you to a cloud environment that is identical to the production environment of your Firebase Functions where everything works out-of-the-box. Together with the [config file](#Config), the cloud environment gives you an access to a copy of your production Firestore database and Firebase Auth users.<br/>
With Foundry CLI, you can feel sure that your code behaves correctly and same as in the production already during the development.


**TODO: GIF HERE**

The key features of Foundry are:
- **Out of the box environment**: Foundry connects you to a pre-configured cloud environment where you can interactively develop your Firebase Functions. No need to configure anything.

- **REPL for you Firebase Functions**: Foundry watches your functions' code for changes. With every change, it sends your code to the cloud environment, evaluates the code there, triggers your functions and sends back the results. Everything is automated, you can just keep coding.

- **Short deploy times and instant feedback**: Your code is always deployed by default. Every code change to your Firebase Functions triggers the CLI that pushes your code to the cloud environment. The output is sent back to you usually within 2 seconds. There isn't any waiting for your code to get deployed, it's always deployed.

- **Access to the production data**: The [config file](#Config) makes it easy to specify what part of your production Firestore and Auth users should be copied to the emulated Firestore in the cloud environment. You access this data the same way as you would in the production - with the official [Firebase Admin SDK](https://firebase.google.com/docs/admin/setup)

- **Continuous feedback**: Pre-define with what data should each Firebase Function be triggered in the [config file](#Config). The functions are then automatically triggered with every code change. This ensures that you always know whether your functions behave correctly against your production data.

- **Discover production bugs**: TODO

## Table of contents
- **[Installation](#Installation)**
- **[Usage](#Usage)**
- **[Config](#Config)**
  - **[Firestore](#Firestore)**
  - **[Authentication](#Authentication)**
  - **[Functions](#Functions)**
- **[FAQ](#FAQ)**

## Download

Download the latest version of Foundry

- **[macOS](https://github.com/FoundryApp/foundry-cli/releases)**

- **[Linux](TODO)**

Add the downloaded binary to one of folders in your system's `PATH` variable.

## Supported languages
Javascript

## Config file `foundry.yaml`
For Foundry to work, it requires that its config file - `foundry.yaml` - is present. You can use `$ foundry init` To generate the initial config file.<br/>
Make sure to call this command from a folder where your `package.json` for your Firebase Functions is placed.

```yaml
# [OPTIONAL] An array of glob patterns for files that should be ignored. The path is relative to the file's dir.
# If the array is changed, the CLI must be restarted for it to take the effect
ignore:
  - node_modules # Skip the whole node_modules directory
  - .git # Skip the whole .git directory
  - "**/*.*[0-9]" # Skip all temp files ending with number
  - "**/.*" # Skip all hidden files
  - "**/*~" # Skip vim's temp files


# [OPTIONAL] Path to your
# serviceAcc: ""


# [OPTIONAL] Describe emulated Firebase Auth users
auth:


# [OPTIONAL] Describe emulated Firestore in your cloud environment
firestore:


# [REQUIRED] An array of Firebase functions that should be evaluated by Foundry. All described functions must be exported in your root index.js
functions:

```


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

### How do I get a service account JSON for my Firebase project?

Go to [Google Cloud Console](https://console.cloud.google.com/iam-admin/serviceaccounts) and choose your project

or from [Firebase Console](https://console.firebase.google.com/project)

### Are you storing my code or data from my production environment?

We don't store your code or data for duration that is longer than the lifetime of your session (specify until pod dies?).

## License
[Mozilla Public License v2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE)
