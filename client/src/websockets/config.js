export const webSocketUrl = "ws"+(window.location.protocol.substr(0,5)==="https"?"s":"")+"://"+ process.env.REACT_APP_CHAT_BACKEND +"/ws";
//export const webSocketUrl = "wss://"+ process.env.REACT_APP_CHAT_BACKEND +"/ws";