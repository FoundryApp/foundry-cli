const functions = require('firebase-functions');
const admin = require('firebase-admin');

admin.initializeApp();

exports.getUserEnvs = functions.https.onCall(async (data, context) => {
  const envsDoc = await admin.firestore()
    .collection('envs')
    .doc(context.auth.uid)
    .get();
  return envsDoc.data().envs;
});
