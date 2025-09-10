import { useState } from "react"
import axios from "axios"
import { useNavigate } from "react-router-dom"
import { useDispatch } from "react-redux"
import { loginSuccess } from "../features/auth/authSlice"
import { API } from "../api/api"

const handleLogin = async (username: string, password: string) => {
    try {
        const res = await axios.post(`${API}/login`, {
            username: username,
            password: password
        })

        if(res.data?.success) {
            console.log("res.data -> ", res.data)
            localStorage.setItem("username", username)
            localStorage.setItem("token", res.data?.token)
            return res
        }else {
            console.log("Error in login -> ", res.data)
            return res
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
            console.log("Signup in Successfully")
            console.log(res.data)
            localStorage.setItem("username", username)
            localStorage.setItem("token", res.data?.token)
            return res
        }else {
            console.log("error -> ", res.data)
            return res
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

    const nav = useNavigate()
    const dispatch = useDispatch()
    const handleSubmit = async () => {
        if(islogin) {
            const res = await handleLogin(username, password)
            if(res == undefined) {
                alert("Cannot login")
                return
            }
            
            if(res.data.success) {
                dispatch(loginSuccess({username, token: res.data.token}))
                nav("/")
            }
        }else {
            const res = await handleSignup(username, email, password)
            if(res == undefined) {
                alert("Cannot signup")
                return
            }

            if(res.data.success) {
                dispatch(loginSuccess({username, token: res.data.token}))
                nav("/")
            }
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