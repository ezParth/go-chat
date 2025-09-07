import { useState } from "react"
import axios from "axios"

const API = "http://localhost:8080"

const handleLogin = async (username: string, password: string) => {
    try {
        const res = await axios.post(`${API}/login`, {
            username: username,
            password: password
        })

        if(res.data?.success) {
            localStorage.setItem("username", username)
            localStorage.setItem("token", res.data?.token)
        }else {
            console.log("Error in login -> ", res.data)
        }
    } catch (error) {
        console.log("Error during login", error)
    }
}

const handleSignup = async (username: string, email: string, password: string) => {
    try {
        const res = await axios.post(`${API}/signup`, {
            username: username,
            email: email,
            password: password
        })

        if(res.data?.success) {
            console.log("Logged in Successfully")
            localStorage.setItem("username", username)
            localStorage.setItem("token", res.data?.token)
        }else {
            console.log("error -> ", res.data)
        }
    } catch (error) {
        console.log("Error during singup", error)
    }
}

const Login = () => {
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [email, setEmail] = useState("")
    const [islogin, setIslogin] = useState(true)

    const handleChangeLogin = () => {
        setIslogin(!islogin)
    }

    const handleSubmit = async () => {
        if(islogin) {
            await handleLogin(username, password)
        }else {
            await handleSignup(username, email, password)
        }
    }

    return (
        <div>
            <input type="text" placeholder="Enter Name" value={username} onChange={(er) => setUsername(er.target.value)}  />
            {!islogin && 
            <input placeholder="Enter Email" value={email} onChange={(er) => setEmail(er.target.value)} />
            }
            <input type="text" placeholder="Enter Password" value={password} onChange={(er) => setPassword(er.target.value)} />
            <button onClick={handleSubmit}>Submit</button>
            <p>{islogin ? "Don't have an account  " : "Already Have an Account  "}<button onClick={handleChangeLogin}>{islogin ? "Signup" : "Login"}</button></p>
        </div>
    )
}

export default Login