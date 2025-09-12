import { BrowserRouter, Routes, Route } from "react-router-dom"
import Chat from "./pages/Chat"
import Login from "./pages/Login"
import CreateGroup from "./pages/CreateGroup"
import GroupChat from "./pages/GroupChat"
import { Toaster } from "react-hot-toast"
import JoinGroup from "./pages/JoinGroup"

const App = () => {
  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Chat />} />
          <Route path="/login" element={<Login />} />
          <Route path="/group/create" element={<CreateGroup />} />
          <Route path="/group/join" element={<JoinGroup />} />
          <Route path="/groupChat/:groupName" element={<GroupChat />} />
        </Routes>
      </BrowserRouter>
      
      <Toaster position="top-right" reverseOrder={false} />
    </>
  )
}

export default App
