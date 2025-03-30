import {useEffect, useState} from "react";
import {useApiClientContext} from "../provider/ApiClientProvider.tsx";

function TopicsPage() {
    const [topics, setTopics] = useState<string[]>([])
    const [error, setError] = useState<string | null>(null)
    const apiClient = useApiClientContext()

    useEffect(() => {
        apiClient.topicsGet().then(
            (response) => {
                setTopics(response)
            }
        ).catch(
            (error) => {
                setError(error.message)
            }
        )
    }, [apiClient])

    return (
        <div>
            <h1>Topics</h1>
            {error && <p>{error}</p>}
            <ul>
                {topics.map((topic) => <li key={topic}>{topic}</li>)}
            </ul>
        </div>
    )
}

export default TopicsPage