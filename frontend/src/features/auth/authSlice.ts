import { createSlice } from "@reduxjs/toolkit"
import type { PayloadAction } from "@reduxjs/toolkit"

interface AuthState {
    username: string | null
    token: string | null
    isAuthenticated: boolean
}

const initialState: AuthState = {
    username: localStorage.getItem("username"),
    token: localStorage.getItem("token"),
    isAuthenticated: !!localStorage.getItem("token"),
}

const authSlice = createSlice({
    name: "auth",
    initialState,
    reducers: {
        loginSuccess: (state, action: PayloadAction<{ username: string; token: string }> ) => {
            state.username = action.payload.username
            state.token = action.payload.token
            state.isAuthenticated = true
            localStorage.setItem("username", action.payload.username)
            localStorage.setItem("token", action.payload.token)
        },
        logout: (state) => {
            state.username = null
            state.token = null
            state.isAuthenticated = false
            localStorage.removeItem("username")
            localStorage.removeItem("token")
        }
    }   
})

export const { loginSuccess, logout } = authSlice.actions
export default authSlice.reducer
