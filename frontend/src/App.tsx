import { BrowserRouter, Routes, Route } from "react-router-dom"
import Chat from "./pages/Chat"
import Login from "./pages/Login"

const App = () => {
  return(
    <div>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Chat />} />
          <Route path="/login" element={<Login />} />
        </Routes>
      </BrowserRouter>
    </div>
  )
}

export default App
