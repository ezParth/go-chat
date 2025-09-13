import axios from "axios";
import { store, type RootState } from "../store/store";

const API_URL = "http://localhost:8080/group"; // Change to your backend URL

// Create Axios instance
const api = axios.create({
  baseURL: API_URL,
});

api.interceptors.request.use((config) => {
    const state: RootState = store.getState();
    let token  = null
    token = state.auth?.token
if (!token) {
    token = localStorage.getItem("token");
}

if(!token) {
    throw new Error("Cannot Find Token, Please Login")
}
//   const username = state.auth?.username;

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

//   if (username) {
//     config.headers["X-Username"] = username;
//   }

  return config;
});

export const groupApi = {
  // createGroup: (groupName: string, avatar?: File) => {
  //   const formData = new FormData();
  //   formData.append("groupname", groupName);
  //   if (avatar) formData.append("avatar", avatar);
  //   return api.post("/create", formData);
  // },

  createGroup: (groupName: string, avatar?: string) => {
    return api.post("/create", {groupName: groupName, avatar: avatar})
  },

  joinGroup: (groupName: string) => {
    return api.post("/join", { groupname: groupName });
  },

  getGroups: () => {
    return api.get("/getGroups");
  },

  getGroupChat: (groupName: string) => {
    return api.get(`/chats/${groupName}`);
  },

  getGroupAvatar: (groupName: string) => {
    return api.get(`/avatar/${groupName}`, { responseType: "blob" });
  },

  getGroupMembersAndAdmin: (groupName: string) => {
    return api.get(`/members/${groupName}`);
  },

  deleteGroup: (groupName: string) => {
    return api.delete("/delete", { data: { groupname: groupName } });
  },

  getUsersByGroupName: (groupName: string) => {
    return api.get(`/online/${groupName}`)
  }
};
