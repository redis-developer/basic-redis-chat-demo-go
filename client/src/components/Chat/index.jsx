// @ts-check
import React, {useCallback, useEffect, useRef, useState} from "react";
import ChatList from "./components/ChatList";
import MessageList from "./components/MessageList";
import TypingArea from "./components/TypingArea";
import useAppStateContext, {useAppState} from '../../state';

/**
 * @param {{
 *  onLogOut: () => void,
 *  onMessageSend: (message: string, roomId: string) => void,
 *  user: import("../../state").UserEntry
 * }} props
 */
export default function Chat({ onLogOut, user, onMessageSend, users }) {
  const [{rooms, currentRoom}] = useAppState();

  const [room, setRoom] = useState({});
  const [message, setMessage] = useState("");

  const messageListElement = useRef(null);

  const scrollToBottom = useCallback(() => {
    if (messageListElement.current) {
      messageListElement.current.scrollTo({
        top: messageListElement.current.scrollHeight,
      });
    }
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [rooms[currentRoom].messages, scrollToBottom]);

  useEffect(() => {
    setRoom(rooms[currentRoom]);
  }, [currentRoom])

  return (
    <div className="container py-5 px-4">
      <div className="chat-body row overflow-hidden shadow bg-light rounded">
        <div className="col-4 px-0">
          <ChatList
            user={user}
            onLogOut={onLogOut}
            rooms={rooms}
            currentRoom={currentRoom}
          />
        </div>
        {/* Chat Box*/}
        <div className="col-8 px-0 flex-column bg-white rounded-lg">
          <div className="px-4 py-4" style={{ borderBottom: "1px solid #eee" }}>
            <h2 className="font-size-15 mb-0">
              {room ? room.name : "Room"}
              {" Room"}
            </h2>
          </div>
          <MessageList
            messageListElement={messageListElement}
            messages={rooms[currentRoom].messages}
            room={rooms[currentRoom]}
            // onLoadMoreMessages={onLoadMoreMessages}
            user={user}
            // onUserClicked={onUserClicked}
          />

          {/* Typing area */}
          <TypingArea
            message={message}
            setMessage={setMessage}
            onSubmit={(e) => {
              e.preventDefault();
              if (message) {
                onMessageSend(message.trim(), room.id);
                setMessage("");

                messageListElement.current.scrollTop =
                    messageListElement.current.scrollHeight;
              }
            }}
          />
        </div>
      </div>
    </div>
  );
}
