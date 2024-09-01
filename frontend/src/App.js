import "./App.css";
import Header from "./Header.js"
import Footer from "./Footer.js"

function App() {
  return (
    <div className="app">
      <Header />
      <form>
        <input placeholder="Enter your name..." type="text" className="input-field"></input>
        <button className="btn">Test</button>
      </form>
      <Footer/>
    </div>
  );
}

export default App;
