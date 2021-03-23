// @ts-check
import React, {useEffect, useState} from "react";
import Login from "./components/Login";
import Chat from "./components/Chat";
import {AppContext} from "./state";
import {LoadingScreen} from "./components/LoadingScreen";
import Navbar from "./components/Navbar";
import {processChannelMessage, processSignIn} from "./websockets/process";
import useAppStateContext from './state';
import {initWebSocket} from './hooks';

const App = () => {
    const [showLogin, setShowLogin] = useState(true)
    const [state, dispatch] = useAppStateContext();

    function onLogout() {
        localStorage.clear();
        setShowLogin(true);
    }

    useEffect(() => {
        initWebSocket({dispatch, state});
    }, []);

    return (
        <AppContext.Provider value={[state, dispatch]}>
            <div
                className={`full-height ${showLogin ? "bg-light" : ""}`}
                style={{
                    backgroundColor: !showLogin ? "#495057" : undefined,
                }}
            >
                <Navbar/>
                {showLogin ? (
                    <Login onLogIn={processSignIn} setShowLogin={setShowLogin}/>
                ) : (
                    <Chat
                        user={state.user}
                        users={state.users}
                        onMessageSend={processChannelMessage}
                        onLogOut={onLogout}
                    />
                )}
            </div>
        </AppContext.Provider>
    );


};

export default App;
