import React from "react";
import styles from "../styles/DataTable.module.css";
import {MessageDetails} from "../model/MessageDetails.ts";
import MessageDataTableCell from "./MessageDataTableCell.tsx";

type MessageDataTableProps = {
    messages: MessageDetails[]
}

const MessageDataTable: React.FC<MessageDataTableProps> = ({ messages }) => {
    const [expandedRow, setExpandedRow] = React.useState<number | null>(null);

    const toggleExpand = (id: number) => {
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
                            key={messageDetails.partition}
                            messageDetails={messageDetails}
                            onToggleExpand={toggleExpand}
                            expanded={expandedRow === messageDetails.partition}
                        />
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export default MessageDataTable;