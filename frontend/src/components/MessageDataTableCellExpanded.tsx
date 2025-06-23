import React from "react";
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
    console.log("Expanded Row State:", expandedRowState);
    const defaultVal = {
        'object': 'Object',
        'key': 'Key',
        'value': 'Value'
    }[expandedRowState ?? 'object'] || 'Object';
    console.log("Default Value:", defaultVal);
    console.log("Message Details:", messageDetails);
    return (
        <>
            <select defaultValue={defaultVal}>
                <option value="Object"> Object </option>
                <option value="Key"> Key </option>
                <option value="Value"> Value </option>
            </select>
            {
                expandedRowState === 'object' ? (
                    <div className="p-2">
                        <div className="text-xs text-gray-500 mb-2">Message Details:</div>
                    </div>
                ) : expandedRowState === 'key' ? (
                    <div className="p-2">
                        <div className="text-xs text-gray-500 mb-2">Key (JSON):</div>
                        <JsonViewer jsonData={messageDetails.keyJsonPayload} />
                    </div>
                ) : expandedRowState === 'value' ? (
                    <div className="p-2">
                        <div className="text-xs text-gray-500 mb-2">Value (JSON):</div>
                        <JsonViewer jsonData={messageDetails.valueJsonPayload} />
                    </div>
                ) : (
                    <div className="p-2">
                        <div className="text-xs text-gray-500 mb-2">Select an option to view details</div>
                    </div>
                )
            }
        </>
    )
};

export default MessageDataTableCellExpanded;