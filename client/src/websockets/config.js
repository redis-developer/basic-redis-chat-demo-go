export const webSocketUrl = "ws"+(window.location.protocol.substr(0,5)==="https"?"s":"")+"://"+ process.env.REACT_APP_CHAT_BACKEND +"/ws";
//export const webSocketUrl = "wss://"+ process.env.REACT_APP_CHAT_BACKEND +"/ws";
console.log("value of REACT_APP_CHAT_BACKEND is : " + process.env.REACT_APP_CHAT_BACKEND);
console.log("value of HTTP_PROXY is : " + process.env.HTTP_PROXY);
console.log("value of webSocketUrl is : " + webSocketUrl);
console.log("value of window.location.protocol is : " + window.location.protocol);
console.log("value of CLIENT_LOCATION is : " + process.env.CLIENT_LOCATION);
console.log("value of SERVER_ADDRESS is : " + process.env.SERVER_ADDRESS);