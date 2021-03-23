// @ts-check
import {messages} from "./websockets/messages";
import {webSocketUrl} from './websockets/config';

const ws = new WebSocket(webSocketUrl);

export function initWebSocket(appState) {
  ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('EVENT', data.type, data);
    if(typeof messages[data.type] == "function") {
      messages[data.type](data, appState);
    } else {
      console.log("Unknown message type: " + data.type);
    }
  }
}

/**
 * main method for send message to websocket server
 * @param {Object} message
 */
export function webSocketSend(message) {
  const wsData = JSON.stringify(message);
  ws.send(wsData);
}

