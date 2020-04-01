# Foundry - The fastest way to develop serverless functions

## Overview
detailed description

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

The you can download it here:

- **[OSX](github release link)**

- **[Windows](github release link)**

(Add to PATH?)

### Go Package

Make sure you have all the [requirements] installed, then clone the repo

    git clone (this repo)

Then build the binary with

[//]: #

    go build

You can then add the binary to your PATH with

    export 

(Add to path optional?)
(Some installing script?)

## Usage

### Run
In the project directory type

    foundry go

and just code. Your code is evaluated on each save and the results are streamed right back.

### Prompt Commands



how to use/what to do

### Registration

The basic version can be used without registration, but

## Config

You can describe the emulated environment in the `foundry.yaml` file in your project root directory.

### Firestore

Before every run you can fill the emulated Firestore with values by adding this code into the config

    firestore:
      - collection: ''
        docs:
          - id: ''
            data: '{"":""}'

or if you added a [service account key](#How%20to%20get%20a%20service%20account%20JSON%20for%20your%20Firebase%20project) to config with 

    serviceAcc: </path/to/serviceAcc.json>

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
      - name: <exportedName>
        type: <https/firestore/auth>
        trigger: <onCreate/onDelete/onUpdate>







## FAQ

### How to get a service account JSON for your Firebase project

Go to

### Are you storing my code or data from my production environment?

We don't store your code or data for duration that is longer than the lifetime of your session (specify until pod dies?).

