import React from "react";
import {MessageDetails} from "../model/MessageDetails.ts";
import {ExpandedRowState} from "./MessageDataTable.tsx";
import MessageDataTableCellExpanded from "./MessageDataTableCellExpanded.tsx";

type MessageDataTableCellProps = {
    messageDetails: MessageDetails;
    onToggleExpand: (partition: number, offset: number, expandedRowState: ExpandedRowState) => void;
    expanded: boolean;
    expandedRowState?: ExpandedRowState;
};

const MessageDataTableCell: React.FC<MessageDataTableCellProps> = (
    { messageDetails, onToggleExpand, expanded, expandedRowState }
) => {
    const maxLength = 45;
    const shouldTruncateKey = messageDetails.key && messageDetails.key.length > maxLength;
    const shouldTruncateValue = messageDetails.value && messageDetails.value.length > maxLength;
    return (
        <>
            <tr>
                <td>
                    <button
                        onClick={() => onToggleExpand(messageDetails.partition, messageDetails.offset, 'object')}
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
                <td>
                    <button
                        onClick={() => onToggleExpand(messageDetails.partition, messageDetails.offset, 'key')}
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
                    {shouldTruncateKey ? messageDetails.key.substring(0, maxLength) + "..." : messageDetails.key}
                </td>
                <td>
                    <button
                        onClick={() => onToggleExpand(messageDetails.partition, messageDetails.offset, 'value')}
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
                    {shouldTruncateValue ? messageDetails.value.substring(0, maxLength) + "..." : messageDetails.value}
                </td>
            </tr>

            {expanded && (
                <tr>
                    <td colSpan={8}>
                        <MessageDataTableCellExpanded
                            messageDetails={messageDetails}
                            expandedRowState={expandedRowState}
                        />
                    </td>
                </tr>
            )}
        </>
    );
};

export default MessageDataTableCell;