import { useState } from "react"
import { groupApi } from "../api/group"
import { showError, showSuccess } from "../servies/toast"
import { useNavigate } from "react-router-dom"

const JoinGroup = () => {
    const [groupName, setGroupName] = useState("")
    const nav = useNavigate()
    const handleJoinGroup = async () => {
        console.log("Hitting HandleJoinGroup")
        if (groupName.trim() === "") {
            showError("Group name cannot be empty")
            return
          }
        const res = await groupApi.joinGroup(groupName)
        console.log("Resoponse -> ",res)
        if (!res?.data.success) {
            showError("Cannot Join the Group")
        }else {
            showSuccess(`Group: ${groupName} joined successfully`)
            nav(`/groupChat/${groupName}`)
        }
    }
    return (
        <div>
            <input type="text" value={groupName} onChange={(er) => setGroupName(er.target.value)} />
            <button onClick={handleJoinGroup}>Join</button>
        </div>
    )
}

export default JoinGroup