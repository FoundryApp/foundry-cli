const foundry = require('@foundryapp/foundry-backend').firebase;

// Fill in auth
foundry.users.add([
  {
    id: '',
    data
  },
]);
foundry.users.copyFromProdById(['DbJD37dhx4VqNF3dtN7zrK6CCB13']);
// Fill in Firestore
foundry.firestore.collection('envs').copyDocsFromProdByCount(5);

////// getUserEnvs
const getUserEnvs = foundry.functions.httpsCallable.register('getUserEnvs');
getUserEnvs.triggerAsUser('DbJD37dhx4VqNF3dtN7zrK6CCB13').onCall({ data: {} });

////// deleteUserEnvs
const deleteUserEnvs = foundry.functions.httpsCallable.register('deleteUserEnvs');
deleteUserEnvs.triggerAsUser('DbJD37dhx4VqNF3dtN7zrK6CCB13').onCall({
  data: {
    delete: ['env1', 'env2'],
  },
});
