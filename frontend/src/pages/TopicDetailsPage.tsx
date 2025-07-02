import React, {useEffect, useMemo, useState} from "react";
import {useApiClientContext} from "../provider/ApiClientProvider.tsx";
import {TopicDetails} from "../model/TopicDetails.ts";
import {PartitionDetails} from "../model/PartitionDetails.ts";
import {mapModelMessageToMessage} from "../service/MessageService.ts";
import {MessageDetails} from "../model/MessageDetails.ts";
import MessageDataTable from "../components/MessageDataTable.tsx";
import {Link, useLocation} from "react-router-dom";
import TopicDetailsOptionBar, {partitionNumberOption} from "../components/TopicDetailsOptionBar.tsx";

const defaultStartOffset = -50;
const defaultEndOffset = -1;

const TopicDetailsPage: React.FC = () => {
    const location = useLocation()
    const topic = location.state?.topicDetails as TopicDetails | undefined;
    const initialPartitionDetails: PartitionDetails[] = Array.from({ length: topic?.numPartitions ?? 0 }, (_, i) => ({
        partition: i,
        startOffset: defaultStartOffset,
        endOffset: defaultEndOffset
    }));

    const [partitionDetails, setPartitionDetails] = useState<PartitionDetails[]>(initialPartitionDetails)
    const [partitionsToShow, setPartitionsToShow] = useState<number[]>(initialPartitionDetails.map(partition => partition.partition))
    const [messages, setMessages] = useState<MessageDetails[]>([])
    const [error, setError] = useState<string | null>(null)
    const apiClient = useApiClientContext()
    const partitionProps = useMemo(() =>
        partitionDetails.map(partition => partition.partition),
        [partitionDetails]
    );

    const handlePartitionNumberChange = (partitionNumber: partitionNumberOption) => {
        if (partitionNumber === 'All') {
            setPartitionsToShow(partitionDetails.map(partition => partition.partition));
        } else {
            setPartitionsToShow([partitionNumber])
        }
    }

    useEffect(() => {
        if (!topic) {
            setError("No topic data found. Please navigate from the topics list.");
            return;
        }
        const apiRequest = {
            topicMessagesInput: {
                topicName: topic.topic,
                partitions: partitionDetails
            }
        }
        apiClient.topicsMessagesPost(apiRequest).then(
            (response) => {
                const mappedResponse = mapModelMessageToMessage(response)
                setMessages(mappedResponse)
            }
        ).catch(
            (error) => {
                setError(error.message)
            }
        )
    }, [apiClient, partitionDetails, topic])

    const messagesToShow = useMemo(() =>
        messages.filter(message => (
            partitionsToShow.includes(message.partition)
        )),
        [messages, partitionsToShow]
    )

    return (
        <div>
            <h1>{topic?.topic ?? "Topic Not Found"}</h1>
            {
                error ?
                    <div>
                        <h1>Error</h1>
                        <p>{error}</p>
                        <Link to="/topics">← Back to Topics</Link>
                    </div>
                :
                    partitionDetails.length ?
                        <div>
                            <TopicDetailsOptionBar
                                partitions={partitionProps}
                                onPartitionChange={handlePartitionNumberChange}
                            />
                            <div>
                                <MessageDataTable messages={messagesToShow}/>
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

export default TopicDetailsPage