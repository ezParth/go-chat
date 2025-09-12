import { toast } from "react-hot-toast"

export const showSuccess = (message: string) => {
  toast.success(message, {
    position: "top-right",
    duration: 3000,
  })
}

export const showError = (message: string) => {
  toast.error(message, {
    position: "top-right",
    duration: 3000,
  })
}

export const showInfo = (message: string) => {
  toast(message, {
    position: "top-right",
    duration: 3000,
  })
}
