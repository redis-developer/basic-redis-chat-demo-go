// @ts-check
import "./style.css";
import React, { useMemo } from "react";
import moment from "moment";
import { useEffect } from "react";
import { getMessages } from "../../../../../../api";
import AvatarImage from "../AvatarImage";
import OnlineIndicator from "../../../OnlineIndicator";
import useAppStateContext, {useAppState} from '../../../../../../state';

/**
 * @param {{ active: boolean; room: import('../../../../../../state').Room; onClick: () => void; }} props
 */
const ChatListItem = ({ room, active = false, onClick }) => {
  const { online, name, userId } = useChatListItemHandlers(room);
  return (
    <div
      onClick={onClick}
      className={`chat-list-item d-flex align-items-start rounded ${
        active ? "bg-white" : ""
      }`}
    >
      <div className="align-self-center mr-3">
        <OnlineIndicator online={online} hide={room.id === "0"} />
      </div>
      <div className="align-self-center mr-3">
        <AvatarImage name={name} id={userId} />
      </div>
      <div className="media-body overflow-hidden">
        <h5 className="text-truncate font-size-14 mb-1">{name}</h5>
      </div>
    </div>
  );
};

const useChatListItemHandlers = (
  /** @type {import("../../../../../../state").Room} */ room
) => {
  const { id, name } = room;
  const [state] = useAppState();
  /** Here we want to associate the room with a user by its name (since it's unique). */
  const [isUser, online, userId] = useMemo(() => {
    try {
      let pseudoUserId = Math.abs(parseInt(id.split(":").reverse().pop()));
      const isUser = pseudoUserId > 0;
      const usersFiltered = Object.entries(state.users)
        .filter(([, user]) => user.Username === name)
        .map(([, user]) => user);
      let online = false;
      if (usersFiltered.length > 0) {
        online = usersFiltered[0].OnLine;
        pseudoUserId = +usersFiltered[0].id;
      }
      return [isUser, online, pseudoUserId];
    } catch (_) {
      return [false, false, "0"];
    }
  }, [id, name, state.users]);

  return {
    isUser,
    online,
    userId,
    name: room.name,
  };
};

export default ChatListItem;
