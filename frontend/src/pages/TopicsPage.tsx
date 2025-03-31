import {useEffect, useState} from "react";
import {useApiClientContext} from "../provider/ApiClientProvider.tsx";
import DataTable from "../components/DataTable.tsx";
import {mapModelTopicDetailsToTopicDetails} from "../service/TopicDetailsService.ts";
import {TopicDetails} from "../model/TopicDetails.ts";

function TopicsPage() {
    const [topics, setTopics] = useState<TopicDetails[]>([])
    const [error, setError] = useState<string | null>(null)
    const apiClient = useApiClientContext()

    useEffect(() => {
        apiClient.topicsGet().then(
            (response) => {
                const mappedResponse = mapModelTopicDetailsToTopicDetails(response)
                setTopics(mappedResponse)
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