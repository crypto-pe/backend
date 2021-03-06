import fetch from 'cross-fetch'
import { API, GetSupportedTokensReturn } from './api.gen'
import { ethers } from 'ethers'
import {
  ETHAuth,
  Claims,
  validateClaims,
  Proof,
  ETHAuthVersion,
  ValidatorFunc,
  IsValidSignatureBytes32MagicValue,
} from "@0xsequence/ethauth";

const client = new API('http://localhost:8000', fetch)

client.ping().then((something) => console.log(something))

client.ping().then((something) => console.log(something));

const authHeaders = {
  Authorization:
    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50IjoiMHhlMGM5ODI4ZGVlMzQxMWEyOGNjYjRiYjgyYTE4ZDBhYWQyNDQ4OWUwIiwiYXBwIjoiYXBpLXRlc3RpbmctY2xpZW50IiwiZXhwIjoxNjc5NjgxMzYyLCJpYXQiOjE2NTM3NjEzNjJ9.4XFQLahOuGixtoa7I_iL8rRwkN0hPb3bwABVHGSVPgw",
};

client
  .ping(authHeaders)
  .then((something) => console.log(something))
  .catch((err) => console.log(err));
// client
//   .getSupportedTokens()
//   .then((something: GetSupportedTokensReturn) =>
//     something.tokens.forEach((token) => console.log(token))
// );

const wallet = ethers.Wallet.fromMnemonic(
  "outdoor sentence roast truly flower surface power begin ocean silent debate funny"
);

const claims: Claims = {
  app: "api-testing-client",
  iat: Math.round(new Date().getTime() / 1000),
  exp: Math.round(new Date().getTime() / 1000) + 60 * 60 * 24 * 300,
  v: ETHAuthVersion,
};

const proof = new Proof({ address: wallet.address });
proof.claims = claims;
const digest = proof.messageDigest();
const digestHex = ethers.utils.hexlify(digest);
console.log("digestHex", digestHex);

async function prooffunc() {
  proof.signature = await wallet.signMessage(digest)
  const ethAuth = new ETHAuth()
  const proofString = await ethAuth.encodeProof(proof)
  console.log('proofStringReturned', proofString)

  // client
  //   .createAccount({
  //     ethAuthProofString: proofString,
  //     name: 'John Doe',
  //     email: 'johndoe@gmail.com',
  //   })
  //   .then((something) => console.log(something))
  //   .catch((err) => console.log(err))

  client
    .login({
      ethAuthProofString: proofString,
    })
    .then((something) => console.log(something))
    .catch((err) => console.log(err))
}

  // client.login({
  //   ethAuthProofString:  proofString,
  // }).then((something) => console.log(something)).catch((err) => console.log(err))

prooffunc();

client
  .getAccount({
    address:
      '0x307865306339383238646565333431316132386363623462623832613138643061616432343438396530',
  })
  .then((data) => console.log("Account data is", data))

  // // client.getAccount(
// //     {
// //         address: '0xd4Bbf5d234CC95441A8Af0a317D8874eE425e74d',
// //     },
// // ).then(account => console.log('ACCOUNT FOUND ', { account }))
// //     .catch(err => console.log('ACCOUNT NOT FOUND ', err))

// client.ping(authHeaders).then(res => console.log('PING OK', res)).catch(err => console.log('PING ERR', err))

// client.getFeed(
//     {
//         req: {
//             accountAddress: '0xd4Bbf5d234CC95441A8Af0a317D8874eE425e74d',
//         }
//     },
//     authHeaders
// ).then(feed => console.log('FEED FOUND ', { feed }))
//     .catch(err => console.log('FEED NOT FOUND ', err))
