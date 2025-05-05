import style from "./components/AuthForm/style.module.css";
import UserInfo from "./components/UserInfo/UserInfo";
import AuthForm from "./components/AuthForm/AuthForm";
import React, {useEffect, useState} from "react";

export default function App() {
  const [myId, setMyId] = useState(null);

  useEffect(() => {
    fetch("/api/me", {
      credentials: 'include',
    })
      .then(res => res.json())
      .then(res => setMyId(res.user_id));
  }, []);

  return <div className={style.auth_page}>
    {myId === null ? "Load" : (myId === undefined ? <AuthForm/> : <UserInfo myId={myId}/>)}
  </div>
}