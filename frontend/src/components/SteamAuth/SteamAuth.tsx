import React from 'react';

import style from './style.module.css';
import {setCookie} from "../../utils/cookie";

export default function SteamAuth() {
  const SteamLogin = () => {
    const realm = window.location.origin;
    const returnTo = `${realm}/api/steam` + window.location.search;

    const nonce = crypto.randomUUID ? crypto.randomUUID() : Math.random().toString(36).substring(2);

    const params = new URLSearchParams({
      'openid.ns': 'http://specs.openid.net/auth/2.0',
      'openid.mode': 'checkid_setup',
      'openid.return_to': returnTo,
      'openid.realm': realm,
      'openid.identity': 'http://specs.openid.net/auth/2.0/identifier_select',
      'openid.claimed_id': 'http://specs.openid.net/auth/2.0/identifier_select',
    });

    setCookie('steam_nonce', nonce, 1);

    window.location.href = `https://steamcommunity.com/openid/login?${params}`;
  };

  return <div>
    <button title={'Войти через Steam'} className={style.button} onClick={SteamLogin}>
    <span>
      <img
        width={'30px'}
        src={'/icons8-steam.svg'}
        alt={'SteamAuth icon'}/>
    </span>
      <span>Steam</span>
    </button>
  </div>;
}