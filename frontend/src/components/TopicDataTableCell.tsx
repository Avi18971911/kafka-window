import React from "react";
import {TopicDetails} from "../model/TopicDetails.ts";
import {useNavigate} from "react-router-dom";

type TopicDataTableCellProps = {
    topicDetails: TopicDetails;
    onToggleExpand: (topic: string) => void;
    expanded: boolean;
};

const TopicDataTableCell: React.FC<TopicDataTableCellProps> = ({ topicDetails, onToggleExpand, expanded }) => {
    const isInternalTopic = topicDetails.isInternal ? 'Yes' : 'No';
    const retentionMs =
        topicDetails.retentionMs ?
            (topicDetails.retentionMs.value === -1 ? 'Indefinite' : topicDetails.retentionMs.value)
        :
            'Not Set';
    const retentionBytes =
        topicDetails.retentionBytes ?
            (topicDetails.retentionBytes === -1 ? 'unknown' : topicDetails.retentionBytes)
        :
            'Not Set';
    let cleanupPolicy;
    switch(topicDetails.cleanupPolicy) {
        case 'delete':
            cleanupPolicy = 'Delete';
            break;
        case 'compact':
            cleanupPolicy = 'Compact';
            break;
        case 'delete,compact':
            cleanupPolicy = 'Both';
            break;
        default:
            cleanupPolicy = 'Unknown';
            break;
    }

    const navigate = useNavigate();

    const handleTopicClick = () => {
        navigate(`/topics/${encodeURIComponent(topicDetails.topic)}`, {
            state: { topicDetails }
        });
    }


    return (
        <>
            <tr>
                <td>
                    <button
                        onClick={() => onToggleExpand(topicDetails.topic)}
                        style={{
                            background: 'none',
                            border: 'none',
                            cursor: 'pointer',
                            color: 'white',
                            fontSize: '16px',
                        }}
                    >
                        {expanded ? '▼' : '▶'}
                    </button>
                </td>
                <td>
                    <button onClick={handleTopicClick} className="topic-link">
                        {topicDetails.topic}
                    </button>
                </td>
                <td>{topicDetails.numPartitions}</td>
                <td>{topicDetails.replicationFactor}</td>
                <td>{cleanupPolicy}</td>
                <td>{isInternalTopic}</td>
                <td>{retentionMs}</td>
                <td>{retentionBytes}</td>
            </tr>

            {expanded && (
                <tr>
                    <td colSpan={8}>
                        <pre>
                            {JSON.stringify(topicDetails, null, 2)}
                        </pre>
                    </td>
                </tr>
            )}
        </>
    );
};

export default TopicDataTableCell;