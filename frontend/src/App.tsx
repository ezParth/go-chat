import { BrowserRouter, Routes, Route } from "react-router-dom"
import Chat from "./pages/Chat"

const App = () => {
  return(
    <div>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Chat />} />
        </Routes>
      </BrowserRouter>
    </div>
  )
}

export default App
