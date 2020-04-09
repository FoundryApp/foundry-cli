# Foundry - The fastest way to build Firebase Functions

<!-- [Foundry](https://foundryapp.co) is a CLI tool that creates the same environment as is the environment where your cloud functions run in the production. The environment is pre-configured, with the copy of your production data, and automatically runs your code with the near-instant feedback.

Foundry let's you to develop your Firebase Functions in the same environment as is your production environemnt with access to the production data from your Firestore database. Your code is evaluated in the cloud environmentEverything happens in your p -->

<!-- Foundry watches your functions' code and creates a copy of your Firebase Function production environment for your development.  -->

<img alt="Foundry" src="https://firebasestorage.googleapis.com/v0/b/foundryapp.appspot.com/o/foundry-logo.svg?alt=media&token=9625306d-3577-4aab-ab12-bbde0daae849" width="600px">


Foundry is a CLI tool for building Firebase Functions. Foundry connects you to a cloud environment that is identical to the production environment of your Firebase Functions where everything works out-of-the-box. Together with the [config file](#Config), the cloud environment gives you an access to a copy of your production Firestore database and Firebase Auth users.<br/>
With Foundry CLI, you can feel sure that your code behaves correctly and same as in the production already during the development.


**TODO: GIF/VIDEO HERE**

Mention that logs are sent back right into your terminal, while you code.

The key features of Foundry are:
- **Out of the box environment**: Foundry connects you to a pre-configured cloud environment where you can interactively develop your Firebase Functions. No need to configure anything.

- **REPL for you Firebase Functions**: Foundry watches your functions' code for changes. With every change, it sends your code to the cloud environment, evaluates the code there, triggers your functions and sends back the results. Everything is automated, you can just keep coding.

- **Short upload times and quick response**: Your code is always deployed in the development environment by default. Every code change to your Firebase Functions triggers the CLI that pushes your code to the cloud environment. The output is sent back to you usually within 1-2 seconds. There isn't any waiting for your code to get deployed, it's always deployed.

- **Access to production data during development**: The [config file](#config-file) makes it easy to specify what part of your production Firestore and Auth users should be **copied** to the emulated Firestore and Firebase Auth in the cloud environment. You access this data the same way as you would in the production - with the official [Firebase Admin SDK](https://firebase.google.com/docs/admin/setup). There's no need to create separate Firebase projects or maintain local scripts to test your functions so you don't corrupt your production data.

- **Continuous feedback**: Pre-define with what data should each Firebase Function be triggered in the [config file](#config-file). The functions are then automatically triggered with every code change. This ensures that you always know whether your functions behave correctly against your production data. There isn't any context switching and no need to leave your coding editor.

- **Discover production bugs**: TODO


## Table of contents
- **[How Foundry works](#how-foundry-works)**
- **[What Foundry doesn't do](#what-foundry-doesnt-do)**
- **[Download and installation](#download-and-installation)**
- **[Supported languages](#supported-languages)**
- **[Config file `foundry.yaml`](#config-file-foundryyaml)**
  - **[Field `functions`](#field-functions)**
  - **[Field `firestore`](#field-firestore)**
  - **[Field `users`](#field-users)**
  - **[Field `ignore`](#field-ignore)**
  - **[Field `serviceAcc`](#field-serviceAcc)**  
  - **[Full config file example](#full-config-file-example)**
- **[How to use the Foundry CLI](#how-to-use-the-foundry-cli)**
  - **[Initialization](#initalization)**
  - **[Interactive prompt](#interactive-prompt)**
  - **[Environment variables](#environment-variables)**
- **[Supported Firebase features](#supported-firebase-features)**
- **[Examples](#examples)**
- **[FAQ](#faq)**
- **[Slack community](#slack-community)**

## How Foundry works
Foundry helps you to develop Firebase Functions faster and with a copy of your production data. It's a command-line tool that connects you to your own cloud environment. This cloud development environment is as much as possible similar to the actual environment where you Firebase Functions run after the deployment and have everything pre-configured.<br/>
Once connected to your development environment, Foundry starts watching your code for changes. Every change notifies the CLI and your code is uploaded to your cloud development environment. Foundry then triggers all your Firebase Functions according to rules specified in your `foundry.yaml` config file. <br/>
Both `stdout` and `stderr` of your functions is sent back to you after each of such runs. The whole upload loop with the transmition of data back to you usually takes about 1-2 seconds. This loop creates a REPL-like tool for you functions and makes it really easy to be sure that your functions' code behave correctly with the production data.<br/>

The [config](#config-file-foundryyaml) `foundry.yaml` file is a critical part of Foundry. There you describe 3 main things:
1. [What Firebase Functions should Foundry register and **how** it should trigger them in each run](#field-functions)
2. [How should the emulated Firestore database look like](#field-firestore)
3. [How should the emulated Firebase Auth users look like](#field-auth)
<br/>

Having an emulated Firestore database and Firebase Auth users in your development environment makes it easy to test new things. 

## What Foundry doesn't do
Foundry doesn't deploy your Firebase Functions onto the production. To deploy functions onto the production use the official [Firebase tool](https://github.com/firebase/firebase-tools).

## Download and installation

- **[macOS](https://github.com/FoundryApp/foundry-cli/releases)**

- **[Linux](TODO)**

To install Foundry, Add the downloaded binary to one of the folders in your system's `PATH` variable.

## Supported languages
JavaScript

## Config file `foundry.yaml`
For Foundry to work, it requires that its config file - `foundry.yaml` - is present. The config file must always be placed next to the Firebase Function's `package.json` file.<br/>
You can run `$ foundry init` to generate a basic config file. Make sure to call this command from a folder where your `package.json` for your Firebase Functions is placed.<br/>

Check out the full config file example [here](#full-config-file-example).

### Field `functions`
An array describing your Firebase functions that should be evaluated by Foundry. All described functions must be exported in the function's root index.js file. In this array, you describe how Foundry should trigger each function in every run.<br/>

Currently, Foundry supports following Firebase functions:
#### HTTPS Functions
Equivalent of - [https://firebase.google.com/docs/functions/http-events]((https://firebase.google.com/docs/functions/http-events))
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

It's important to understand how trigger functions work in Foundry. Everything happens against the emulated Firestore database or the emulated Firebase Auth users. Both can be specified in the config file under fields `firestore` and `users` respectively. So all of your code where you manipulate with Firestore or Firebase Auth happens against the emulated Firestore database and emulated Firebase Auth.<br/>
The same is true for function triggers you describe in the Foundry config file. The triggers usually describe how should the emulated Firestore database or emulated Firebase Auth users be mutated. In return, these mutations will trigger your functions.

### Field `firestore`
The field `firestore` gives you an option to have a separate Firestore database from your production Firestore database. This separate Firestore is an emulated Firestore database that lives in your cloud environment for the duration of your session. <br/>

The `firestore` field expects an array of collections.
You have two options how to fill an emulated Firestore database. 

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

To create a nested collections specify full collection's path
```yaml
firestore:
  - collection: my/nested/collection
    docs: ...
```

### Field `users`
Same as you can [emulate](#field-firestore) Firestore database for your development you can also emulate Firebase Auth users.<br/>

You have two options how to fill the emulated uses:
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
Often there are files and folders that you don't want to upload are watch. To ignore these, you can use an array of glob patterns. For example:
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
For Foundry to able to copy some of your production data to your development cloud environment it must have an access to your Firebase project. This is done through a [service account](https://firebase.google.com/support/guides/service-accounts).
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

  # You can use both the direct and 'geFromProd' approach simultaneously
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
  # Foundry currently supports following types of Firebase Functions:
  # - https
  # - httpsCallable
  # - auth
  # - firestory
  
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
All available commands can be found by calling `$ foundry --help`.

``

### Initialization
TODO

### Interactive prompt
TODO

### Environment variables
TODO


## Supported Firebase features
TODO

## FAQ

### How long is my cloud development environment active after I end the session?
The cloud environment exists only for the time being your session is active. Once your session ends, the environment is terminated.

### Are you storing my code or data?
We don't store your code or any data for a duration longer than is the lifetime duration of your session. Once your session ends, your cloud environment is terminated and only metadata (environment variables) is preserved.

### Why do you need a service account to my Firebase project?
TODO

### How do I get a service account JSON for my Firebase project?
Go to [Google Cloud Console](https://console.cloud.google.com/iam-admin/serviceaccounts) and choose your project

or from [Firebase Console](https://console.firebase.google.com/project)

## License
[Mozilla Public License v2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE)
