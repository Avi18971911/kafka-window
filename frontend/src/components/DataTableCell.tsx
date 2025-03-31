import React from "react";
import {TopicDetails} from "../model/TopicDetails.ts";

type DataTableCellProps = {
    topicDetails: TopicDetails;
    onToggleExpand: (topic: string) => void;
    expanded: boolean;
};

const DataTableCell: React.FC<DataTableCellProps> = ({ topicDetails, onToggleExpand, expanded }) => {
    const isInternalTopic = topicDetails.isInternal ? 'Yes' : 'No';
    const retentionMs = topicDetails.retentionMs.indefinite ? 'indefinite' : topicDetails.retentionMs.value;
    const retentionBytes = topicDetails.retentionBytes === -1 ? 'unknown' : topicDetails.retentionBytes;
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
                <td>{topicDetails.topic}</td>
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

export default DataTableCell;