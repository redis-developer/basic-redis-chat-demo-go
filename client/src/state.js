// @ts-check
import { createContext, useContext, useReducer } from "react";

/**
 * @typedef {{
 *  from: string
 *  date: number
 *  message: string
 *  roomId?: string
 * }} Message
 *
 * @typedef {{
 *   name: string;
 *   id: string;
 *   messages?: Message[]
 *   connected?: boolean;
 *   offset?: number;
 *   forUserId?: null | number | string
 *   lastMessage?: Message | null
 * }} Room
 *
 * @typedef {{
 *   username: string;
 *   id: string;
 *   online?: boolean;
 *   room?: string;
 * }} UserEntry
 *
 * @typedef {{
 *  currentRoom: string;
 *  rooms: {[id: string]: Room};
 *  users: {[id: string]: UserEntry}
 * }} State
 *
 * @param {State} state
 * @param {{type: string; payload: any}} action
 * @returns {State}
 */
const reducer = (state, action) => {
  switch (action.type) {
    case "clear":
      return { currentRoom: "0", rooms: {}, users: {} };
    case "set user": {
      console.log('set user', action.payload)
      return {...state, user: action.payload};
    }
    case "set users": {
      return {
        ...state,
        users: action.payload,
      };
    }
    case "make user online": {
      return {
        ...state,
        users: {
          ...state.users,
          [action.payload]: { ...state.users[action.payload], online: true },
        },
      };
    }
    case "append users": {
      return { ...state, users: { ...state.users, ...action.payload } };
    }
    case "set messages": {
      return {
        ...state,
        rooms: {
          ...state.rooms,
          [state.currentRoom]: {
            ...state.rooms[state.currentRoom],
            messages: action.payload,
            offset: action.payload.length,
          },
        },
      };
    }
    case "prepend messages": {
      const messages = [
        ...action.payload.messages,
        ...state.rooms[action.payload.id].messages,
      ];
      return {
        ...state,
        rooms: {
          ...state.rooms,
          [action.payload.id]: {
            ...state.rooms[action.payload.id],
            messages,
            offset: messages.length,
          },
        },
      };
    }
    case "append message":
      if (state.rooms[action.payload.id] === undefined) {
        return state;
      }
      return {
        ...state,
        rooms: {
          ...state.rooms,
          [action.payload.id]: {
            ...state.rooms[action.payload.id],
            lastMessage: action.payload.message,
            messages: state.rooms[action.payload.id].messages
              ? [
                ...state.rooms[action.payload.id].messages,
                action.payload.message,
              ]
              : undefined,
          },
        },
      };
    case 'set last message':
      return { ...state, rooms: { ...state.rooms, [action.payload.id]: { ...state.rooms[action.payload.id], lastMessage: action.payload.lastMessage } } };
    case "set current room":
      return { ...state, currentRoom: action.payload };
    case "add room":
      return {
        ...state,
        rooms: { ...state.rooms, [action.payload.id]: action.payload },
      };
    // case "set rooms": {
    //   /** @type {Room[]} */
    //   const newRooms = action.payload;
    //   const rooms = { ...state.rooms };
    //   newRooms.forEach((room) => {
    //     rooms[room.id] = {
    //       ...room,
    //       messages: rooms[room.id] && rooms[room.id].messages,
    //     };
    //   });
    //   return { ...state, rooms };
    // }
    // temporary, while no rooms
      case "set rooms": {
        /** @type {Room[]} */
        const newRooms = action.payload;
        const rooms = { ...state.rooms };
        newRooms.forEach((room) => {
          rooms[room.UUID] = {
            ...room,
            id: room.UUID,
            name: room.Username,
            messages: rooms[room.id] && rooms[room.id].messages,
          };
        });
        return { ...state, rooms };
      }
    default:
      return state;
  }
};

/** @type {State} */
const initialState = {
  currentRoom: "",
  rooms: {"": {id: "", name: "General"}},
  users: {},
  user: {}
};

const useAppStateContext = () => {
  return useReducer(reducer, initialState);
};

// @ts-ignore
export const AppContext = createContext();

/**
 * @returns {[
 *  State,
 *  React.Dispatch<{
 *   type: string;
 *   payload: any;
 * }>
 * ]}
 */
export const useAppState = () => {
  const [state, dispatch] = useContext(AppContext);
  return [state, dispatch];
};

export default useAppStateContext;
