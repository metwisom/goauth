export default function UserInfo({myId}: { myId: number }) {
  const realm = window.location.origin;

  return <div>UserInfo {myId}<br/>
    <a href={'/api/logout?redirect_uri=' + realm}>Выйти</a>
  </div>
}