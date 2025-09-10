import { BrowserRouter, Routes, Route } from "react-router-dom"
import Chat from "./pages/Chat"
import Login from "./pages/Login"
import CreateGroup from "./pages/CreateGroup"
import GroupChat from "./pages/GroupChat"

const App = () => {
  return (
    <div>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Chat />} />
          <Route path="/login" element={<Login />} />
          <Route path="/createGroup" element={<CreateGroup />} />
          <Route path="/groupChat/:groupName" element={<GroupChat />} />
        </Routes>
      </BrowserRouter>
    </div>
  )
}

export default App
