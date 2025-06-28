import React, {useState} from "react";
import {MessageDetails} from "../model/MessageDetails.ts";
import {ExpandedRowState} from "./MessageDataTable.tsx";
import JsonViewer from "./JsonViewer.tsx";

type MessageDataTableCellExpandedProps = {
    messageDetails: MessageDetails;
    expandedRowState?: ExpandedRowState;
};

const MessageDataTableCellExpanded: React.FC<MessageDataTableCellExpandedProps> = (
    { messageDetails, expandedRowState }
) => {
    const [currentRowState, setCurrentRowState] = useState<ExpandedRowState>(
        expandedRowState ?? 'object'
    );

    const handleOnChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        setCurrentRowState(event.target.value.toLowerCase() as ExpandedRowState);
    };

    const stringCurrentRowState = currentRowState ?? 'object';

    return (
        <>
            <select value={stringCurrentRowState} onChange={handleOnChange}>
                <option value="object">Object</option>
                <option value="key">Key</option>
                <option value="value">Value</option>
            </select>
            {
                currentRowState === 'object' ? (
                    <div style={{ padding: '0.5rem' }}>
                    <pre style={{
                        fontSize: '0.75rem',
                        backgroundColor: 'none',
                        padding: '0.75rem',
                        borderRadius: '0.25rem',
                        border: '1px solid #374151',
                        color: '#e5e7eb',
                        whiteSpace: 'pre-wrap'
                    }}>
                        {JSON.stringify(messageDetails, null, 2)}
                    </pre>
                    </div>
                ) : currentRowState === 'key' ? (
                    <div style={{ padding: '0.5rem' }}>
                        <JsonViewer jsonData={messageDetails.keyJsonPayload} />
                    </div>
                ) : currentRowState === 'value' ? (
                    <div style={{ padding: '0.5rem' }}>
                        <JsonViewer jsonData={messageDetails.valueJsonPayload} />
                    </div>
                ) : (
                    <div style={{ padding: '0.5rem' }}>
                        <div style={{
                            fontSize: '0.75rem',
                            color: '#6b7280',
                            marginBottom: '0.5rem'
                        }}>
                            Select an option to view details
                        </div>
                    </div>
                )
            }
        </>
    );
};

export default MessageDataTableCellExpanded;