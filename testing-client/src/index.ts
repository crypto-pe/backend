import fetch from "cross-fetch";
import { API, GetSupportedTokensReturn } from "./api.gen";

const client = new API("http://localhost:8000", fetch);

client.ping().then((something) => console.log(something));
client
  .getSupportedTokens()
  .then((something: GetSupportedTokensReturn) =>
    something.tokens.forEach((token) => console.log(token))
  );

// const authHeaders = {
//     Authorization: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHAiOiJldm9zdmVyc2UiLCJleHAiOjE2ODI4Nzg4NjAsImlhdCI6MTY1MTMyODQ2MH0.KkgXsEQLBCP8e8GrBpUHwNeWvbx60TL4tYqR6u7AUC8Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50IjoiMHgwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwN2QwIn0.FjLk2uKwKJOvhfa61qzvUJxwZs_qWl6AqjpWR3QDHkQ'
// }

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
