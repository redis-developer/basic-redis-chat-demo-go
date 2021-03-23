// @ts-check
import React from "react";
import {login, MESSAGES_TO_LOAD} from "../../../../api";
import InfoMessage from "./components/InfoMessage";
import MessagesLoading from "./components/MessagesLoading";
import NoMessages from "./components/NoMessages";
import ReceiverMessage from "./components/ReceiverMessage";
import SenderMessage from "./components/SenderMessage";

const MessageList = ({
  messageListElement,
  messages,
  room,
  onLoadMoreMessages,
  user = {},
  onUserClicked,
}) => {
  return (
  <div
    ref={messageListElement}
    className="chat-box-wrapper position-relative d-flex"
  >
    {messages === undefined ? (
      <MessagesLoading />
    ) : messages.length === 0 ? (
      <NoMessages />
    ) : (
      <></>
    )}
    <div className="px-4 pt-5 chat-box position-absolute">
      {messages && messages.length !== 0 && (
        <>
          {room.offset && room.offset >= MESSAGES_TO_LOAD ? (
            <div className="d-flex flex-row align-items-center mb-4">
              <div
                style={{ height: 1, backgroundColor: "#eee", flex: 1 }}
              ></div>
              <div className="mx-3">
                <button
                  aria-haspopup="true"
                  aria-expanded="true"
                  type="button"
                  onClick={onLoadMoreMessages}
                  className="btn rounded-button btn-secondary nav-btn"
                  id="__BVID__168__BV_toggle_"
                >
                  Load more
                </button>
              </div>
              <div
                style={{ height: 1, backgroundColor: "#eee", flex: 1 }}
              ></div>
            </div>
          ) : (
            <></>
          )}
          {messages.map((message, x) => {
            const key = message.Message + message.CreatedAt + message.SenderUUID + x;
            if (message.SenderUUID === "info") {
              return <InfoMessage key={key} message={message.Message} />;
            }
            if (message.SenderUUID !== user.uuid) {
              return (
                <SenderMessage
                  onUserClicked={() => onUserClicked(message.SenderUUID)}
                  key={key}
                  message={message.Message}
                  date={message.CreatedAt}
                  user={message.Sender}
                />
              );
            }
            return (
              <ReceiverMessage
                username={
                  message.Sender ?  message.Sender.Username : ""
                }
                key={key}
                message={message.Message}
                date={message.CreatedAt}
              />
            );
          })}
        </>
      )}
    </div>
  </div>
)};
export default MessageList;
