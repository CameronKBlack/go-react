import "./Header.css"

export default function Header(){
    return (
      <header className="app-header">
        <TopBar/>
      </header>
    )
  }
  
  function TopBar(){
    return (
      <div>
        <button className="login-btn">Login</button>
      </div>
  )
}