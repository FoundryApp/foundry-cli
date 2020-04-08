const functions = require('firebase-functions');
const admin = require('firebase-admin');

admin.initializeApp();
let FieldValue = admin.firestore.FieldValue;

exports.getUserEnvs = functions.https.onCall(async (data, context) => {
  const envsDoc = await admin.firestore()
    .collection('envs')
    .doc(context.auth.uid)
    .get();
  return envsDoc.data().envs;
});

exports.deleteUserEnvs = functions.https.onCall(async (data, context) => {
  const toDeleteArr = data.delete;
  if (!toDeleteArr) {
    throw new functions.https.HttpsError('invalid-argument', `Expected "delete" array in the body. Got: ${toDeleteArr}`);
  }

  const envsDocRef = admin.firestore()
    .collection('envs')
    .doc(context.auth.uid);

  const currentEnvs = (await envsDocRef.get()).data().envs;

  console.log(`Current envs: "${Object.keys(currentEnvs)}", for user "${context.auth.uid}"`);
  console.log(`Will delete envs "${toDeleteArr}"`);

  const newEnvs = {}
  Object.keys(currentEnvs).forEach(envName => {
    if (!toDeleteArr.includes(envName)) {
      newEnvs[envName] = currentEnvs[envName]
    }
  });

  try {
    await envsDocRef.update({ envs: newEnvs });
    console.log("New envs:", newEnvs);
  } catch (error) {
    throw new functions.https.HttpsError('internal', `Error updating user envs: ${error}`);
  }
});
