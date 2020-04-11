# Foundry - The fastest way to build Firebase Functions

<!-- [Foundry](https://foundryapp.co) is a CLI tool that creates the same environment as is the environment where your cloud functions run in the production. The environment is pre-configured, with the copy of your production data, and automatically runs your code with the near-instant feedback.

Foundry let's you to develop your Firebase Functions in the same environment as is your production environemnt with access to the production data from your Firestore database. Your code is evaluated in the cloud environmentEverything happens in your p -->

<!-- Foundry watches your functions' code and creates a copy of your Firebase Function production environment for your development.  -->

<img alt="Foundry" src="https://firebasestorage.googleapis.com/v0/b/foundryapp.appspot.com/o/foundry-logo.svg?alt=media&token=9625306d-3577-4aab-ab12-bbde0daae849" width="600px">

Foundry is a command-line tool for building Firebase Functions. Foundry connects you to a cloud environment where everything works out-of-the-box and that is as much as possible identical to the production environment of your Firebase Functions. Together with the [config file](#config-file-foundryyaml), the cloud environment gives you access to a copy of your production Firestore database and Firebase Auth users.<br/>
With Foundry, you can feel sure that your code behaves correctly and the same as in the production already during the development.


<!-- **TODO: GIF/VIDEO HERE** -->

### Watch the 5 minutes video explaining Foundry
<!-- [![Watch the 5 min video explaining Foundry](https://img.youtube.com/vi/wYPbR8MnNfE/maxresdefault.jpg)](https://youtu.be/wYPbR8MnNfE) -->
[![Watch the 5 min video explaining Foundry](https://firebasestorage.googleapis.com/v0/b/foundryapp.appspot.com/o/video-thumbnail.png?alt=media&token=a0273107-e55c-42a6-b6d2-bb24a1da722c)](https://youtu.be/wYPbR8MnNfE)


The key features of Foundry are:
- **Out of the box environment**: Foundry connects you to a pre-configured cloud environment where you can interactively develop your Firebase Functions. You don't have to configure anything to debug & test your functions.

- **REPL for you Firebase Functions**: Foundry watches your functions' code for changes. With every change, it sends your code to the cloud environment, evaluates the code there, triggers your functions and sends back the results to your terminal. Everything is automated, no need for context switching, you can just keep coding.

- **Short upload times and quick response**: Your code is always deployed in the development cloud environment by default. Every code change notifies Foundry to push your code to your environment. The output is sent back to you usually within 1-2 seconds. There isn't any waiting for your code to get deployed, it's always deployed and you always know whether it's working.

- **Access to production data during development**: The [config file](#config-file-foundryyaml) makes it easy to specify what part of your production Firestore database and Firebase Auth users should be **copied** to the emulated Firestore and Firebase Auth in your development cloud environment. You access this data the same way as you would in the production - with the official [Firebase Admin SDK](https://firebase.google.com/docs/admin/setup). There's no need to create separate Firebase projects or maintain local scripts to test your functions so you don't corrupt your production data.

- **Continuous feedback**: Pre-define with what data should each Firebase Function be triggered in the [config file](#config-file-foundryyaml). The functions are then automatically triggered with every code change. This ensures that you always know whether your functions behave correctly with your production data.

- **Discover production bugs**: A lot of bugs happen only after you deploy your Functions onto the production. The biggest reason for this is not having easy access to your production data. Foundry solves that by copying your production Firestore database and production Firebase Auth users based on your [`foundry.yaml` config file](#config-file-foundryyaml).


## TL;DR to start Foundry
1. `$ cd <directory where is a package.json for your Firebase Functions>`
2. `$ foundry init`
3. [Add your Firebase Functions into Foundry `foundry.yaml` config file](#field-functions)
3. `$ foundry go`
<br/>

Once your cloud environment is ready, you can start coding. Once Foundry installs all You will see that Foundry triggers all functions you mentioned in the config file each time you save your code.

- [Read more on how to access production Firestore data and users](#config-file-foundryyaml)
- [Read more on how to set environment variables](#environment-variables)
- [Read more on how you to filter Foundry's output](#filtering-functions)
- [Read more on currently supported Firebase features](#supported-firebase-features)


## Table of contents
- **[How Foundry works](#how-foundry-works)**
- **[What Foundry doesn't do](#what-foundry-doesnt-do)**
- **[Download and installation](#download-and-installation)**
- **[Examples](#examples)**
- **[Slack community](#slack-community)**
- **[Supported languages](#supported-languages)**
- **[Config file `foundry.yaml`](#config-file-foundryyaml)**
  - **[Field `functions`](#field-functions)**
  - **[Field `firestore`](#field-firestore)**
  - **[Field `users`](#field-users)**
  - **[Field `ignore`](#field-ignore)**
  - **[Field `serviceAcc`](#field-serviceAcc)**
  - **[Full config file example](#full-config-file-example)**
- **[How to use the Foundry CLI](#how-to-use-the-foundry-cli)**
  - **[Initialization](#initialization)**
  - **[Connecting to your cloud environment](#connecting-to-your-cloud-environment)**
  - **[Environment variables](#environment-variables)**
- **[Supported Firebase features](#supported-firebase-features)**
- **[FAQ](#faq)**


## How Foundry works
Foundry makes the development of Firebase Functions faster and with access to a [copy of your production data](#field-firestore). It's a command-line tool that connects you to your own cloud environment for the development of Firebase Functions. This environment is as much as possible identical to the actual environment where you Firebase Functions run after the deployment. Your environment is pre-configured and your functions work out-of-the-box.<br/>
Once connected to your development environment, Foundry starts watching your code for changes. Every change notifies the CLI and your code is uploaded to your cloud environment. Foundry then triggers all of your Firebase Functions based on the [rules](#field-functions) specified in your `foundry.yaml` [config file](#config-file-foundryyaml). <br/>
Both `stdout` and `stderr` of your functions are sent back to your terminal after each of such runs. The whole upload loop with the transmission of data back to you usually takes about 1-2 seconds. This loop creates a REPL-like tool for your functions and makes it easy to be sure that your functions' code behaves correctly with the production data.<br/>

The [config](#config-file-foundryyaml) `foundry.yaml` file is a critical part of Foundry. It describes 3 main things:
1. [What Firebase Functions should Foundry register and **how** it should trigger them in each run](#field-functions)
2. [How should the emulated Firestore database look like](#field-firestore)
3. [How should the emulated Firebase Auth users look like](#field-auth)
<br/>

Having an emulated Firestore database and Firebase Auth users in your development environment makes it easy to test the full logic of your functions.

## What Foundry doesn't do
Foundry doesn't deploy your Firebase Functions onto the production. To deploy functions onto the production use the official [Firebase tool](https://github.com/firebase/firebase-tools).

## Download and installation

- **[macOS](https://github.com/FoundryApp/foundry-cli/releases)**

- **[Linux](TODO)**

To install Foundry, add the downloaded binary to one of the folders in your system's `PATH` variable.

## Examples
Check out a separate repo with [example projects](https://github.com/FoundryApp/examples)

## Slack community
[Join Foundry community on Slack](https://join.slack.com/t/community-foundry/shared_invite/zt-dcpyblnb-JSSWviMFbRvjGnikMAWJeA)

## Supported languages
JavaScript

## Config file `foundry.yaml`
For Foundry to work, it requires that its config file - `foundry.yaml` - is present. The config file must always be placed next to the Firebase Function's `package.json` file.<br/>
You can run `$ foundry init` to generate a basic config file. Make sure to call this command from a folder where your `package.json` for your Firebase Functions is placed.<br/>

Check out the full config file example [here](#full-config-file-example).

### Field `functions`
An array describing your Firebase functions that should be evaluated by Foundry. All described functions must be exported in the function's root `index.js` file. In this array, you describe how Foundry should trigger each function in every run.<br/>

It's important to understand how trigger functions work in Foundry. Everything happens against the emulated Firestore database or the emulated Firebase Auth users. Both can be specified in the config file under fields [`firestore`](#field-firestore) and [`users`](#field-users) respectively. So all of your code where you manipulate with Firestore or Firebase Auth happens against the emulated Firestore database and emulated Firebase Auth.<br/>
The same is true for function triggers you describe in the Foundry config file. The triggers usually describe how should the emulated Firestore database or emulated Firebase Auth users be mutated. In return, these mutations will trigger your functions.<br/>

Currently, Foundry supports following Firebase functions:
#### HTTPS Functions
Equivalent of - [https://firebase.google.com/docs/functions/http-events](https://firebase.google.com/docs/functions/http-events)
```yaml
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
```

#### HTTPS Callable Functions
Equivalent of - [https://firebase.google.com/docs/functions/callable](https://firebase.google.com/docs/functions/callable)
```yaml
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
```

#### Auth Trigger Functions
Equivalent of - [https://firebase.google.com/docs/functions/auth-events](https://firebase.google.com/docs/functions/auth-events)<br/>

`onCreate` trigger
```yaml
- name: myAuthOnCreateFunction
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
  # production by using the 'getFromProd' field.
  # This user will be copied to the emulated Firebase auth users
  # and will trigger this auth function:
  createUser:
    getFromProd:
      id: user-id-in-production
```

`onDelete` trigger
```yaml
- name: myAuthOnDeleteFunction
  type: auth
  trigger: onDelete
  # This auth function will be triggered by deleting
  # a user with the specified ID from your emulated
  # Firebase Auth users.
  # Keep in mind that this user will actually be deleted
  # from the emulated Firebase Auth users!
  deleteUser:
    # A user with this ID must be present in the emulated Firebase Auth users
    id: existing-user-id
```

#### Firestore Trigger Functions
Equivalent of - [https://firebase.google.com/docs/functions/firestore-events](https://firebase.google.com/docs/functions/firestore-events)<br/>

`onCreate` trigger
```yaml
- name: myFirestoreOnCreateFunction
  type: firestore
  trigger: onCreate
  # The 'createDoc' field creates a new specified document
  # that will trigger this firestore function.
  # Keep in mind that this document will actually be
  # create in the emulated Firestore!
  createDoc:
    collection: path/to/collection
    id: new-doc-id
    data: '{}'
  # You can also reference a document from your production
  # Firestore database.
  # This document will be copied to the emulated Firestore
  # database and will trigger this function.
  createDoc:
    getFromProd:
      collection: path/to/collection
      id: existing-doc-id
```

`onDelete` trigger
```yaml
- name: myFirestoreOnDeleteFunction
  type: firestore
  trigger: onDelete
  # The 'deleteDoc' field deletes a specified document from
  # the emulated Firestore database. The deletion will
  # trigger this firestore function.
  # Keep in mind that this document will actually be
  # deleted from the emulated Firestore database. So it must
  # exist first!
  deleteDoc:
    # A document inside this collection must exist in the
    # emulated Firestore database
    collection: path/to/collection
    id: existing-doc-id
```

`onUpdate` trigger
```yaml
- name: myFirestoreOnUpdateFunction
  type: firestore
  trigger: onUpdate
  # The 'updateDoc' field updates a specified document
  # from the emulated Firestore database with a new
  # data. The update will trigger this firestore function.
  # Keep in mind that this document will actually be
  # updated in the emulated Firestore database. So it must
  # exist first!
  updateDoc:
    collection: path/to/collection
    id: existing-doc-id
    # A JSON string specifying new document's data
    data: '{}'
```

### Field `firestore`
The field `firestore` gives you an option to have a separate Firestore database from your production Firestore database. This separate Firestore is an emulated Firestore database that lives in your cloud environment for the duration of your session. <br/>

The `firestore` field expects an array of collections.
You have two options on how to fill an emulated Firestore database.

1. Specify documents directly with JSON strings<br/>
```yaml
firestore:
  - collection: workspaces
    docs:
      - id: ws-id-1
        data: '{"userId": "user-id-1"}'
      - id: ws-id-2
        data: '{"userId": "user-id-2"}'
```

2. Specify what documents should be copied from your production Firestore database<br/>
Note that this approach requires you to specify the [`serviceAcc`](#field-serviceacc) field.
```yaml
firestore:
  - collection: workspaces
    docs:
      - getFromProd: [workspace-id-1, workspace-id-2]
      # This option tells Foundry to take first 2 documents from collection 'workspaces'
      # in your production Firestore database
      - getFromProd: 2
```

You can combine both the direct approach and `getFromProd` approach.
<br/>

To create a nested collection specify a full collection's path
```yaml
firestore:
  - collection: my/nested/collection
    docs: ...
```

### Field `users`
The same way you can [emulate](#field-firestore) Firestore database for your development you can also emulate Firebase Auth users.<br/>

You have two options on how to fill the emulated uses:
1. Directly with JSON strings
```yaml
users:
    - id: user-id-1
      # The 'data' field takes a JSON string
      data: '{"email": "user-id-1-email@email.com"}'
```

2. Specify what users should be copied from your production Firebase Auth<br/>
Note that this approach requires you to specify the [`serviceAcc`](#field-serviceacc) field.
```yaml
users:
    - getFromProd: [user-id-1, user-id-2]
      # If the value is a number, Foundry takes first N users from Firebase Auth
    - getFromProd: 2
```

### Field `ignore`
Often there are files and folders that you don't want Foundry to upload or watch. To ignore them, you can use an array of glob patterns. For example:
```yaml
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
```

### Field `serviceAcc`
For Foundry to be able to copy some of your production data to your development cloud environment it must have access to your Firebase project. This is done through a [service account](https://firebase.google.com/support/guides/service-accounts).
The field `serviceAcc` expects a path to a service account JSON file for your Firebase project.<br/>

Of course, if you aren't copying any of your production data you don't need to specify `serviceAcc`.<br/>


[See FAQ](#how-do-i-get-a-service-account-json-for-my-firebase-project) to learn how to obtain a service account to your Firebase project.

### Full config file example
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
# See https://github.com/FoundryApp/foundry-cli#field-serviceAcc
# for more info on how to obtain your service account.
serviceAcc: path/to/service/account.json


# [OPTIONAL]
# An array describing emulated Firebase Auth users in your cloud environment
users:
  # You can describe your emulated Auth users either directly
  - id: user-id-1
    # The 'data' field takes a JSON string
    data: '{"email": "user-id-1-email@email.com"}'
  # Or you can copy your production users from Firebase Auth by using 'getFromProd'
  # (WARNING: service account is required!):
  # If the value is a number, Foundry takes first N users from Firebase Auth
  - getFromProd: 2
  # If the value is an array, Foundry expects that the array's elements
  # are real IDs of your Firebase Auth users
  - getFromProd: [id-of-a-user-in-production, another-id]

  # You can use both the direct and 'getFromProd' approach simultaneously
  # The final Firebase Auth users will be a merge of these

# [OPTIONAL]
# An array describing emulated Firestore in your cloud environment
firestore:
  # You can describe your emulated Firestore either directly
  - collection: workspaces
    docs:
      - id: ws-id-1
        data: '{"userId": "user-id-1"}'
      - id: ws-id-2
        data: '{"userId": "user-id-2"}'
      # Or you can copy data from your production Firestore by using 'getFromProd'
      # (WARNING: service account is required!):
      # If the value is a number, Foundry takes first N documents from the
      # specified collection (here 'workspaces')
      - getFromProd: 2
      # If the value is an array, Foundry expects that the array's elements
      # are real IDs of documents in the specified collection (here documents
      # in the 'workspaces' collection)
      - getFromProd: [workspace-id-1, workspace-id-2]

  # You can use both the direct and 'getFromProd' approach simultaneously
  # The final documents will be a merge of these

  # To create a nested collection:
  - collection: collection/doc-id/subcollection
    docs:
      - id: doc-in-subcollection
        data: '{}'

# [REQUIRED]
# An array describing your Firebase functions that should be evaluated by Foundry.
# All described functions must be exported in the function's root index.js file.
# In this array, you describe how Foundry should trigger each function in every run.
functions:
  # Foundry currently supports the following types of Firebase Functions:
  # - https
  # - httpsCallable
  # - auth
  # - firestore

  # Each function has at least 2 fields:
  # 'name' - the same name under which a function is exported from your function's root index.js file
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
  - name: myAuthOnCreateFunction
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
    # production by using the 'getFromProd' field.
    # This user will be copied to the emulated Firebase auth users
    # and will trigger this auth function:
    createUser:
      getFromProd:
        id: user-id-in-production

  - name: myAuthOnDeleteFunction
    type: auth
    trigger: onDelete
    # This auth function will be triggered by deleting
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
  - name: myFirestoreOnCreateFunction
    type: firestore
    trigger: onCreate
    # The 'createDoc' field creates a new specified document
    # that will trigger this firestore function.
    # Keep in mind that this document will actually be
    # create in the emulated Firestore!
    createDoc:
      collection: path/to/collection
      id: new-doc-id
      data: '{}'
    # You can also reference a document from your production
    # Firestore database.
    # This document will be copied to the emulated Firestore
    # database and will trigger this function.
    createDoc:
      getFromProd:
        collection: path/to/collection
        id: existing-doc-id

  - name: myFirestoreOnDeleteFunction
    type: firestore
    trigger: onDelete
    # The 'deleteDoc' field deletes a specified document from
    # the emulated Firestore database. The deletion will
    # trigger this firestore function.
    # Keep in mind that this document will actually be
    # deleted from the emulated Firestore database. So it must
    # exist first!
    deleteDoc:
      # A document inside this collection must exist in the
      # emulated Firestore database
      collection: path/to/collection
      id: existing-doc-id

  - name: myFirestoreOnUpdateFunction
    type: firestore
    trigger: onUpdate
    # The 'updateDoc' field updates a specified document
    # from the emulated Firestore database with a new
    # data. The update will trigger this firestore function.
    # Keep in mind that this document will actually be
    # updated in the emulated Firestore database. So it must
    # exist first!
    updateDoc:
      collection: path/to/collection
      id: existing-doc-id
      # A JSON string specifying new document's data
      data: '{}'
```


## How to use the Foundry CLI
All available commands can be printed by calling `$ foundry --help`.

- `$ foundry init`<br/>
Creates an initial `foundry.yaml` config file in the current directory

- `$ foundry go`<br/>
Start an interactive prompt and connects you to your cloud development environment

- `$ foundry env-set ENV_NAME=VAL`<br/>
Sets the environment variable(s) in your cloud development environment

- `$ foundry env-delete ENV_NAME`<br/>
Deletes the environment variable(s) in your cloud development environment

- `$ foundry env-print`<br/>
Prints environment variable(s) in your cloud development environment

- `$ foundry sign-up`<br/>
Create a new Foundry account in your terminal

- `$ foundry sign-in`<br/>
Sign in to your Foundry account

- `$ foundry sign-out`<br/>
Sign out from your Foundry account



### Initialization
For Foundry to work it needs to have a [`foundry.yaml` config file](#config-file-foundryyaml). This config file must placed in the same directory as is the `package.json` file for your Firebase Functions.<br/>

To generate a basic config file run `$ foundry init`.

### Connecting to your cloud environment
To connect to your cloud development environment and start your session run `$ foundry go`. If you aren't signed in, this will create an anonymous account for you that can be linked to your actual account later once you sign up.<br/>

Foundry starts an interactive prompt, connects you to your cloud environment and starts watching your code for changes. Each change notifies the CLI to send your code into the environment where your functions are triggered. You see the output and errors from your functions inside the interactive prompt.

#### Filtering functions
Often, you have many functions and orienting in the output is hard. To trigger only specific functions in each run you can execute the `watch` command inside the prompt. Its format is `watch function_name_1 function_name_2`.<br/>
To stop filtering functions execute the command `watch:all`.

### Environment variables
You can set, delete, and print all  environment variables in your cloud development environment with the following commands respectively:

`$ foundry env-set ENV_1=VAL_1 ENV_2=VAL_2`<br/>
Sets the specified environment variables.

`$ foundry env-delete ENV_1 ENV_2`<br/>
Deletes the specified environment variables.

`$ foundry env-print`<br/>
Prints all existing environment variables.

## Supported Firebase features
We support all Firebase functions triggered by Firestore changes, except for the 'onWrite' function, and all functions triggered by Firebase Auth changes. We support both the HTTPS and the HTTPS callable functions.

The other features we currently support are Firestore and parts of Firebase Auth. You can access emulated version of these services through [`firebase-admin` SDK](https://firebase.google.com/docs/admin/setup) as you would normally.

### Functions
We support following Firebase Functions' methods from the [`firebase-functions` SDK](https://firebase.google.com/docs/reference/functions):
  - Firestore
    - `firestore().document(<documentPath>).onCreate`
    - `firestore().document(<documentPath>).onUpdate`
    - `firestore().document(<documentPath>).onDelete`
- Firebase Auth
  - `auth().user().onCreate`
  - `auth().user().onDelete`
- Https
  - `https.onRequest`
  - `https.onCall` - without providing the `context.auth.token` and the `context.instanceIdToken` arguments

### Firestore
All functionality excluding the security rules is supported.

### Firebase Auth
Methods we support:

- `getUser`
- `getUserByEmail`
- `createUser`
- `deleteUser`
- `updateUser`
- `verifyIdToken` - accepting user's `uid` in place of an `idToken`

Note: The [`UserRecord` type](https://firebase.google.com/docs/reference/admin/node/admin.auth.UserRecord) isn't fully supported yet. Following properties aren't implemented:
- `customClaims`
- `metadata`
- `multiFactor`
- `passwordHash`
- `passwordSalt`
- `providerData`
- `tenantId`
- `tokensValidAfterTime`

If you need any of those properties to be supported, please open an issue. We will happily implement them.


## FAQ

### How long is my development cloud environment active after I end the session?
The cloud environment exists only for the time of your session. Once your session ends, the environment is terminated.

### Do you store my code or data?
We don't store your code or any data for a duration longer than is the duration of your session. Once your session ends, your cloud environment is terminated and only metadata (environment variables) is preserved.

### Why do you need a service account to my Firebase project?
The service account is needed for any action that requires copying data from your production Firestore database or production Firebase Auth to their emulated equivalents.<br/>
You can definitely use Foundry without specifying a path to your service account. Some features just won't be available though.
### How do I get a service account JSON for my Firebase project?

1. Go to [https://console.firebase.google.com/](https://console.firebase.google.com/) and select your project.<br/>
2. In your Firebase dashboard, select the settings icon. The icon is placed at the top of the left sidebar.<br/>
<img alt="Service-Account-Step-1" src="https://firebasestorage.googleapis.com/v0/b/foundryapp.appspot.com/o/service-acc-tutorial%2Fservice-acc-step-1.png?alt=media&token=3887c551-2aa6-4555-bd8e-7e0fe6e73c5e">
3. A popup menu will appear, select "Project settings".<br/>
<img alt="Service-Account-Step-2" src="https://firebasestorage.googleapis.com/v0/b/foundryapp.appspot.com/o/service-acc-tutorial%2Fservice-acc-step-2.png?alt=media&token=f6ae3f0c-069c-4f76-8712-203f4ba07e39">
4. This will take you to your projects settings with links to different sections on top. Select "Service accounts."<br/>
<img alt="Service-Account-Step-3" src="https://firebasestorage.googleapis.com/v0/b/foundryapp.appspot.com/o/service-acc-tutorial%2Fservice-acc-step-3.png?alt=media&token=7726fba0-ba90-4aa4-89d4-86727df595df">
5. A window with info about the Firebase Admin SDK will appear. There, click on "Generate new private key" at the bottom.<br/>
<img alt="Service-Account-Step-4" src="https://firebasestorage.googleapis.com/v0/b/foundryapp.appspot.com/o/service-acc-tutorial%2Fservice-acc-step-4.png?alt=media&token=caa3f666-59fc-4917-b736-1169e95b471e">
6. This will present a modal window informing you about security implications. Read the text and then click on the "Generate key" button at the bottom. This will download a JSON file that is your service account key.<br/>
<img alt="Service-Account-Step-5" src="https://firebasestorage.googleapis.com/v0/b/foundryapp.appspot.com/o/service-acc-tutorial%2Fservice-acc-step-5.png?alt=media&token=05e630f8-4dd8-4c43-baed-2a80b8096691">

7. Copy the path of your service account key file you just download to the field's `serviceAcc` value in the `foundry.yaml` config file.

## License
[Mozilla Public License v2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE)
