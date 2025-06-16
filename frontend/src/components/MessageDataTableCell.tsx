import React from "react";
import {MessageDetails} from "../model/MessageDetails.ts";

type MessageDataTableCellProps = {
    messageDetails: MessageDetails;
    onToggleExpand: (partition: number, offset: number) => void;
    expanded: boolean;
};

const MessageDataTableCell: React.FC<MessageDataTableCellProps> = ({ messageDetails, onToggleExpand, expanded }) => {
    return (
        <>
            <tr>
                <td>
                    <button
                        onClick={() => onToggleExpand(messageDetails.partition, messageDetails.offset)}
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
                <td>{messageDetails.offset}</td>
                <td>{messageDetails.partition}</td>
                <td>{messageDetails.timestamp}</td>
                <td>{messageDetails.key}</td>
                <td>{messageDetails.value}</td>
            </tr>

            {expanded && (
                <tr>
                    <td colSpan={8}>
                        <pre>
                            {JSON.stringify(messageDetails, null, 2)}
                        </pre>
                    </td>
                </tr>
            )}
        </>
    );
};

export default MessageDataTableCell;