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

    const handleOnChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setCurrentRowState(event.target.value.toLowerCase() as ExpandedRowState);
    };

    const stringCurrentRowState = currentRowState ?? 'object';

    return (
        <>
            <div style={{ display: 'flex', borderBottom: '1px solid #444', marginBottom: '1rem' }}>
                {['object', 'key', 'value'].map(option => (
                    <label key={option} style={{
                        padding: '0.5rem 1rem',
                        backgroundColor: stringCurrentRowState === option ? '#555' : 'transparent',
                        color: 'white',
                        cursor: 'pointer',
                        borderBottom: stringCurrentRowState === option ? '2px solid #fff' : '2px solid transparent',
                        textTransform: 'capitalize'
                    }}>
                        <input
                            type="radio"
                            value={option}
                            checked={stringCurrentRowState === option}
                            onChange={handleOnChange}
                            style={{ display: 'none' }}
                        />
                        {option}
                    </label>
                ))}
            </div>
            {
                currentRowState === 'object' ? (
                    <div style={{ padding: '0.5rem' }}>
                    <pre style={{
                        fontSize: '0.75rem',
                        backgroundColor: 'none',
                        padding: '0.75rem',
                        borderRadius: '0.25rem',
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