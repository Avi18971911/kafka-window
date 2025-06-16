import React from "react";
import styles from "../styles/DataTable.module.css";
import {MessageDetails} from "../model/MessageDetails.ts";
import MessageDataTableCell from "./MessageDataTableCell.tsx";

type MessageDataTableProps = {
    messages: MessageDetails[]
}

const mapPartitionAndOffsetToKey = (partition: number, offset: number): string => {
    return `${partition}-${offset}`;
}

const MessageDataTable: React.FC<MessageDataTableProps> = ({ messages }) => {
    const [expandedRow, setExpandedRow] = React.useState<string | null>(null);

    const toggleExpand = (partition: number, offset: number) => {
        const id = mapPartitionAndOffsetToKey(partition, offset);
        setExpandedRow((prev) => (prev === id ? null : id));
    };

    return (
        <div className={styles.tableContainer}>
            <table className={styles.table}>
                <thead>
                    <tr>
                        <th></th>
                        <th>Offset</th>
                        <th>Partition</th>
                        <th>Timestamp</th>
                        <th>Key</th>
                        <th>Value</th>
                    </tr>
                </thead>
                <tbody>
                    {messages.map((messageDetails) => (
                        <MessageDataTableCell
                            key={mapPartitionAndOffsetToKey(messageDetails.partition, messageDetails.offset)}
                            messageDetails={messageDetails}
                            onToggleExpand={toggleExpand}
                            expanded={expandedRow === mapPartitionAndOffsetToKey(
                                messageDetails.partition,
                                messageDetails.offset,
                            )}
                        />
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export default MessageDataTable;