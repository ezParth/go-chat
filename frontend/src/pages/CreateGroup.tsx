import { useEffect, useState } from "react"
import { groupApi } from "../api/group.ts"
import { useSelector } from "react-redux"
import type { RootState } from "../store/store"
import { useNavigate } from "react-router-dom"

const CreateGroup = () => {
  const [groupName, setGroupName] = useState<string>("")
  const { isAuthenticated, token } = useSelector((state: RootState) => state.auth)
  const nav = useNavigate()

  useEffect(() => {
    if (!isAuthenticated || !token) {
      nav("/login")
    }
  }, [isAuthenticated, token, nav])

  const CreateNewGroup = async () => {
    if (groupName.trim() === "") {
      alert("Group's name cannot be empty")
      return
    }

    if (!token) return
    try {
      const res = await groupApi.createGroup(groupName)
      if (res?.data.success) {
        console.log("Group created -> ", res.data)
        nav(`/groupChat/${groupName}`)
      } else {
        console.log("Error -> ", res?.data)
      }
    } catch (err) {
      console.error("Error creating group", err)
    }
  }

  return (
    <div style={{ padding: "2rem" }}>
      <h2>Create New Group</h2>
      <input
        type="text"
        placeholder="Enter Group Name"
        value={groupName}
        onChange={(e) => setGroupName(e.target.value)}
      />
      <button onClick={CreateNewGroup} style={{ marginLeft: "1rem" }}>
        Create
      </button>
    </div>
  )
}

export default CreateGroup
