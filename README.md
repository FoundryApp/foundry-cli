# Foundry - The fastest way to build Firebase Functions

<!-- [Foundry](https://foundryapp.co) is a CLI tool that creates the same environment as is the environment where your cloud functions run in the production. The environment is pre-configured, with the copy of your production data, and automatically runs your code with the near-instant feedback.

Foundry let's you to develop your Firebase Functions in the same environment as is your production environemnt with access to the production data from your Firestore database. Your code is evaluated in the cloud environmentEverything happens in your p -->

<!-- Foundry watches your functions' code and creates a copy of your Firebase Function production environment for your development.  -->

<img alt="Foundry" src="https://firebasestorage.googleapis.com/v0/b/foundryapp.appspot.com/o/foundry-logo.svg?alt=media&token=9625306d-3577-4aab-ab12-bbde0daae849" width="600px">


Foundry is a CLI tool for building Firebase Functions. Foundry connects you to a cloud environment that is identical to the production environment of your Firebase Functions where everything works out-of-the-box. Together with the [config file](#Config), the cloud environment gives you an access to a copy of your production Firestore database and Firebase Auth users.<br/>
With Foundry CLI, you can feel sure that your code behaves correctly and same as in the production already during the development.


**TODO: GIF HERE**

Mention that logs are sent back right into your terminal, while you code.

The key features of Foundry are:
- **Out of the box environment**: Foundry connects you to a pre-configured cloud environment where you can interactively develop your Firebase Functions. No need to configure anything.

- **REPL for you Firebase Functions**: Foundry watches your functions' code for changes. With every change, it sends your code to the cloud environment, evaluates the code there, triggers your functions and sends back the results. Everything is automated, you can just keep coding.

- **Short deploy times and instant feedback**: Your code is always deployed by default. Every code change to your Firebase Functions triggers the CLI that pushes your code to the cloud environment. The output is sent back to you usually within 1-2 seconds. There isn't any waiting for your code to get deployed, it's always deployed.

- **Access to the production data**: The [config file](#config-file) makes it easy to specify what part of your production Firestore and Auth users should be **copied** to the emulated Firestore and Firebase Auth in the cloud environment. You access this data the same way as you would in the production - with the official [Firebase Admin SDK](https://firebase.google.com/docs/admin/setup). There's no need to create separate Firebase projects or maintain local scripts to test your functions so you don't corrupt your production data.

- **Continuous feedback**: Pre-define with what data should each Firebase Function be triggered in the [config file](#config-file). The functions are then automatically triggered with every code change. This ensures that you always know whether your functions behave correctly against your production data. There isn't any context switching and no need to leave your coding editor.

- **Discover production bugs**: TODO

## Table of contents
- **[How Foundry works](#how-foundry-works)**
- **[Download](#download)**
- **[Supported languages](#supported-languages)**
- **[Config file](#config-file)**
  - **[Functions](#functions)**
  - **[Firestore](#firestore)**
  - **[Auth](#auth)**
  - **[Ignore files](#ignore-files)**
  - **[Service account](#authentication)**
- **[Getting started](#getting-started)**
- **[Interactive prompt](#interactive-prompt)**
- **[Supported Firebase features](#supported-firebase-features)**
- **[Examples](#examples)**
- **[FAQ](#faq)**
- **[Slack community](#slack-community)**

## How Foundry works
TODO

## Download

- **[macOS](https://github.com/FoundryApp/foundry-cli/releases)**

- **[Linux](TODO)**

Add the downloaded binary to one of the folders in your system's `PATH` variable.

## Supported languages
Javascript

## Config file
For Foundry to work, it requires that its config file - `foundry.yaml` - is present. You can run `$ foundry init` to generate a basic config file.<br/>
Make sure to call this command from a folder where your `package.json` for your Firebase Functions is placed - `foundry.yaml` must always be placed next to the Firebase Function's `package.json` file. 
<br/>

Here's a full example of the config file:
```yaml
# [OPTIONAL]
# An array of glob patterns for files that should be ignored. The path is relative 
# to the config file's path.
# If the array is changed, the CLI must be restarted for it to take the effect.
ignore:
    # Skip all node_modules directories
  - "**/node_modules"
    # Skip the whole .git directory
  - .git 
    # Skip all temp files ending with number
  - "**/*.*[0-9]" 
    # Skip all hidden files
  - "**/.*"
    # Skip Vim's temp files
  - "**/*~"


# [OPTIONAL] 
# A path to a service account for your Firebase project. 
# See <TODO:URL> for more info on how to obtain your service account.
serviceAcc: path/to/service/account.json



# [OPTIONAL] 
# An array describing emulated Firebase Auth users in your cloud environment
users:
  # You can describe your emulated Auth users either directly
  - id: user-id-1
    # The 'data' field takes a JSON string
    data: '{"email": "user-id-1-email@email.com"}'
  # Or you can copy your production users from Firebae Auth by using 'geFromProd'
  # (WARNING: service account is required!):
  # If the value is a number, Foundry takes first N users from Firebase Auth
  - getFromProd: 2
  # If the value is an array, Foundry expects that the array's elements
  # are real IDs of your Firebase Auth users
  - getFromProd: [id-of-a-user-in-production, another-id]

  # You can use both the direct and 'geFromProd' approach simultaneously
  # The final Firebase Auth users will be a merge of these

# [OPTIONAL]
# An array describing emulated Firestore in your cloud environment
firestore:
  # TODO:
  # - collection: workspaces
  #   docs:
  #     - id: ws-id-1
  #       data: '{"userId": "user-id-1"}'
  #     - id: ws-id-2
  #       data: '{"userId": "user-id-2"}'
  #     - getFromProd: 2
  #     - getFromProd: [wp-id-in-prod]
  
  # You can describe your emulated Firestore either directly
  - collection: workspaces
    docs:
      - id: ws-id-1
        data: '{"userId": "user-id-1"}'
      - id: ws-id-2
        data: '{"userId": "user-id-2"}'
      
  # Or you can copy data from your production Firestore by using 'getFromProd'
  # (WARNING: service account is required!):
  - collection: workspaces
    # If the value is a number, Foundry takes first N documents from the 
    # specified collection (here 'workspaces')
    getFromProd: 2
    # If the value is an array, Foundry expects that the array's elements
    # are real IDs of documents in the specified collection (here documents
    # in the 'workspaces' collection)
    getFromProd: [workspace-id-1, workspace-id-2]
    
  # To create a collection or a document that is inside another collection:
  - collection: collection/doc-id/subcollection
    docs:
      - id: doc-in-subcollection
        data: '{}'


# [REQUIRED] 
# An array describing your Firebase functions that should be evaluated by Foundry. 
# All described functions must be exported in the function's root index.js file.
# In this array, you describe how Foundry should trigger each function in every run.
functions:
  # Foundry currently supports following types of Firebase Functions:
  # - https
  # - httpsCallable
  # - auth
  # - firestory
  
  # Each function has at least 2 fields:
  # 'name' - the same name under which your function is exported from your function's root index.js file
  # 'type' - one of the following: https, httpsCallable, auth, firestore
  
  
  # -----------------------
  # Type: https
  # A 'https' functions is the equivalent of
  # https://firebase.google.com/docs/functions/http-events
  - name: myHttpsFunction
    type: https
    # The payload field can either be a JSON string
    payload: '{"field":"value"}'
    # or you can reference a document from your production Firestore
    payload:
      doc:
        getFromProd:
          collection: path/to/collection
          id: doc-id
    # or you can reference a document from the emulated Firestore
    payload:
      doc:
        collection: path/to/collection
        id: doc-id
  
  # -----------------------
  # Type: httpsCallable
  # A 'httpsCallable' is the equivalent of  
  # https://firebase.google.com/docs/functions/callable
  - name: myHttpsCallableFunction
    type: httpsCallable
    # Since 'httpsCallable' function is meant to be called
    # from inside your app by your users, it expects a 
    # field 'asUser'
    # With this field you specify as what user should this
    # function be triggered
    asUser:
      # A user with this ID must be present in the emulated Firebase Auth users
      id: user-id
    # The 'payload' field is the same as in the 'https' function
    payload: '{}'
      
    
  # -----------------------
  # Type: auth
  # An 'auth' function is the equivalent of
  # https://firebase.google.com/docs/functions/auth-events
  # Based on the 'trigger' field, there are 2 sub-types of 
  # an 'auth' function: onCreate, onDelete
  - name: myAuthOnDreateFunction
    type: auth
    trigger: onCreate
    # The 'createUser' field specifies a new user record
    # that will trigger this auth function.
    # Keep in mind that this user will actually be created
    # in the emulated Firebase Auth users!
    createUser:    
      id: new-user-id
      data: '{"email": "new-user@email.com"}'
    # You can also reference a Firebase auth user from your
    # production. This user will get copied to the emulated
    # Firebase auth users and triggers this auth function:
    createUser:
      getFromProd:
        id: user-id-in-production
    
  - name: myAuthOnDeleteFunction
    type: auth
    trigger: onDelete    
    # This auth function will get triggered by deleting
    # a user with the specified ID from your emulated
    # Firebase Auth users.
    # Keep in mind that this user will actually be deleted
    # from the emulated Firebase Auth users!
    deleteUser:
      # A user with this ID must be present in the emulated Firebase Auth users
      id: existing-user-id  
  
  # -----------------------
  # Type: firestore
  # A 'firestore' function is the equivalent of
  # https://firebase.google.com/docs/functions/firestore-events
  # Based on the 'trigger' field, there are 3 sub-types of 
  # a 'firestore' function: onCreate, onDelete, onUpdate
```

### Field `functions`
It's important to understand how trigger functions work. Everything happens against the emulated Firestore database or the emulated Firebase Auth users. Both of these can be specified in the config file under fields `firestore` and `users` respectively.
### Field `firestore`
### Field `auth`
### Field `ignore`
### Field `serviceAcc`


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

## Supported Firebase features
TODO

## FAQ

### How do I get a service account JSON for my Firebase project?

Go to [Google Cloud Console](https://console.cloud.google.com/iam-admin/serviceaccounts) and choose your project

or from [Firebase Console](https://console.firebase.google.com/project)

### Why do you need a service account to my Firebase project?
TODO

### Are you storing my code or data from my production environment?

We don't store your code for a duration longer than the lifetime of your session. Once your session ends, your cloud environment is killed and only metadata (environment variables) is preserved.

## License
[Mozilla Public License v2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE)
