const functions = require('firebase-functions');
const admin = require('firebase-admin');
// Import Foundry
const foundry = require('@foundryapp/foundry-backend').firebase;

admin.initializeApp();

// Fill in the emulated auth users
foundry.users.add([
  {
    id: 'user-id-1',
    data: { email: 'user@email.com' },
  },
]);

// Fill in the emulated Firestore
foundry.firestore.collection('posts').addDocs([
  {
    id: 'post-doc-id-1',
    data: {
      ownerId: 'user-id-1',
      content: 'Hello World!',
    },
  },
]);

// Register 'myCloudFunc' with Foundry
// The name under which your cloud function is registered
// with Foundry must be the same under which you export
// your cloud function
const createPost = foundry.functions.httpsCallable.register('createPost');

// Now specify how Foundry should trigger your function
createPost.triggerAsUser('user-id-1').onCall({
  data: {
    content: 'Content of a new post',
  },
});

// Cloud Function for creating posts
exports.createPost = functions.https.onCall(async (data, context) => {
  if (!context.auth) {
    throw new functions.https.HttpsError('permission-denied', 'User isn\'t authenticted');
  }

  const { uid } = context.auth;
  await admin.firestore().collection('posts').add({
    ownerId: uid,
    content: data.content,
  });
});

/////////
const getPosts = foundry.functions.httpsCallable.register('getPosts');
getPosts.triggerAsUser('user-id-1').onCall({ data: {} });

// Cloud Function for retrieving all user's posts
exports.getPosts = functions.https.onCall(async (data, context) => {
  if (!context.auth) {
    throw new functions.https.HttpsError('permission-denied', 'User isn\'t authenticted');
  }

  const { uid } = context.auth;
  const postDocs = await admin.firestore().collection('posts').where('ownerId', '==', uid).get();
  return postDocs.docs.map(d => d.data());
});