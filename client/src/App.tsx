import { Route, Routes } from "react-router-dom"
import { Toaster } from "./components/ui/sonner"
import Home from "./pages/home"
import SendEmailForm from "./pages/email-form-page"

function App() {


  return (
    <div>
    <Routes>
    <Route path="/" element={<Home />} />
    <Route path="/email" element={<SendEmailForm />} />
    </Routes>
    <Toaster />
    </div>
  )
}

export default App
