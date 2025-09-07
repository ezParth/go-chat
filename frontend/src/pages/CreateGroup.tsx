import { useEffect, useState } from "react"
import { API } from "../API/api"
import axios from "axios"
import { useSelector } from "react-redux"
import type { RootState } from "../store/store"
import { useNavigate } from "react-router-dom"

const HandleCreateGroup = async (groupName: string, token: string) => {
    try {
        const res = await axios.post(`${API}/createGroup`, {
            groupname: groupName,
            headers: {
                Authorization: `Bearer ${token}`
            }
        })

        return res
    } catch (error) {
        console.log("Error in creating the group -> ",error)
    }
}

// const HandleDeleteGroup = async (groupName: string, token: string) => {
//     try {
//         const res = await axios.post(`${API}/deleteGroup`, {
//             groupName: groupName,
//             headers: {
//                 Authorization: `Bearer ${token}`
//             }
//         })

//         return res
//     } catch (error) {
//         console.log("Error deleting the group -> ", error)
//     }
// }

const CreateGroup = () => {
    const [groupName, setGroupName] = useState<string>("")
    const { isAuthenticated, token } = useSelector((state: RootState) => state.auth)
    const nav = useNavigate()
    useEffect(() => {
        if(!isAuthenticated || !token) {
            nav("/login")
        }
    }, [isAuthenticated, token, nav])

    const CreateNewGroup = async () => {
        if(groupName == "") {
            alert("Group's name cannot be empty")
        }

        if(!token) return
        const res = await HandleCreateGroup(groupName, token)
        if(res?.data.success) {
            console.log("res -> ", res.data)
            nav(`/group/${groupName}`)
        }else {
            console.log("res -> ", res?.data)
        }
    }

    return (
        <div>
            <div>
                <input type="text" placeholder="Enter Group Name" value={groupName} onChange={(er) => setGroupName(er.target.value)} />
                <button onClick={CreateNewGroup}>Create</button>
            </div>
        </div>
    )
}

export default CreateGroup