import React, {useEffect, useState} from "react";
import {useApiClientContext} from "../provider/ApiClientProvider.tsx";
import TopicDataTable from "../components/TopicDataTable.tsx";
import {mapModelTopicDetailsToTopicDetails} from "../service/TopicDetailsService.ts";
import {TopicDetails} from "../model/TopicDetails.ts";

const TopicsPage: React.FC = () => {
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
                                <TopicDataTable topics={topics}/>
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