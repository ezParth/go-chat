import { useState } from "react"

const GroupChat = () => {
    const [peopleOnline, setPeopleOnline] = useState<string[]>([])
    return(
        <div>
            <div>
                {peopleOnline.map((val, key) => (
                    <li key={key}>{val}</li>
                ))}
            </div>
        </div>
    )
}

export default GroupChat