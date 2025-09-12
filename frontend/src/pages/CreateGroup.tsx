import { useEffect, useState } from "react"
import { groupApi } from "../api/group.ts"
import { useSelector } from "react-redux"
import type { RootState } from "../store/store"
import { useNavigate } from "react-router-dom"
import { avatarlink } from "../assets/Image.ts"

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
    let res
    try {
      console.log("Making an API call to go-server")
      res = await groupApi.createGroup(groupName, avatarlink)
      console.log("Create Group Response -> ", res)
      if (res?.data.success) {
        console.log("Group created -> ", res.data)
        nav(`/groupChat/${groupName}`)
      } else {
        console.log("Error -> ", res?.data)
      }
    } catch (err) {
      console.error("Error creating group", err)
    } finally {
      console.log("Create Group Response -> ", res)
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
