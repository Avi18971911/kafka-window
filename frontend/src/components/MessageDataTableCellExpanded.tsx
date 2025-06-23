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
                    <div className="p-2">
                        <div className="text-xs text-gray-500 mb-2">
                            {JSON.stringify(messageDetails, null, 10)}
                        </div>
                    </div>
                ) : currentRowState === 'key' ? (
                    <div className="p-2">
                        <JsonViewer jsonData={messageDetails.keyJsonPayload} />
                    </div>
                ) : currentRowState === 'value' ? (
                    <div className="p-2">
                        <JsonViewer jsonData={messageDetails.valueJsonPayload} />
                    </div>
                ) : (
                    <div className="p-2">
                        <div className="text-xs text-gray-500 mb-2">Select an option to view details</div>
                    </div>
                )
            }
        </>
    );
};

export default MessageDataTableCellExpanded;