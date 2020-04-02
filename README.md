# Foundry - The fastest way to develop serverless functions

## Overview
CLI that creates the same environment as is the environment where your cloud functions run in the production. The environment is pre-configured, with the copy of your production data, and automatically runs your code with the near-instant feedback.

Foundry [website](https://www.foundryapp.co/)

### Table of content
- **[Installation](#Installation)**
- **[Usage](#Usage)**
- **[Config](#Config)**
  - **[Firestore](#Firestore)**
  - **[Authentication](#Authentication)**
  - **[Functions](#Functions)**
- **[FAQ](#FAQ)**

## Installation

### Homebrew

    brew

### Standalone Binary

The you can download it here

- **[OSX](github release link)**

- **[Windows](github release link)**

(Add to PATH?)

### Go Package (Compile from source)

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
