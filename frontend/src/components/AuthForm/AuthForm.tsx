import React from 'react';
import style from './style.module.css';
import SteamAuth from "../SteamAuth/SteamAuth";

function AuthForm() {
  return (
    <>
      <h1 className={style.title}>Авторизация</h1>
      <form method={'POST'} action={'/api/login' + window.location.search} className={style.auth_form}>
        <input placeholder={"Логин"} name={'login'} />
        <input placeholder={"Пароль"} name={'password'} type={'password'}/>
        <input type={'submit'} value={'Войти'}/>
      </form>
      <SteamAuth/>
    </>
  );
}

export default AuthForm;
