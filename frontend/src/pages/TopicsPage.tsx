import {useEffect, useState} from "react";
import {useApiClientContext} from "../provider/ApiClientProvider.tsx";
import DataTable from "../components/DataTable.tsx";

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
            {
                error ?
                    <p>{error}</p>
                :
                    topics.length ?
                        <div>
                            <div>
                                <DataTable topics={topics}/>
                            </div>
                        </div>
                    :
                        <div>
                            <h2> No Data </h2>
                        </div>
            }
        </div>
    )
}

export default TopicsPage